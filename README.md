# Location Share Backend
This is the backend for the location share app, that allows users to share their location with their friends. The backend is built using go and the gin framework.

## Development

### Running the backend

To run the backend, you need to have [go](https://golang.org/doc/install) installed on your machine.
You also need a postgres database running on your machine. You can install postgres from [here](https://www.postgresql.org/download/) or use the docker-compose file in the root directory to run a postgres database in a docker container for development purposes.

After installing go and postgres, you need to create a `app.env` file in the root directory of the project with the contents from the `app.env.example` file.

Then you can run the following command to start the backend server:

```bash
go run main.go
```

### Making changes to the database

After making changes to the existing models, run the following command to update the database schema:

```bash
go run migrate/migrate.go
```

When adding new models, you need to add them to the `migrate/migrate.go` file to be able to update the database schema.