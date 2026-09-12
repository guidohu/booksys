#!/bin/bash

source .env
export BOOKSYS_MYNAUTIQUE_API_KEY=${BOOKSYS_MYNAUTIQUE_API_KEY}
export BOOKSYS_DATABASE_HOST=${LOCAL_DATABASE_HOST}
export BOOKSYS_DATABASE_PORT=${LOCAL_DATABASE_PORT}
export BOOKSYS_DATABASE_DBNAME=${MYSQL_DATABASE}
export BOOKSYS_DATABASE_USER=${MYSQL_USER}
export BOOKSYS_DATABASE_PASSWORD=${MYSQL_PASSWORD}

# Set your module path (e.g., github.com/user/project)
export MODULE_PATH="server" 

# Get version/commit/date
export BOOKSYS_VERSION=$(git describe --tags --abbrev=0)
export BOOKSYS_COMMIT=$(git rev-parse --short HEAD)
export BOOKSYS_DATE=$(date +'%Y-%m-%d_%T%z')
export BOOKSYS_ENVIRONMENT="local_dev"

dir=`pwd`
cd ../../backend/api/v2
go run -ldflags="-X '${MODULE_PATH}/version.Release=${BOOKSYS_VERSION}' -X ${MODULE_PATH}/version.Commit=${BOOKSYS_COMMIT} -X ${MODULE_PATH}/version.BuildDate=${BOOKSYS_DATE}" \
   main.go --http_port 9090 --http_websetup --http_secure_cookie=false
cd ${dir}