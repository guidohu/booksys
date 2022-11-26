# Booksys Backend

The backend is a collection of PHP scripts providing an API for the frontend.

```
UI   -->  Backend --> Database
```

## How to build the backend container

```
docker build -t booksys_api .
```

## How to run the backend container

```
docker run -d --name bs_api --rm -p 80:80 booksys_api
```

## Run for development

```
docker run -d --name bs_api --rm -p 80:80 \
    -v api:/var/www/html/api \
    -v classes:/var/www/html/classes \
    -v config:/var/www/html/config \
    booksys_api
```