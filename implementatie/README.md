
## Overview

- RESTful API endpoints for CRUD operations.
- JWT Authentication.
- Rate Limiting.
- Swagger Documentation.
- PostgreSQL database integration using GORM.
- Redis cache.
- MongoDB for logging storage.
- Dockerized application for easy setup and deployment.


## Starting the application
1. Navigate to the directory

```bash
cd implementatie
```

2. Build and run the Docker containers

```bash
make up
```

Please refer to the [Makefile](./Makefile) if you need to build in the local environment.


### Authentication

To use authenticated routes, you must include the `Authorization` header with the JWT token.

