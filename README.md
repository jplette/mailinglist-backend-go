# mailinglist-backend-go
A backend service to manage mailing list subscriptions with Mailgun.

## Docker images via GitHub Actions
This repo builds and pushes Docker images to Docker Hub via GitHub Actions:
- On Release (published): pushes two tags to Docker Hub – `latest` and the release tag (e.g., `v1.2.3`).
- On push to `main`: pushes a `dev` tag to Docker Hub.

Configure the following GitHub repository secrets before running the workflows:
- `DOCKERHUB_USERNAME`: Your Docker Hub username or organization.
- `DOCKERHUB_TOKEN`: A Docker Hub access token or password.
- `DOCKERHUB_REPO`: The full Docker Hub repository name, e.g., `yourname/mailinglist-backend-go`.

The workflows are located at:
- `.github/workflows/release-docker.yml`
- `.github/workflows/dev-docker.yml`

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
