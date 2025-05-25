## Fluxend
Fluxend is a fast, no BS backend server built with Go. It cuts the noise and gives you raw power to build, scale, and own your backend — your way.

## Features
- Built-in Org & Role Management
- Instantly Generated Endpoints
- Plug-and-Play Auth
- Realtime Database
- Row-Level Access Control
- Import CSV/XLSX as APIs
- DB Functions, Triggers & Hooks
- Smart Forms with Validations & Triggers
- Multi-Driver Storage (S3, Dropbox, BackBlaze, FS)
- Detailed Audit Logs
- Built-in Search Engine (upcoming)
- Zapier Integration (upcoming)
- And much much more

## How it works
You can refer to [Wiki](https://github.com/fluxend/fluxend/wiki) to understand how different Fluxend components work and how they can be integrated into your existing stack. These explain basic functionality and detailed inner workings backed by flowcharts. Some of the topics include:
- [Dynamic REST endpoints](https://github.com/fluxend/fluxend/wiki/Dynamic-REST-Endpoints)
- [Authentication](https://github.com/fluxend/fluxend/wiki/Authentication)
- [Forms](https://github.com/fluxend/fluxend/wiki/Forms)
- [Storage](https://github.com/fluxend/fluxend/wiki/Storage)
- [Backup](https://github.com/fluxend/fluxend/wiki/Backups).

## Installation

### Method 1: Using Docker (Recommended for Easy Setup)
To get up and running with Fluxend in just a few minutes, simply follow these steps:

Clone the Fluxend repository:
```bash
git clone https://github.com/fluxend/fluxend.git fluxend
cd fluxend
make setup
   ```
This might take a while during first run. This will start two Docker containers:

- **Database Container**: A Postgres database to store your data.
- **Fluxend Server**: A backend server running on port 80.

Once the server is up, you can access the API documentation at http://localhost/docs/index.html.

### Method 2: Standalone Binary (For Self-Contained Deployments)
Prefer a single binary to run without Docker? You can easily build Fluxend and run it as a standalone executable:

Build with `make build` and then `./bin/fluxend` to start the server.

## Commands
Fluxend has several commands to perform operations and make your experience smoother. Fluxend binary supports core commands which is further augmented by make commands

### CLI commands
```
Fluxend CLI allows you to start the server, run seeders, and inspect routes.

Usage:
  fluxend [command]

Available Commands:
  help        Help about any command
  about       Prints information about the application
  optimize    Flush all caches and optimize the application
  routes      List all registered API routes
  seed        Seed the database with initial data
  server      Start the Fluxend API server
  udb.stats   Pull stats from given database
  udb.restart Restart all PostGREST instances

Flags:
  -h, --help   help for fluxend

Use "fluxend [command] --help" for more information about a command.
```

### Make commands
```
make help

help                           Shows this help
setup                          Setup the project
build                          Build the project
up                             Start the project
down                           Stop the project

login.app                      Login to fluxend container
login.db                       Login to database container

pgr.list                       List all postgrest containers
pgr.destroy                    Destroy all postgrest containers

docs.generate                  Generate docs
lint                           Run linter
lint.fix                       Run linter and fix

server                         Run the project in development mode
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
Fluxend is open-source! If you're passionate about building a blazing-fast backend and want to make Fluxend even better, we welcome contributions. Please send PRs our way.
