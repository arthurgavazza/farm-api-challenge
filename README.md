# Traive Challenge - Backend Developer

## Overview

This repository contains a solution to the backend developer challenge. The project implements a JSON API to manage a collection of **Farms** and their associated **Crop Productions**. The solution demonstrates backend development best practices, including modular code structure, environment configuration, and optional stretch goal implementations.

## Features

- **Create a Farm** with nested Crop Productions.
- **Delete a Farm** by its ID.
- **List all Farms** with pagination and filtering.

## Technologies Used

- **Programming Language**: Golang
- **Frameworks/Libraries**: [Mention relevant libraries]
- **Database**: PostgreSQL
- **Containerization**: Docker
- **Testing Tools**: [e.g., Go testing, Jest, pytest, etc.]

## Project Structure
```
├── Dockerfile
├── README.md
├── cmd
│   └── main.go
├── docker-compose.dev.yml
├── docker-compose.local.yml
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   └── app
│       ├── domain
│       │   ├── crop_production.go
│       │   ├── farm.go
│       │   ├── farm_repository.go
│       │   └── usecases
│       │       ├── create_farm.go
│       │       ├── create_farm_test.go
│       │       ├── delete_farm.go
│       │       ├── list_farms.go
│       │       └── module.go
│       ├── dto
│       │   └── create_farm_dto.go
│       ├── infra
│       │   ├── config
│       │   │   ├── config.go
│       │   │   └── module.go
│       │   ├── database
│       │   │   ├── database.go
│       │   │   ├── entities
│       │   │   │   ├── crop_production_entity.go
│       │   │   │   └── farm_entity.go
│       │   │   ├── mappers
│       │   │   │   ├── mappers.go
│       │   │   │   └── mappers_test.go
│       │   │   ├── module.go
│       │   │   └── repositories
│       │   │       ├── farm_repository.go
│       │   │       ├── farm_repository_test.go
│       │   │       └── module.go
│       │   └── httpapi
│       │       ├── controllers 
│       │       │   ├── farm_controller.go
│       │       │   ├── farm_controller_test.go
│       │       │   └── module.go
│       │       ├── module.go
│       │       ├── routers     
│       │       │   ├── farm.go
│       │       │   ├── module.go
│       │       │   └── router.go
│       │       └── server.go    
│       ├── models
│       │   └── models.go
│       └── shared
│           ├── errors
│           │   └── errors.go
│           ├── utils
│           └── validation
│               └── validation.go
└── testutils
    ├── fakes.go
    ├── helpers.go
    └── matchers.go
```
## Top-Level Directory Structure

### `cmd`
Contains the entry point of the application. The `main.go` file initializes and runs the API, pulling together configurations, dependencies, and modules.

### `internal`
Houses the core logic, domain models, infrastructure code, and shared utilities for the application. It is structured to follow a modular architecture, separating concerns between domain, infrastructure, and HTTP API layers.

#### `internal/app/domain`
Defines the domain layer of the application, including:
- **Business Logic (`usecases`)**: Implements core use cases like `create_farm` and `list_farms`.
- **Domain Entities**: Represents business concepts (e.g., `farm.go`, `crop_production.go`).
- **Interfaces**: Includes repository interfaces (e.g., `farm_repository.go`).

#### `internal/app/infra`
Implements infrastructure concerns such as configuration, database interactions, and HTTP APIs. It bridges the domain layer and external systems.

- **`config`**: Handles application configuration (e.g., environment variables and settings).  
- **`database`**: 
  - Manages database connections and schema definitions (entities).  
  - Includes mappers for converting between database models and domain models.  
  - Contains repository implementations.  
- **`httpapi`**:  
  - **Controllers**: API route handlers (e.g., `farm_controller.go`).  
  - **Routers**: Defines API routes and sets up routing logic.  
  - **Server**: Handles HTTP server setup.

#### `internal/app/dto`
Defines Data Transfer Objects (DTOs) used for API requests and responses. These objects provide a clear contract for data structures exchanged between the API and its consumers.

#### `internal/app/models`
Contains shared data structures and models used across the application. These models represent business objects and are independent of infrastructure concerns.

