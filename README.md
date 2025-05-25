# ⚡️ Fluxend
**Fluxend** is a **blazing-fast, self-hosted BaaS** built with Go — no fluff, no bloat, no lock-in. Ship production-grade backends in minutes with full control over your data, logic, and storage. It's your backend. Done your way.

## 🚀 Why Fluxend?
Tired of Firebase’s handcuffs? Supabase too slow or limited? Fluxend doesn’t babysit you — it gives you raw backend firepower out of the box:

- ✅ Fully open-source
- 🧠 Built with Go for max performance
- 🔩 Dead-simple setup with Docker
- 🧱 Modular & extendable
- 🧨 Ready for production on Day 1

## 🔥 Features
| Feature | Description |
|--------|-------------|
| 🧑‍💼 Org & Role Management | Built-in multi-tenant support with fine-grained RBAC |
| 🔐 Plug-and-Play Auth | OAuth, JWT, Magic Links — pick your poison |
| 🔄 Realtime Database | Instant updates pushed to clients |
| 🔥 Dynamic REST APIs | Define tables, get CRUD endpoints automagically |
| 🧮 Row-Level Permissions | Control access down to the individual row |
| 📥 Import CSV/XLSX as APIs | Upload a file → Get a full API. Done. |
| ⚙️ DB Triggers & Hooks | Run server-side logic without extra services |
| 🧾 Smart Forms | Auto-generated forms with validations and logic |
| ☁️ Multi-Driver Storage | S3, Dropbox, Backblaze, or local FS — your call |
| 🔍 Built-in Search Engine *(soon)* | Typesense/Sphinx powered indexing and search |
| 🔁 Zapier Integration *(soon)* | Automate anything with Fluxend events |
| 📜 Audit Logs | Every action tracked. No black boxes. |


## ⚙️ Installation
Clone the Fluxend repository:
```bash
git clone https://github.com/fluxend/fluxend.git fluxend
cd fluxend
make setup
   ```
This might take a while during first run. Once setup is done, the following containers will spin up:

- 🐘 **Postgres** – stores your application data (`fluxend_db`)
- 🧠 **Fluxend API Server** – backend engine (`fluxend_app`)
- 🌐 **Fluxend Frontend** – admin panel (`fluxend_frontend`)
- 🚦 **Traefik** – reverse proxy for routing requests (`fluxend_traefik`)

Access the app via the `APP_URL` defined in your `.env` file. Swagger docs available at: `http://{APP_URL}/docs/index.html`

## 📚 Learn How It Works
You can refer to [Wiki](https://github.com/fluxend/fluxend/wiki) to understand how different Fluxend components work and how they can be integrated into your existing stack. These explain basic functionality and detailed inner workings backed by flowcharts. Some of the topics include:
- [Dynamic REST endpoints](https://github.com/fluxend/fluxend/wiki/Dynamic-REST-Endpoints)
- [Authentication](https://github.com/fluxend/fluxend/wiki/Authentication)
- [Forms](https://github.com/fluxend/fluxend/wiki/Forms)
- [Storage](https://github.com/fluxend/fluxend/wiki/Storage)
- [Backup](https://github.com/fluxend/fluxend/wiki/Backups).

## 🧠 Commands & CLI
Fluxend has several commands to perform operations and make your experience smoother. Fluxend binary supports core commands which is further augmented by make commands

### 🔧 CLI commands
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

### 🛠 Make commands
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

## 🤝 Want to Contribute?
We're building the most badass backend tool in the open. If you:

- Hate boilerplate
- Love Go
- Want to build the next-gen BaaS engine
- Want influence in an early-stage rocket

Then you're in the right place.

- 🛠 Check out [issues](https://github.com/fluxend/fluxend/issues)
- 📬 Drop a PR. We'll review it FAST.
