# Payment API Exacmple

## Start db with schema
```
docker run \
    --mount type=bind,source="$(pwd)/database",target=/docker-entrypoint-initdb.d \
    --publish 5432:5432 \
    postgres
```
