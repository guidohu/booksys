# About booksys_db

This is a sample MySQL database that can be used as a database backend for the Booksys App.

# Create Docker Container

Use the Dockerfile and adjust as you like, then build the container.

`
docker build -t booksys_db .
`

# Create Deployment

Use `docker-compose` to do a database deployment. Adjust docker-compose.yaml as you like (especially passwords), then start the deployment.

`
docker-compose up -d
`
