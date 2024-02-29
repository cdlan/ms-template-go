# Migrations
To run migrations this library is used [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)

## Create a new migration
To create a new migration file, run this command from the project root directory:
```shell
migrate create -ext sql -dir internal/migrations/sql -seq create_<something>
```
> replace `<something>` with a descriptive name of what does the migration do