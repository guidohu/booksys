# About dev_fresh

`dev_fresh` is a setup that can be used to develop the project and have a clean setup every time that leads you through the setup procedure.

# Usage

Start with a fresh deployment:

```
docker compose up --build
```

Take down the current deployment and remove the containers:

```
docker compose down
```

# Notes

This setup will let you make live changes on the frontend as the frontend code is mapped into the frontend container. Whenever changes are saved, npm will rebuild the frontend app.

The `node_modules` of the repository are not used by the frontend container, as they
contain binaries for the platform of the host (e.g. `@rollup/rollup-darwin-arm64` on a
Mac) which do not run on linux. The container has its own `node_modules` in the
`frontend_node_modules` volume, which is filled by `npm install` on every start of the
container. The first start therefore takes a while. To throw those packages away, remove
the volume:

```
docker compose down -v
```