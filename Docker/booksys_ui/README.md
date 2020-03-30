# Create Docker Container

Use the Dockerfile and adjust as you like, then build the container

`
docker build -t booksys_ui .
`

# Create Deployment

Use `docker-compose` to do a server deployment. Adjust docker-compose.yaml as you like, then start the deployment.

`
docker-compose up -d
`

Note: this deployment will not work without a given database backend.
