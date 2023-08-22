# ms template
Template x go microservices


## Setup repo

1. Edit `go.mod` with the name of your project (and update all import statements)
2. Edit `serviceName` in [server.go](cmd/server/server.go) with the name of the service
3. Edit `app.image` in [docker-compose.yaml](deployments/docker/docker-compose.yaml)
4. Uncomment all rows from [.gitlab-ci.yml](.gitlab-ci.yml) to enable pipelines

## Start developing
1. Create .proto files in [api](api/) folder
2. Run script x classes generation [gen_grpc_classes.sh](scripts/gen_grpc_classes.sh)
3. In [grpc](internal/grpc) create a file for each service that you defined and implement the service servers and add a NewXYZServer() that return a pointer to the server
4. In [server.go](cmd/server/server.go) register the newly created servers

## Tips
- If you need types or operations from the db create them inside the [database](internal/database) package