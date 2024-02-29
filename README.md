# ms template
Template x go microservices

## Badges
[![Latest Release](https://git.cdlan.net/applications/uservices/ms-template-go/-/badges/release.svg)](https://git.cdlan.net/applications/uservices/ms-template-go/-/releases)

## Deploy instruction
The microservice will be packaged in a docker image, to deploy use the provided configs in the [deployments](./deployments) folder

## Develop Requirements
- install [go](https://go.dev/dl/)
> Make sure to install the version indicated in the [go.mod](./go.mod) file
- Docker for building

# Instructions - README
## Setup repo

- [ ] Edit `go.mod` with the name (like `git.cdlan.net/<group>/<reponame>`) of your project (and update all import statements)
- [ ] Edit `serviceName` in [server.go](cmd/server/server.go) with the name of the service
- [ ] Edit `app.image` in [docker-compose.yaml](deployments/docker/docker-compose.yaml)
- [ ] (optional) Uncomment all rows from [.gitlab-ci.yml](.gitlab-ci.yml) to enable pipelines
- [ ] Register project in [sonarqube](https://sonar.cub-otto.it/)

## Quick Start
1. Edit .proto files in [api](api/) folder
2. Run this command to generate go code from the protos
```shell
make grpc
```
3. In [grpc](internal/grpc) create a file for each service that you defined and implement the service servers and add a NewXYZServer() that return a pointer to the server
4. In [server.go](cmd/server/server.go) register the newly created servers

## Docs
We use [mkdocs](https://www.mkdocs.org/) for documentation
1. Update the content of [mkdocs.yml](mkdocs.yml) for title and url
2. Place md files inside the [docs](docs/) folder eventually in subfolders
3. Run this command to build the container locally and serve it
```shell
make docs
```

## Tests
We use [testify](https://github.com/stretchr/testify) for tests.
Each package should have this two test suite defined:
```go
type IntegrationTestSuite struct {
	suite.Suite
}

type UnitTestSuite struct {
	suite.Suite
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
```

- **UnitTestSuite** should contains all tests that can be run on code alone (no db, no ext services connections) these can always be run inside the CI pipeline
- **IntegrationTestSuite** contains all tests that should be run with some dependencies (db/ext-services/...)

To add a new test inside a suite:
```go
func (suite *IntegrationTestSuite) TestCustomerGroupDAO_RmNonExistentResFromCustomerGroup() {

	res := model.CustomerGroupResource{
		CustomerGroupId: 99,
		Resource:        "test",
		AccessType:      "read",
	}

	dao := CustomerGroupDAO{}
	ctx := context.Background()

	// rm res
	err := dao.RemoveResourceAccess(ctx, res)
	assert.NotNil(suite.T(), err)
}
```