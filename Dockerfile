FROM ubuntu:18.04 

ENV PGSQLVER 10 
ENV DEBIAN_FRONTEND 'noninteractive' 

RUN echo 'Europe/Moscow' > '/etc/timezone' 

RUN apt-get -y update 
RUN apt install -y gcc git wget 
RUN apt install -y postgresql-$PGSQLVER 

RUN wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz 
RUN tar -xvf go1.11.2.linux-amd64.tar.gz 
RUN mv go /usr/local 

ENV GOROOT /usr/local/go 
ENV GOPATH /opt/go 
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH 

WORKDIR /escapade
COPY . . 

EXPOSE 5000 

USER postgres 

RUN /etc/init.d/postgresql start &&\ 
psql --echo-all --command "CREATE USER rolepade WITH SUPERUSER PASSWORD 'escapade';" &&\ 
createdb -O rolepade escabase &&\ 
psql --echo-all --command "\\c escabase;" &&\
psql --echo-all --command "\\i /escapade/internal/database/pgsql/create.pgsql;" &&\ 
psql --echo-all --command "\\i /escapade/internal/database/pgsql/test_insert.pgsql;" 

RUN echo "host all all 0.0.0.0/0 md5" » /etc/postgresql/$PGSQLVER/main/pg_hba.conf &&\ 
echo "listen_addresses='*'" » /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\ 
echo "fsync = off" » /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\ 
echo "synchronous_commit = off" » /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\ 
echo "shared_buffers = 512MB" » /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\ 
echo "random_page_cost = 1.0" » /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\ 
echo "wal_level = minimal" » /etc/postgresql/$PGSQLVER/main/postgresql.conf &&\ 
echo "max_wal_senders = 0" » /etc/postgresql/$PGSQLVER/main/postgresql.conf 

RUN service postgresql restart

EXPOSE 5432 

USER root 

RUN service postgresql start 
RUN /etc/init.d/postgresql start
RUN go run main.go 

#CMD DSN: "db://postgres:postgres@db:5432/postgres?sslmode=disable"
#RUN cd internal/services/api 
#RUN cd internal/services/api &&\ 
#  go test -coverprofile test/cover.out &&\
#  go tool cover -html=test/cover.out -o test/coverage.html

#RUN go test ./.../

# docker rmi escapade # удалим исходник(image)
# docker rm esc_cont # удалим контейнер
# docker build -t esc_image .
# docker run --name escemon -d esc_image /bin/sh -c "while true; do echo hello world; sleep 1; done"
# docker run -p 5000:5000 --name esc_image -t esc_cont 
# docker start -a esc_cont
# docker exec -it <name/id> /bin/sh

# launch!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
# docker build -t esc_image .
# docker run --name escemon -d esc_image /bin/sh -c "while true; do echo hello world; sleep 1; done"
# docker exec -it escemon /bin/sh
# /etc/init.d/postgresql start
# cd internal/services/api
# go test -coverprofile test/cover.out
