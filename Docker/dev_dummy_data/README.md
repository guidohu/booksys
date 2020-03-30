# Development Environment

This deployment instructions are not recommended for a production setup due to various issues:

- No secure passwords for database access are used
- Dummy data is already be added to the database
- No data persistence for database content

## Pre-requisites
Make sure that the two docker images `booksys_ui` and `booksys_dummy_db` exist in one of your docker registries. On how to build and publish them, please refer to the respective documentation.

## Run the development environment
Run:
```
docker-compose up -d
```

Go to `localhost` or the IP under which you run your docker environment.

Login using credentials:
```
email1@domain1.com / password123
email2@domain2.com / password123
...
```