# About dev_dummy

`dev_dummy` is a setup that can be used to develop the project and have a setup that is already filled with dummy data.

# Usage

Start with a dummy deployment:

```
docker compose up --build
```

The app is then served on http://localhost, the traefik dashboard on http://localhost:8080.
Log in with any of the anonymized accounts of the dummy data set, e.g.
`email1@domain1.com` / `password123` (see `database/README.md`).

Take down the current deployment and remove the containers:

```
docker compose down
```

# Notes

The dummy data in `database/testdata_anonym.sql` is a dump of a MySQL 5.7 server, while
the `db` container runs MySQL 8.4. The dump still imports there, and the backend migrates
the schema it creates to the current one on its first connection, so the log of `api_v2`
shows a run of migration steps on every fresh start.

This setup will let you make live changes on the frontend as the frontend code is mapped into the frontend container. Whenever changes are saved, vite will rebuild the frontend app.

The `node_modules` of the repository are not used by the frontend container, as they
contain binaries for the platform of the host (e.g. `@rollup/rollup-darwin-arm64` on a
Mac) which do not run on linux. The container has its own `node_modules` in the
`frontend_node_modules` volume, which is filled by `npm install` on every start of the
container. The first start therefore takes a while. To throw those packages away, remove
the volume:

```
docker compose down -v
```
