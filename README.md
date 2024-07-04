# Golang RESTful API Templete
This project is a Golang RESTful API template designed with the same simplicity and flexibility as the **Express.js** framework. It provides a robust starting point for building scalable and maintainable APIs in Go, with clear and concise patterns for routing, middleware, and request handling. Whether you are transitioning from JavaScript to Go or starting fresh with Go, this template helps streamline the development process and ensures best practices are followed, making it easier to build performant and reliable web services.

## Tech Stack
- Programming Language: [GO](https://go.dev/)
- HTTP Web Framework: [Gin-Gonic](https://gin-gonic.com/)
- Database: [PostgreSQL](https://www.postgresql.org/)
- Database Driver: [Gorm](https://gorm.io/)

## Feature
- Refresh - Access token auth system
- Middleware
- Rate Limitter
- Cors
- Strict body request
- Auto Migration
- Seeder
- Filter by date/string builder
- String Search builder
- Pagination builder
- Mailer

## Base User Routes
- [GET] /v1/users
- [GET] /v1/users/:uuid
- [POST] /v1/signup
- [PUT] /v1/users/:uuid
- [POST] /v1/signin
- [POST] /v1/signout
- [POST] /v1/check-email
- [POST] /v1/check-username
- [POST] /v1/reset-password
- [GET] /v1/check-reset-token/:token
- [PUT] /v1/reset-password

## Application Flow
The application flow is structured as follows:

- Model: Acts as the blueprint for the data structures, similar to TypeScript or DTOs, defining the shape of the data used throughout the application.
- Controller: Handles routing and binds request data from the body or parameters, directing incoming requests to the appropriate endpoints
- Service: Contains all the business logic, processes data received from the controllers, interacts with the models, and executes the core functionalities of the application.
    - Sub-service: Contains a simple function such as checking user, get a single field of data, etcm that mostly return a boolean or a string.
    - Filter-service: Contains a builder for the filter. Any logic to build a filter such as filter by username, filter by date, pagination, etc.

## Development Setup

You need atleast using Golang v1.22.0

### Database
To set up the database for this project, you have two options: using Docker or WSL for a Linux environment. Follow the steps below to initialize the database.

#### Option 1: Using Docker
1. Install docker
2. Run this command:
```bash
docker pull postgres
docker run --name my_database -e POSTGRES_PASSWORD=mysecretpassword -d postgres
docker ps
```

#### Option 2: Using WSL for Linux Environment
1. Install WSL2 for windows (If you already using Linux as your main OS, then you don't need to download WSL)
2. Run this command:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo service postgresql start
sudo service postgresql status
```

By following these steps, you can quickly set up the database for your Golang RESTful API project using either Docker or WSL.

### Project

#### Installation
To run this project, you need to install the required go package by run this command:
```bash
go mod tidy
```

#### env
For the environment file, simply copy `example.env` and rename it to `.env`.
- For the Server section, you can costumize the port and frontend host.
- For the Database section, you need to fill it up based on the database that you have initialize before.
- For the Services related section:
    - `JWT_SECRET_KEY`: you can generate random string for this.
    - `USER_ADMIN_EMAIL` & `USER_ADMIN_PASSWORD`: for the seeder as initial admin account.
    - `MAILER_EMAIL`: for the mailing system (login notification, reset password, etc.)
    - `MAILER_PASSWORD`: you can get your gmail App Password by follow this steps:
        - Create Gmail Account
        - Enable [2FA](https://myaccount.google.com/signinoptions/two-step-verification/enroll-welcome)
        - Create [App Password](https://myaccount.google.com/apppasswords)
    - `RESETPW_FE_ENDPOINT`: this is the endpoint from frontend for the reset password link that sent to user email.
- Lastly, for the Linux build section, this is only required when you want to build the app at linux OS. If you build using windows, you can comment this section.

#### Run the project
Doing all of the steps above, you have two options to run the app:
1. Use the [Makefile](https://gnuwin32.sourceforge.net/packages/make.htm) to start the application:
2. Run this command:
```bash
make
```
Or
```bash
make dev
```

Or simply run the Go command:
```bash
go run main.go
```

## Development Flow
The development flow is similar to Express.js. First, you need to declare the model -> route -> request binding -> database query -> and data return. This is the detailed development flow for this project:

1. Table Declaration: Define a main model representing the table to be created, using singular naming conventions.
2. Migration: Import the created model into app/config/database/migration.go and add it to the modelList variable.
3. Route: Declare routes in the controller with plural names corresponding to the entities, specifying the endpoints for various operations.
4. CRUD: For operations requiring request payloads, create specific models such as "UserCreate" for creating users, and implement corresponding functions like "Create User" in the service layer to handle business logic.