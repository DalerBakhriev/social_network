FROM mysql:8.0

RUN rm -r -f /docker-entrypoint-initdb.d/
ADD init.sql /docker-entrypoint-initdb.d/