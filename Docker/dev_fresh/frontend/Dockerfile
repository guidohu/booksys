FROM node:24

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
