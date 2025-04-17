package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"io"
	"net/http"
	"os"
	"time"
)

type Product struct {
	Id        int64
	CreatedAt int64
	Name      string
}

func DeleteProducts(rw http.ResponseWriter, req *http.Request) {
	timeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rw.Header().Set("Content-Type", "application/json")

	l := zerolog.New(os.Stdout).With().Stack().Timestamp().Logger()

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		l.Error().Stack().Err(err).Msg("cannot get request body")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "cannot get request body"}`))

		return
	}

	var input Product

	err = json.Unmarshal(rawBody, &input)
	if err != nil {
		l.Error().Stack().Err(err).Msg("cannot unmarshal request body")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "cannot unmarshal request body"}`))

		return
	}

	creds := ycsdk.InstanceServiceAccount()
	token, err := creds.IAMToken(timeCtx)
	if err != nil {
		l.Error().Err(err).Msg("failed to get token")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "failed to get token"}`))

		return
	}
	l.Info().Str("token", token.GetIamToken()).Msg("got IAM token")

	database := os.Getenv("DATABASE_URL")
	l.Info().Str("database", database).Msg("database url")

	db, err := ydb.Open(timeCtx,
		database,
		ydb.WithAccessTokenCredentials(token.GetIamToken()),
	)
	if err != nil {
		l.Error().Err(err).Msg("failed to connect to database")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "failed to connect to database"}`))

		return
	}

	result, err := db.Query().Query(timeCtx, fmt.Sprintf("delete from products where id = %v;", input.Id))
	if err != nil {
		l.Error().Err(err).Msg("failed to query products")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(`{"error": "failed to query products"}`))

		return
	}
	defer func() {
		result.Close(timeCtx)
	}()

	response, _ := json.Marshal(struct {
		Data string `json:"data"`
	}{
		Data: "product",
	})

	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
	return
}
