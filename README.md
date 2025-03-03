# Fluxton
**Blazing-Fast, Futuristic Backend for the Modern Web – Deploy with Just One File!**

## What is Fluxton?
**Fluxton** is a cutting-edge backend server that is as fast as it is simple to use. Built with the power of Go and the flexibility of the Echo framework, Fluxton allows you to create scalable and dynamic backend solutions with minimal effort. The best part? You get everything in a single file – no need for complex infrastructure setups or tedious server configurations.

Fluxton automatically handles backend logic, form validation, submissions, and even provides an API for your database. Simply connect your database, and Fluxton dynamically generates a fully-functional RESTful API based on the tables you create. It's truly backend development for the future.

## Why Choose Fluxton?
- **Blazing Fast**: Powered by Go and optimized for performance, Fluxton delivers lightning-fast responses.
- **Zero Hassle Deployment**: With just **one file**, you can deploy and be up and running in minutes.
- **Automatic API Generation**: Fluxton automatically generates RESTful API endpoints based on your database schema. No more manually creating routes!
- **Dynamic Query Builder**: Effortlessly build complex database queries with Fluxton's dynamic query builder.
- **Database Management Made Easy**: Fluxton provides an intuitive database editor for seamless table management.
- **Form Validation & Backend Management**: Handle form submissions and validate data without writing a single line of validation logic.

## Installation

Getting started with **Fluxton** is as easy as 1-2-3. Choose the installation method that suits you best:

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

### Want to Contribute?
Fluxton is open-source! If you're passionate about building a blazing-fast backend and want to make Fluxton even better, we welcome contributions. Please send PRs our way.