# SCUBA Divers Log Microservice

This project is a microservice-based application for maintaining SCUBA divers' diving logs. It provides a RESTful API to create diver profiles, log new dives, retrieve dive logs, generate dive reports, and retrieve divers' information.

## Features

- Create diver profiles and store them in the database.
- Log new dives with depth and timestamp, while performing necessary validations.
- Retrieve dive logs for a specific diver.
- Generate dive reports with the total number of dives per diver.
- Retrieve divers' information based on provided diver IDs.
- Concurrency support for handling concurrent connections using Goroutines.
- Database: MySQL with the use of the native database/sql package.
- Framework: Gin, a lightweight web framework for Go.

## Installation

1. Clone the repository: git clone `https://github.com/AmmrFX/scuba-divers.git`
2. Install the dependencies: `go mod download`
3. Set up the MySQL database:
- Create a MySQL database and update the database configuration in the `main.go` file (`connectToDB` function).
- Or you can run `bash prepare.sh` to use automate the installation of mysql - setting up the password ; table and columns
4. Build and run the application:
- `go build -o scuba-divers-app`
- `./scuba-divers-app`

5. The application will be accessible at `http://localhost:8080`.

## API Documentation

The API documentation for the microservice is available in the OpenAPI 3.0 format. You can find the specification in the `api/openapi-doc.yaml` file. You can import this specification into tools like Swagger UI or Postman for better visualization and testing of the API endpoints.

## Docker Support
This project includes Docker support for containerization. You can use the provided Dockerfile files to build and run the application in a Docker container.




