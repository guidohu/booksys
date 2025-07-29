#!/bin/bash

cd /usr/app
rm -rf node_modules
rm package-lock.json
npm install
npm run serve-port