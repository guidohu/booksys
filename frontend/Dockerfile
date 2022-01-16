# build the app with node 16 LTS
FROM node:16 as build-stage

WORKDIR /home/node/app
COPY package*.json     ./
COPY babel.config.js   ./
COPY vue.config.js     ./
COPY src               ./src
COPY public            ./public
RUN npm install -g npm
RUN npm install
RUN npm run build

# run the frontend behind nginx
FROM nginx:alpine as production-stage
RUN mkdir /app
COPY --from=build-stage /home/node/app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf