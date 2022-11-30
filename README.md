# dittofunc

HTTP server implementation of Ditto <https://github.com/juanvillacortac/ditto> code generation.

## Server

The Ditto `httpserver.Handler` is a simple HTTP handler that receives a Ditto configuration file from GitHub and generates the code.

### Self-hosting

Requirements:

- [Go](https://go.dev/dl/)

```bash
go run .
```

## Deployment