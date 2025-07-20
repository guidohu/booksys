# About dev_fresh

`dev_fresh` is a setup that can be used to develop the project and have a clean setup every time that leads you through the setup procedure.

# Usage

Start with a fresh deployment:

```
docker compose up
```

Take down the current deployment and remove the containers:

```
docker compose down
```

# Notes

This setup will let you make live changes on the frontend as the frontend code is mapped into the frontend container. Whenever changes are saved, npm will rebuild the frontend app.