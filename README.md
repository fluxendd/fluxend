## Fluxton
Fluxton is a fast, no BS backend server built with Go. It cuts the noise and gives you raw power to build, scale, and own your backend — your way.

## Features
- Built-in Org & Role Management
- Instantly Generated Endpoints
- Plug-and-Play Auth
- Built-in Search Engine
- Realtime Database
- Zapier Integration
- Row-Level Access Control
- Import CSV/XLSX as APIs
- DB Functions, Triggers & Hooks
- Smart Forms with Validations & Triggers
- Multi-Driver Storage (S3, Dropbox, BackBlaze, FS)
- Detailed Audit Logs
- And much much more

## Flowcharts
You can refer to the `.flowcharts` directory to understand how Fluxton works and how it can be integrated into your existing stack. These provide insights into dynamic REST endpoints, authentication, forms, storage, and backup processes.

## Installation

### Method 1: Using Docker (Recommended for Easy Setup)
To get up and running with Fluxton in just a few minutes, simply follow these steps:

Clone the Fluxton repository:
```bash
git clone https://github.com/fluxton-io/fluxton.git fluxton
cd fluxton
make setup
   ```
This might take a while during first run. This will start two Docker containers:

- **Database Container**: A Postgres database to store your data.
- **Fluxton Server**: A backend server running on port 80.

Once the server is up, you can access the API documentation at http://localhost/docs/index.html.

### Method 2: Standalone Binary (For Self-Contained Deployments)
Prefer a single binary to run without Docker? You can easily build Fluxton and run it as a standalone executable:

Build with `make build` and then `./bin/fluxton` to start the server.

## Commands
Fluxton has several commands to perform operations and make your experience smoother. Fluxton binary supports core commands which is further augmented by make commands

### CLI commands
```
Fluxton CLI allows you to start the server, run seeders, and inspect routes.

Usage:
  fluxton [command]

Available Commands:
  help        Help about any command
  about       Prints information about the application
  optimize    Flush all caches and optimize the application
  routes      List all registered API routes
  seed        Seed the database with initial data
  server      Start the Fluxton API server
  udb.stats   Pull stats from given database
  udb.restart Restart all PostGREST instances

Flags:
  -h, --help   help for fluxton

Use "fluxton [command] --help" for more information about a command.
```

### Make commands
```
make help

help                           Shows this help
setup                          Setup the project
build                          Build the project
up                             Start the project
down                           Stop the project

login.app                      Login to fluxton container
login.db                       Login to database container

pgr.list                       List all postgrest containers
pgr.destroy                    Destroy all postgrest containers

docs.generate                  Generate docs
lint                           Run linter
lint.fix                       Run linter and fix

serve                          Run the project in development mode
seed                           Seed the database
about                          Show the project information
optimize                       Optimize the project
udb.stats                      Show the database stats
udb.restart                    Restart all PostGREST instances
routes.list                    Show all the available routes
drop.user.dbs                  Drop all user-created databases

migration.create               Create a new database migration
migration.up                   Run database migrations
migration.down                 Rollback database migrations
migration.status               Show the status of the database migrations
migration.reset                Rollback all migrations and run them again
migration.redo                 Rollback the last migration and run it again
migration.fresh                Rollback all migrations and run them again

seed.fresh                     Seed the database with fresh data
```

### Want to Contribute?
Fluxton is open-source! If you're passionate about building a blazing-fast backend and want to make Fluxton even better, we welcome contributions. Please send PRs our way.
