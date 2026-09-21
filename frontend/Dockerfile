# build the app with node 24
FROM node:24-alpine AS build-stage

WORKDIR /home/node/app

COPY package*.json     ./
RUN npm install

COPY . .

# Stamp the image's version into package.json. BOOKSYS_VERSION carries the git
# ref name, so a release tag arrives as v3.5.0 and the leading v is dropped to
# keep the field a plain npm version. The ARG sits after npm install so that a
# new tag does not invalidate the dependency layer, and after COPY . . because
# that copy restores the committed package.json. Without the arg (a local
# docker build) the committed version is left alone.
ARG BOOKSYS_VERSION
RUN if [ -n "${BOOKSYS_VERSION}" ]; then \
        npm pkg set version="${BOOKSYS_VERSION#v}"; \
    fi

RUN npm run build

# run the frontend behind nginx
FROM nginx:1.30-alpine AS production-stage

RUN addgroup -S booksys && adduser -S -G booksys booksys

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build-stage /home/node/app/dist /app

RUN touch /tmp/nginx.pid && \
    chown -R booksys:booksys /var/cache/nginx /var/log/nginx /etc/nginx /tmp/nginx.pid && \
    chmod -R 755 /var/log/nginx && \
    chmod 644 /tmp/nginx.pid

USER booksys

EXPOSE 80