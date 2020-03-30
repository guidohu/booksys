FROM mysql:5.7
COPY testdata_anonym.sql /docker-entrypoint-initdb.d/testdata.sql


EXPOSE 3306 33060

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["mysqld"]

