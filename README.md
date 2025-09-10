# mailinglist-backend-go
A backend service to manage mailing list subscriptions with Mailgun.

## Swagger / OpenAPI documentation
The project includes Swagger (OpenAPI) annotations and can generate a swagger.json file locally.

### Generate Swagger locally
You can generate the Swagger JSON in the repository at `./swagger/swagger.json`.

- Quick one-liner (no global install required):
  - `go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -o ./swagger -ot json`

- Or install the `swag` CLI once and use it:
  1. `go install github.com/swaggo/swag/cmd/swag@latest`
  2. `swag init -g main.go -o ./swagger -ot json`

- Alternatively, you can use Go generate from the project root:
  - `go generate ./...`

Result:
- The generated file will be at `./swagger/swagger.json`.

Notes:
- Ensure you run the commands from the repository root (where `main.go` is).
- If you update handler comments or add endpoints, re-run the command to refresh the spec.
