# build the app with node 24
FROM node:24-alpine AS build-stage

WORKDIR /home/node/app

COPY package*.json     ./
RUN npm install

COPY . .
RUN npm run build

# run the frontend behind nginx
FROM nginx:1.27-alpine AS production-stage

RUN addgroup -S booksys && adduser -S -G booksys booksys

COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build-stage /home/node/app/dist /app

RUN touch /tmp/nginx.pid && \
    chown -R booksys:booksys /var/cache/nginx /var/log/nginx /etc/nginx /tmp/nginx.pid && \
    chmod -R 755 /var/log/nginx && \
    chmod 644 /tmp/nginx.pid

USER booksys

EXPOSE 80