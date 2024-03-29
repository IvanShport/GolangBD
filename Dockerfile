FROM golang:1.12.6-alpine AS build

WORKDIR /opt

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -ldflags="-w -s" -mod vendor -o forum

FROM ubuntu:18.04 AS release

MAINTAINER Artem V. Navrotskiy

#
# Установка postgresql
#
ENV PGVER 10

RUN apt-get -y update &&\
    apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
# EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

# Объявлем порт сервера
EXPOSE 5000

# Собранный ранее сервер
WORKDIR /opt
COPY --from=build /opt/forum forum
COPY migrations migrations
#
# Запускаем PostgreSQL и сервер
#
CMD service postgresql start && ./forum