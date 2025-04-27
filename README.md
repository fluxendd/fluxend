# Fluxton
**Blazing-Fast, Futuristic Backend for the Modern Web – Deploy with Just One File!**

## Table of contents
- [What is Fluxton](#what-is-fluxton)
- [Installation](#installation)
- [Commands](#commands)
- [Contribute](#want-to-contribute)

## What is Fluxton?
Fluxton is a lightweight, high-performance backend server built with Go. It’s designed to simplify backend development without sacrificing flexibility or speed.

With Fluxton, you can connect a database and immediately get a dynamic, auto-generated RESTful API for all your tables—no manual routing or boilerplate code required. It also handles form validation and submissions out of the box, making it easier to build robust APIs and admin panels quickly.

Everything runs from a single binary — no complex setup or external dependencies. Ideal for prototyping, internal tools, or building production-ready backends with minimal overhead.

## Why Choose Fluxton?
- Fast by Design: Built with Go and Echo, Fluxton is optimized for performance and low-latency APIs out of the box.
- Simple Deployment: Just one binary file—no external dependencies or setup scripts. Ideal for quick prototyping or deploying to production.
- Auto-Generated REST API: Connect your database, and Fluxton instantly exposes RESTful endpoints based on your schema. No manual route definitions needed.
- Flexible Query Builder: Easily construct advanced queries without writing raw SQL.
- Built-in Database UI: Includes a minimal interface for managing tables and records directly, useful for internal tools or quick data edits.
- Integrated Form Handling: Validate and handle form submissions server-side with minimal configuration—no extra validation libraries or middlewares required.

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
  make:model  Creates a new model file in the models directory
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
