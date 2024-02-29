# Go Microservices

## Project structure
Use this layout -> https://github.com/golang-standards/project-layout

## Docs folder
The docs folder contains all documentation for project (better if readable directly from gitlab/github)

## Dockerfile & docker-compose
- Add Dockerfile in `/build` folder (or subfolder)
- If the microservice has dependencies (DB or other ms) add a `docker-compose.yml` file in `/deployments/docker`

## Config
- Add in `/config` example file for configs (EG: if config is config.yml -> config-example.yml) with safe values (no prod values and no passwords)
- Add real config files to gitignore
- Add config from multiple source, usually I do:
  1. Load from Default
  2. Load from Yaml (overrides previous value only if found, on a per-var basis)
    - we usually use [viper](https://github.com/spf13/viper) with yaml, but it is not mandatory, if you want an example read [here](https://cdlab.cdlan.net/cdlan/users/-/blob/main/internal/config/config.go#L36)
  3. Load from ENV var (Docker) (Overrides previous values only if found, on a per-var basis)
  4. (I have never done it, But I have seen it suggested) Flags at launch (Optional)

- [Example](https://cdlab.cdlan.net/cdlan/users/-/blob/main/internal/config/config.go)
  - I have a struct for each configuration I need and It always have the same three methods:
    - Default() -> return instance of class with default values
    - class.loadVarsFromYaml() -> ONLY for ROOT config class, nested config structs will be loaded by top automatically
    - class.loadVarsFromEnv()
  - We keep the root config class in a global variable so that it can be accessed anywhere in the app. I don't think it is best practice but greatly simplify stuff

## Open Tracing
Use this package https://cdlab.cdlan.net/cdlan/users/-/tree/main/pkg/otel as an idea on how to implement otel x ms. Deployment will likely use OLTP exporter with [grafana tempo](https://grafana.com/oss/tempo/) or [jaeger](https://www.jaegertracing.io/)
- Each trace should correspond to the whole call (mistra -> ms1 -> ms2 ...)
- Each function should have a span with the function name (Ideally), errors should be logged by the caller function
- there should be a way to disable all telemetry so that it is simpler to test (I had a bool inside the package that bypass all functions)

## DB
We do not have preferences as to what to use for DB, but we have already used and deployed MariaDB and PostgreSQL

## Statefulness
The ms will be run inside K8s, so the service itself needs to be stateless. All the state information (ES: session information) must be held in some external place (Usually some DB or for files S3 on NFS).