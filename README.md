# Wake and Surf Booksys

![Build](https://github.com/guidohu/booksys/workflows/Docker%20Image%20CI/badge.svg) ![CodeQL](https://github.com/guidohu/booksys/workflows/CodeQL/badge.svg)

A booking and management tool small boat communities can use to schedule boat usage and simplify maintenance as well as accounting.

![dashboard view](https://raw.githubusercontent.com/guidohu/booksys/master/docs/img/dashboard.png)

It comes with the following features:
- **Boat Management** Track engine hours, fuel consumption and maintenance logs.
- **Session Management** Create sessions, invite users to sessions, track time behind the boat and much more.
- **Multi User Support** with different roles such as Administrators, Members, Guests and flexible pricing models for different users.
- **Accounting** to track the finances of the boat community.

## Installation

### Prerequisites
It is strongly recommended to run the application behind a simple reverse proxy (nginx, apache, ...) to enforce https.

To have a simple deployment scenario it is suggested to run the application in a docker container as described below.

### Setup
To run it on your local machine, follow these steps:

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

Go to [http://localhost](http://localhost). You will be guided through the setup process.

### Security Considerations

- Change passwords in `docker/booksys_db/Dockerfile` before building the container.
- Expose the application through a reverse proxy to enforce HTTPS.
- Restrict access to `classes` and `config` folders, in case you do not run the default docker containers.

## Limitations

- Some minor configuration parameters can only be updated in the database directly (see configuration table).
- The setup should also work with mariaDB but this has not been tested.
- So far, no official payment system (paypal, ...) is integrated. Please get in contact with me in case this is required or feel free to submit a pull request.
- It is optimized for the use within a non commercial boat community. To allow guests to reserve slots, adjustments would be needed.
