# API Gateway

Manages all other microservices, and exposes an API for frontend clients to use.

## Goals of this project:
- Learn about microservices
- Learn about gRPC
- Be able to use this service along with other microservice (have an API Gateway, a frontend in React, and interact with a To-Dos service.

## To-do:

- [X] Expose a REST Server
- [X] Add a validator external package
- [X] Interact with other gRPC servers such as the Users gRPC, and the To-dos gRPC (microservices)
- [X] Validate authorization tokens (JWTs) for protected endpoints.
- [ ] Add data from to-dos gRPC
- [ ] Refactor repeated code in middleware.

## Overall project architecture

#### API Gateway
- REST API (returns in JSON)
- A client for all the gRPC microservices (reads all the gRPC services).
- Responsible for the interactions between all the microservices.

#### Users gRPC
- Manages the users in the architecture

#### To-dos gRPC
- Manages all the to-dos.
- Maybe will be used for note-taking as well.

#### Frontend
- In React, will only call and use data from the API Gateway.
