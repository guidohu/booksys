# Wake and Surf Booksys

A booking and management tool for small boat communities to schedule boat usage and simplify maintenance as well as accounting.

![dashboard view](https://raw.githubusercontent.com/guidohu/booksys/master/img/dashboard.png)

## Installation

### Prerequisites
It is strongly recommended to run the application behind a simple reverse proxy (nginx, apache, ...) to enforce https.

### Setup
To run it on your machine, follow these steps:

```
# build the application itself
docker build -t booksys_ui .

# build a sample backend database (official mysql 5.7)
cd docker/booksys_db
docker build -t booksys_db .

# start your deployment
cd ../prod
echo '' > config.php
docker-compose up -d
```

Go to [http://localhost](http://localhost) assuming docker runs on your local machine. You will be guided through the setup process.

### Security Considerations

- Change passwords in `docker/booksys_db/Dockerfile` before building the container.
- Expose the application through a reverse proxy to enforce HTTPS.
- Restrict access to `classes` and `config` folders

## Limitations

- Customization of configuration parameters can only be done in the database directly (configuration table).
- Not every site is available as mobile and desktop version yet.
- Some third party components (such as highcharts.js) are using restricted licenses and thus the project cannot be used commercially currently.
- The setup should also work with mariaDB but this has not been tested.
