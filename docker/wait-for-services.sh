#! /bin/sh

# Wait for Mongo
until nc -z -v -w30 $DATABASE_HOST $DATABASE_PORT
do
  echo 'Waiting for Mongo...'
  sleep 1
done
echo "Mongo is up and running"
