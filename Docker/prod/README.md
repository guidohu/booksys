# Production Environment

## Pre-requisites
Make sure that the two docker images `booksys_ui` and `booksys_db` exist in one of your docker registries. On how to build and publish them, please refer to the respective documentation.

## Run the development environment
Adjust parameters following parameters:
- MYSQL_DATABASE
- MYSQL_USER (recommended)
- MYSQL_PASSWORD (recommended)
- MYSQL_ROOT_PASSWORD (recommended)
- DB_NAME (same as MYSQL_DATABASE)
- DB_USER (same as MYSQL_USER)
- DB_PASSWORD (same as MYSQL_PASSWORD)

Run:
```
docker-compose up -d
```

This command runs the database and the UI server on your system, whereas the database volume will be persistent.