FROM postgis/postgis:16-3.4-alpine
COPY migrations/*.sql /docker-entrypoint-initdb.d/

EXPOSE 5432
