# Blog API (Golang)

The **Blog API** is a high-performance, RESTful backend service built in **Go (Golang)**, designed for managing blog content and user interactions.  
It provides the core functionality for creating, reading, updating, and deleting blog posts, along with secure user authentication and role-based access control.

## Overview

This API is built with **clean architecture principles**, ensuring that the business logic is separated from the infrastructure layer.  
It’s lightweight, efficient, and scalable — ideal for production environments where performance and maintainability are key.

For detailed request/response formats and example calls, see the **[Postman Documentation](https://documenter.getpostman.com/view/46777269/2sB3BEnpzU)**.

## Features

-   **High Performance** — Built with Go’s concurrency and lightweight footprint for minimal latency.
-   **Clean Architecture** — Clear separation between domain, use cases, and infrastructure layers.
-   **Secure Authentication** — JWT-based authentication with password hashing for user safety.
-   **Role-Based Access Control** — Different permissions for admins, authors, and regular users.
-   **Robust Validation** — Enforced input validation and error handling across endpoints.
-   **Extensible** — Designed to support additional features such as comments, categories, or tags.
-   **Database Integration** — Optimized queries with MongoDB (can be adapted to other databases).

## Technology Stack

-   **Language**: Go (Golang)
-   **Framework**: [Gin](https://gin-gonic.com/) for fast HTTP routing
-   **Database**: MongoDB (NoSQL)
-   **Authentication**: JWT
-   **Testing**: Go testing suite with mocks
-   **Documentation**: Postman

## Use Cases

-   Blogging platforms
-   Content management systems (CMS)
-   Knowledge base or article publishing tools

## Architecture

The API follows a **Clean Architecture** approach:

-   **Domain Layer** — Core entities and business rules.
-   **Use Case Layer** — Application-specific business logic and workflows.
-   **Infrastructure Layer** — Database connections, external APIs, and other implementations.
-   **Interface Layer** — HTTP handlers, request/response mapping.

This structure makes the codebase easier to maintain, test, and extend.

---

For complete API request and response examples, visit the **[Postman API Documentation](https://documenter.getpostman.com/view/46777269/2sB3BEnpzU)**.
