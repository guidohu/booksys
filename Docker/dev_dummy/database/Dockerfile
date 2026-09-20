# Keep this in step with Docker/dev_fresh/database/Dockerfile. mysql:5.7 is end
# of life and has no arm64 image, so it cannot be built on an Apple Silicon
# host. The dump below is a mysqldump from a 5.7 server, which 8.4 still reads:
# the display widths (int(11)) and the utf8/utf8_unicode_ci aliases it contains
# are deprecated there but accepted. The backend migrates the schema it creates
# to the current one on its first connection.
FROM mysql:8.4
COPY testdata_anonym.sql /docker-entrypoint-initdb.d/testdata.sql

EXPOSE 3306 33060

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["mysqld"]
