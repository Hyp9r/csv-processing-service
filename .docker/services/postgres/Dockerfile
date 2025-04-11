FROM postgres:13.15-alpine

# Set enviroment variables
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=backend

COPY ./.docker/services/postgres/init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

CMD [ "postgres" ]