#### `internal/app/shared`
Includes shared utilities and components used across the application:
- **`errors`**: Custom error types and error-handling logic.  
- **`utils`**: General-purpose utility functions.  
- **`validation`**: Input validation logic for various application components.

### `testutils`
Provides helper functions, fake objects, and custom matchers to facilitate writing and organizing tests across the application.

---

## Top-Level Files

### `Dockerfile`
Defines the instructions to build a Docker image for the application.

### `docker-compose.*.yml`
Contains configurations for running the application in different environments (`dev`, `local`, etc.) using Docker Compose.

### `README.md`
Provides an overview of the project, including setup instructions, usage examples, and other relevant details.

### `go.mod` & `go.sum`
Manage dependencies and module requirements for the Go application.



## API Endpoints

The API includes the following endpoints:

### **Farm Endpoints**

#### Create a Farm

- **URL**: `/farms`
- **Method**: `POST`
- **Payload**:
  ```json
  {
    "name": "test2",
    "land_area": 550.5,
    "unit_measure": "hectares",
    "address": "123 Farm Lane, Countryside",
    "crop_productions": [
      {
        "crop_type": "COFFEE",
        "is_irrigated": true,
        "is_insured": false
      },
      {
        "crop_type": "CORN",
        "is_irrigated": false,
        "is_insured": true
      }
    ]
  }

  ```
- **Response**: Returns the created farm object.

#### Delete a Farm

- **URL**: `/farms/{id}`
- **Method**: `DELETE`
- **Response**: Confirmation of deletion.

#### List Farms

- **URL**: `/farms`
- **Method**: `GET`
- **Query Parameters** (optional):
  - `crop_type` (filter by crop type)
  - `minimum_land_area` (filter farms with land area greater than or equal to this value)
  - `maximum_land_area` (filter farms with land area less than or equal to this value)
  - `page` (pagination page number)
  - `per_page` (number of records per page)
- **Response**: 
  ```json
  {
    "items": [
        {
            "id": "264e0463-0d15-410b-9bc5-17e5e0741519",
            "name": "Sunny Farm",
            "land_area": 120.5,
            "unit_measure": "hectares",
            "address": "123 Farm Lane, Countryside",
            "created_at": "2024-12-09T22:07:44.357163-03:00",
            "updated_at": "2024-12-09T22:07:44.357163-03:00",
            "crop_productions": [
                {
                    "id": "05ef1eac-a763-4f1e-9556-1010d9f0c879",
                    "farm_id": "264e0463-0d15-410b-9bc5-17e5e0741519",
                    "crop_type": "CORN",
                    "is_irrigated": false,
                    "is_insured": true
                },
                {
                    "id": "bcf21d06-8fd6-4eea-b347-21f4d28fc7e1",
                    "farm_id": "264e0463-0d15-410b-9bc5-17e5e0741519",
                    "crop_type": "COFFEE",
                    "is_irrigated": true,
                    "is_insured": false
                }
            ]
        }
    ],
    "total_count": 4,
    "current_page": 1,
    "per_page": 2
}
  ```

## Setup Instructions

### Prerequisites

- Go 1.23+
- Docker

### Installation Steps

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-name>
   ```
2. Install dependencies:
   ```bash
   [Insert command to install dependencies, e.g., `go mod tidy`]
   ```
3. Set up the environment variables:
   - Copy `.env.example` to `.env` and configure the values.
     ```bash
     cp .env.example .env
     ```
4. Run the application:
   ```bash
   [Insert command to start the server, e.g., `go run main.go`]
   ```
5. Access the API at `http://localhost:PORT`.

### Testing

Run the test suite with:

```bash
go test ./...
```
## Stretch Goals

The following stretch goals were implemented:

- **Pagination and Filtering** for the List Farms endpoint.
- **Unit and Integration Tests** to validate functionality.
- **Containerization** using Docker.
- **API Documentation** with OpenAPI.
- **Production Readiness** with input validation, logging, and error handling.

## Acknowledgments

Thank you for reviewing this solution! Please feel free to reach out with any feedback or questions.
