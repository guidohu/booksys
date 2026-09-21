# Keep this in step with Docker/dev_fresh/frontend/Dockerfile.
FROM node:26

# npm's version is pinned by the `node:26` tag above, so that the dev stack and
# the production build (frontend/Dockerfile) use the same one. Upgrading npm
# from inside the container would drift from that, so silence the notice it
# prints on every container start, where `npm install` runs (docker-compose.yaml).
ENV NPM_CONFIG_UPDATE_NOTIFIER=false

# The frontend sources are bind mounted from the host into /home/node/app
# (see docker-compose.yaml). The node_modules of the host contain platform
# specific binaries (e.g. @rollup/rollup-darwin-arm64, @esbuild/darwin-arm64)
# which do not run inside this container. Therefore node_modules is masked
# with a named volume that is populated on container start.
# The directory has to exist in the image and be owned by the node user, so
# that docker creates the volume with the correct ownership.
RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app

WORKDIR /home/node/app

USER node
