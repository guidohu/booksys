# About dummy database creation

This is a database that contains synthetic data for debugging and development purposes.

# Create Docker Container

Use the Dockerfile and adjust as you like, then build the container.

`
docker build -t booksys_dummy_db .
`

# Import new production data to development

It is not recommended to work with production data, thus every database dump that is created with mysql_dump can be anonymized automatically using:

`
./anonymize_dump.sh -i <your_dump.sql> -o testdata_anonym.sql
`

After that, the docker image will be built using your anonymized test data. You can login with using the following credentials:

`
email1@domain1.com / password123
email2@domain2.com / password123
...
`
