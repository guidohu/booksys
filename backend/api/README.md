# Booksys Backend

The backend is written in Go and provides an API for the frontend.

```
UI   -->  Backend --> Database
```

## How to build the backend container

```
docker build -t booksys_api .
```

## How to run the backend container

```
docker run -d --name booksys_api --rm -p 9090:9090 booksys_api
```

## Run for development

```
docker run -d --name booksys_api --rm -p 9090:9090 \
    booksys_api
```