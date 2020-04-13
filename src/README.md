## Directory Structure Explained

- `.` website structure
- `./api` public facing PHP API
- `./bootstrap` framework
- `./classes` PHP classes used by the API
- `./config` application configuration and database schemas
- `./font-awesome` icons used in some parts of the application
- `./img` images
- `./js` javascript libraries, both 3rd party and own libraries
- `./res` mostly static HTML content and vue templates

## Run without Docker

The content of the `./src` folder can be installed into a webserver's directory. The webserver needs to support PHP (incl. mysqli extension). Also the webserver needs to have `sendemail` installed.
