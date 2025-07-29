FROM node:24

WORKDIR /usr/app

COPY setup.sh /setup.sh

RUN npm cache clean --force
RUN npm install --save bootstrap
RUN npm install --save bootstrap-vue
RUN npm install --save vite