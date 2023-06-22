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
By Running `docker run -d -p 80:8080 -e SWAGGER_JSON_URL=https://gist.githubusercontent.com/AmmrFX/99d25fbfa8f9a16ad52aa5efca93c66f/raw/60c189fe5992875936d8c3a512d9dd7ee45f1d48/openapi-scuba.json swaggerapi/swagger-ui`
and accessing your server at `http://localhost:80/#/` you can see the API documentation as fellow 

![image](https://github.com/AmmrFX/scuba-divers/assets/55325468/01b3bad6-1888-43a6-b45b-2d844c861bfb)
![image](https://github.com/AmmrFX/scuba-divers/assets/55325468/5ca95718-f6f4-46ee-853c-b8faf5b9e621)
![image](https://github.com/AmmrFX/scuba-divers/assets/55325468/10fc4ff6-8480-469e-ab4c-02b2b7fff380)






## Docker Support
This project includes Docker support for containerization. You can use the provided Dockerfile files to build and run the application in a Docker container.
1. clone the repository
2. run `docker-compose up --build -d`
   ![image](https://github.com/AmmrFX/scuba-divers/assets/55325468/48249cf9-e81e-4cac-8d0d-3e597efe30eb)

4. to remove the containers alongside with the used volume run `docker-compose down -v`





