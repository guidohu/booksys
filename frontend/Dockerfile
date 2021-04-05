FROM node:15 as build-stage

# copy project configuration
# COPY  package.json /home/node/app/
# COPY  vue.config.js /home/node/app/
# COPY  webpack.config.js /home/node/app/
# COPY  babel.config.js /home/node/app/

# copy source code and resources
# COPY  src    /home/node/app/
# COPY  public /home/node/app/

WORKDIR /home/node/app
COPY package*.json     ./
COPY babel.config.js   ./
COPY vue.config.js     ./
COPY src               ./src
COPY public            ./public
RUN npm install
RUN npm run build

# # setup vue build environment for the app
# RUN npm install -g @vue/cli
# RUN npm install
# RUN vue inspect

# # build the app
# RUN cd /home/node/app; npm run-script build



FROM nginx:alpine as production-stage
RUN mkdir /app
COPY --from=build-stage /home/node/app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf