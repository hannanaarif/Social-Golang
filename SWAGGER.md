# Swagger Documentation Guide

This guide explains how to install, configure, and run Swagger documentation for this project.

## 1. Installation

### Install `swag` CLI
To generate the documentation, you need the `swag` command-line tool.
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
Ensure your `GOBIN` (usually `~/go/bin`) is in your system `PATH`.

### Install Dependencies
The project uses `http-swagger` to serve the UI.
```bash
go get github.com/swaggo/http-swagger/v2
```

## 2. Generating Documentation

Run the following command from the project root to generate the `docs` folder:
```bash
swag init -g cmd/api/main.go
```
*Note: The `-g` flag points to the file containing your General API Info (usually where `main.go` is).*

## 3. Serving Swagger UI

In your router configuration (e.g., `cmd/api/api.go`), register the Swagger handler:

```go
import (
    _ "github.com/hannanaarif/Social/docs" // Import generated docs
    httpSwagger "github.com/swaggo/http-swagger/v2"
)

// ... inside your router
r.Get("/swagger/*", httpSwagger.Handler())
```

The documentation will be available at:
`http://localhost:8080/swagger/index.html`

## 4. Adding Annotations

### General API Info (in `main.go`)
```go
// @title           Social API
// @version         1.0
// @description     API Server for Social Application
// @host            localhost:8080
// @BasePath        /v1
```

### Handler Annotations
```go
// @Summary      Create a post
// @Description  Create a new post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        payload  body      createPostPayload  true  "Post payload"
// @Success      201      {object}  store.Post
// @Router       /posts [post]
```

## 5. Using the Makefile
You can run `make gen-docs` to regenerate the documentation.
