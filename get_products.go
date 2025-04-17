package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"io"
	"net/http"
	"os"
)

type Response struct {
	StatusCode int `json:"statusCode"`
	Body       any `json:"body"`
}

type Product struct {
	Id        int64
	CreatedAt int64
	Name      string
}

func GetProducts(ctx context.Context) (*Response, error) {
	l := zerolog.New(os.Stdout).With().Stack().Timestamp().Logger()

	creds := ycsdk.InstanceServiceAccount()
	token, err := creds.IAMToken(ctx)
	if err != nil {
		l.Error().Err(err).Msg("failed to get token")
		return nil, fmt.Errorf("failed to get IAM token: %w", err)
	}

	l.Info().Str("token", token.GetIamToken()).Msg("got IAM token")

	database := os.Getenv("DATABASE_URL")
	l.Info().Str("database", database).Msg("database url")

	db, err := ydb.Open(ctx,
		database,
		ydb.WithAccessTokenCredentials(token.GetIamToken()),
	)
	if err != nil {
		l.Error().Err(err).Msg("failed to connect to database")

		return nil, err
	}

	result, err := db.Query().Query(ctx, "select * from products")
	if err != nil {
		l.Error().Err(err).Msg("failed to query products")

		return nil, err
	}
	defer func() {
		err = result.Close(ctx)
	}()

	var products []Product

	for {
		var product Product

		set, err := result.NextResultSet(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			l.Error().Err(err).Msg("failed to fetch set of products")
			return nil, err
		}

		for {
			row, err := set.NextRow(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				l.Error().Err(err).Msg("failed to fetch row of products")
				return nil, err
			}

			err = row.Scan(&product.Id, &product.CreatedAt, &product.Name)
			if err != nil {
				l.Error().Err(err).Msg("failed to scan products")
				return nil, err
			}

			products = append(products, product)
		}

	}

	return &Response{StatusCode: http.StatusOK, Body: products}, nil
}
