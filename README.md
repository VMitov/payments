# Payment API Exacmple

## Start db with schema
```
docker run \
    --mount type=bind,source="$(pwd)/database",target=/docker-entrypoint-initdb.d \
    --publish 5432:5432 \
    postgres
```

## Start the service in container
```
docker build --tag=payments .
docker run \
    --network=host \
    payments -db "postgres://postgres@127.0.0.1:5432/payments?sslmode=disable"
```

## Run tests
```
go test ./...
```

### in Docker container
```
docker run --rm \
    --mount type=bind,source="$(pwd)",target=/go/src/github.com/VMitov/payments
    golang go test github.com/VMitov/payments/...
```
