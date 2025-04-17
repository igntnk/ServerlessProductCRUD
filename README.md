This is a Go-based serverless application that provides CRUD (Create, Read, Update, Delete) operations for products using Yandex Database (YDB) as the backend.

## Features

- ✅ **Create** new products
- 🔍 **Retrieve** products (single or all)
- ✏️ **Update** existing products
- 🗑️ **Delete** products
- 🛡️ Built with **YDB SDK** for database operations
- 🔐 Uses **Yandex Cloud IAM** for authentication
- 📝 Structured logging with **zerolog**

## API Endpoints

Prerequisites

- 🐹 Go 1.21 or higher
- ☁️ Yandex Cloud account
- 🗄️ YDB database instance
- ⚙️ Environment variables:
    DATABASE_URL: Connection string for YDB database
