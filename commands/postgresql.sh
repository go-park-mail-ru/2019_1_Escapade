#!/bin/sh
echo ""
echo "  --------------------------"
echo "  ------PostgreSQL set------"
echo "  --------------------------"
echo ""
#chmod +x postgresql.sh && ./postgresql.sh

echo "  1. Create our postgreSQL user and database"
#/etc/init.d/postgresql start &&\
sudo -u postgres psql postgres postgres --echo-all --command \
    "CREATE USER rolepade WITH SUPERUSER PASSWORD 'escapade';
     CREATE DATABASE escabase OWNER rolepade;" &&\
#/etc/init.d/postgresql stop

echo "  2. Drop old tables"
PGPASSWORD=escapade psql \
    -h 127.0.0.1 -p 5432 -U rolepade -d escabase -f \
    "../internal/database/scripts/drop.psql"

echo "  3. Create new tables"
PGPASSWORD=escapade psql \
    -h 127.0.0.1 -p 5432 -U rolepade -d escabase -f \
    "../internal/database/scripts/create.psql"

#echo "  4. Look at tables"
#PGPASSWORD=escapade psql \
#     -h 127.0.0.1 -p 5432 -U rolepade -d escabase -f \
#    "../internal/database/scripts/look.psql"

echo "  ----------Done!-----------"