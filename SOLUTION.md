# Solution Steps

1. Create a .env file specifying consistent Postgres and API environment variables (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB).

2. Modify docker-compose.yml to use the .env file for both db and api services, pass all needed variables through environment, and ensure DB and API env variables match.

3. For the db service in docker-compose.yml, add a healthcheck using pg_isready that waits for the database to accept connections.

4. Set the depends_on for the api service so it only starts when the db service is healthy (service_healthy).

5. Update Dockerfile to follow best practices: perform a multi-stage build, copy in built binary, expose the API port, and use ENTRYPOINT rather than CMD.

6. In main.go, implement a retry mechanism for connecting to the Postgres DB (attempt to connect & ping DB several times, with delays, before giving up and exiting with an error).

7. Load all environment variables in Go (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, DB_HOST, DB_PORT) to build the DSN, and use PORT for server listen port.

8. In Go, add a /healthz endpoint that checks DB connectivity.

9. Ensure application startup does not proceed until database connection is established, thus aligning startup order between Docker and the Go application.

