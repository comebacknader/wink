FROM postgres:11.3-alpine
ENV POSTGRES_USER wink-db-username
ENV POSTGRES_PASSWORD wink-db-password
ENV POSTGRES_DB wink-db-name

COPY init.sql /docker-entrypoint-initdb.d/


