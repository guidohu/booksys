# Docker Environment

## Components

### booksys_dummy_db
Is a docker container with dummy data in it, mainly used for running the dev environment with some data to play around.

### booksys_db
Is a docker container that contains a simple mysql DB. The application will work with a custom mysql server too, only the configuration parameters need to be changed accordingly.

### booksys_ui
Is the docker container that runs the UI.

### dev_dummy_data
Is a setup of both database and the UI and is also though to be non-persistent. It can be started for development purposes. It will use the booksys_dummy_db container that already contains dummy data.

### dev_fresh_setup
Is a non-persistent setup that can be used to test the setup process or to have a short look at the app without having anything persistent.

### prod
Is a simple production ready installation with data persistency configured. Credentials need to be changed for security reasons.