docker-compose up -d
source ./.env

# TODO: Fix this magic value
sleep 5
docker exec -it mongo bash -c "echo 'rs.initiate()' | mongo 'mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@127.0.0.1'"