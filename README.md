# ms template
Template x go microservices

## Badges
[![Latest Release](https://cdlab.cdlan.net/uservices/ms-template/-/badges/release.svg)](https://cdlab.cdlan.net/cdlan/uservices/ms-template/-/releases)

## Deploy instruction
The microservice is packaged in a docker image, to deploy use the provided configs in the [deployments](./deployments) folder

## Develop Requirements
- install [go](https://go.dev/dl/)
> Make sure to install the version indicated in the [go.mod](./go.mod) file

# Instructions - README
## Setup repo

- [ ] Edit `go.mod` with the name of your project (and update all import statements)
- [ ] Edit `serviceName` in [server.go](cmd/server/server.go) with the name of the service
- [ ] Edit `app.image` in [docker-compose.yaml](deployments/docker/docker-compose.yaml)
- [ ] (optional) Uncomment all rows from [.gitlab-ci.yml](.gitlab-ci.yml) to enable pipelines
- [ ] Register project in [sonarqube](https://sonar.cdlan.net/)

## Quick Start
1. Edit .proto files in [api](api/) folder
2. Run script x classes generation [gen_grpc_classes.sh](scripts/gen_grpc_classes.sh)
3. In [grpc](internal/grpc) create a file for each service that you defined and implement the service servers and add a NewXYZServer() that return a pointer to the server
4. In [server.go](cmd/server/server.go) register the newly created servers

## Tips
- If you need types or operations from the db create them inside the [database](internal/database) package