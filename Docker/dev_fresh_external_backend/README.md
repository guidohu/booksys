# About dev_fresh

`dev_fresh_external` is a setup that can be used to develop the project and have a clean setup every time that leads you through the setup procedure.

# Usage

Start with a fresh deployment:

```
docker compose up
```

Take down the current deployment and remove the containers:

```
docker compose down
```

The above command will run the database and traefik. To run the backend,
either manually run it (run on port 9090) or execute:

```
./run_backend.sh
```

Traefik is configured such that the backend is expected to run on port :9090.

To start the frontend run:

```
npm run server-port
```

This will run the frontend on port 8080.

# Notes

Whenever changes are saved, npm will rebuild the frontend app.

For changes in the backend, you need to manually rebuild and restart the backend.