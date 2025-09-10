# mailinglist-backend-go
A backend service to manage mailing list subscriptions with Mailgun.

## Swagger / OpenAPI documentation
The project includes Swagger (OpenAPI) annotations and can generate a swagger.json file both locally and automatically during the Docker image build.

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

### Swagger generated during Docker build
The Dockerfile is configured to generate the Swagger JSON during the image build. In the builder stage, it runs:

- `swag init -g main.go -o ./swagger -ot json`

The generated files are copied into the final image at `/app/swagger/swagger.json`.

Example build and ways to access the file:

- Build the image:
  - `docker build -t mailinglist-backend-go:latest .`

- Print the swagger.json from the built image:
  - `docker run --rm mailinglist-backend-go:latest cat /app/swagger/swagger.json`

- Or copy it out of a running container:
  - `CID=$(docker create mailinglist-backend-go:latest)`
  - `docker cp "$CID":/app/swagger/swagger.json ./swagger-from-image.json`
  - `docker rm "$CID"`
