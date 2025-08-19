# chirpy

A Twitter-like REST API and backend created for "Learn HTTP Servers in Go" [boot.dev course](https://www.boot.dev/courses/learn-http-servers-golang).

## Tools

- postgres database
- [goose](https://github.com/pressly/goose) used for db migrations
- [sqlc](https://sqlc.dev/) used for compiling sql queries to type-safe Go code
- [golang-jwt](https://github.com/golang-jwt/jwt) used for authentication

## Project Structure

### main.go

main.go initialises the environment by:

- loading environment variables from .env
  - database connection string
  - platform
    - set to "dev" to enable reset endpoints
  - JWT secret
  - "POLKA_KEY"
    - API key for imaginary 3rd party service sending webhooks
  - goose migration config
    - set GOOSE_DRIVER="postgres"
    - set GOOSE_DBSTRING=\<YOUR DB CONNECTION STRING\>
    - set GOOSE_MIGRATION_DIR="sql/schema"
- creating an instance of the api config struct
- initializing database connection

It then initialises a NewServeMux() and implements handlers for the various supported endpoints.

Finally, it creates the http server and begins listening and serving on port 8080.

### /sql/

- sql/schema/
  - this is where the goose migration files live
- sql/queries/
  - this is where the sqlc query files live
  - if changes are made here, run "sqlc generate" from the project root to generate the associated Golang code, which populates in /internal/database

### /internal/

#### /internal/database/

This is where the sqlc-generated Go models and functions are created.

#### /internal/auth/

Comprises the "auth" package.

Authentication-related files, including:

- JWT token creation and validation
- password hashing and checking
- API key helper functions

#### /internal/config/

Comprises the "config" package.

The endpoint handler functions are defined here, each function getting its own file.
Note that all handler functions are defined as methods on the apiConfig struct defined in config.go.

There are also some helper functions; e.g. "helper_authenticateUser.go".

"helper_dataSerialization.go" is a crucial file that helps define the way JSON responses are marshalled, managing how the sqlc-generated structs are given json struct tags and censored for public API consumption (e.g. removing hashed password fields).

### Other

The root dir has miscellaneous other files.

- sqlc.yaml configures sqlc
  - sql driver
  - where to look for query files
  - where to generate models and functions
- Makefile
  - makes it very easy to run the server; simply execute "make run" from the project root
- index.html
  - placeholder web page served on the /app/ path
- /assets/
  - placeholder dir for public assets; not used
