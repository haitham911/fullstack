<h1 align="center">Task</h1>
<p align="center">
<a href="https://codecov.io/gh/haitham911/fullstack">
  <img src="https://codecov.io/gh/haitham911/fullstack/branch/master/graph/badge.svg?token=oZzvTyKCYc"/>
</a>
<p align="center">
 <img src="![workflow](https://github.com/haitham911/fullstack/actions/workflows/test.yml/badge.svg)"/>
</p>
</p>
## Task
> **sample:**
>
> REST API should be implemented consuming and producing “application/json”
> Implement user model with username, password, deposit and role (seller or buyer) fields
> Implement product model with amountAvailable, cost, productName and sellerId fields
> Implement an authentication method (basic, oAuth, JWT or something else, thechoice is yours)
> All endpoints should be authenticated except the user endpoints
> Implement CRUD for user and product models
> Implement /buy endpoint(accepts productId,amount of products) 
> so users with a “buyer” role can buy products with the money they’ve deposited.

> **Evaluation criteria:**
>
> Language/Framework of choice best practises
> Edge cases covered
> Write tests for buy and one CRUD endpoint of your choice
> Add Readme file

These samples provide a starting point for how to integrate different services using a Compose file and to manage their deployment with Docker Compose.

> **Note:**
>
> The following samples are intended for use in local development environments such as project setups, tinkering with software stacks, etc. These samples must not be deployed in production environments.

<!--lint disable awesome-toc-->
## Samples applications with multiple integrated services
- [`Go  / PostgreSQL`] - Sample Go

<!--lint disable awesome-toc-->
## Getting started

These instructions will get you through the bootstrap phase of creating and
deploying samples of containerized applications with Docker Compose.

### Prerequisites

- Make sure that you have Docker and Docker Compose installed
  - Windows or macOS:
    [Install Docker Desktop](https://www.docker.com/get-started)
  - Linux: [Install Docker](https://www.docker.com/get-started) and then
    [Docker Compose](https://github.com/docker/compose)
- Download some or all of the samples from this repository.

### Running a sample

The root directory of each sample contains the `docker-compose.yaml` which
describes the configuration of service components. All samples can be run in
a local environment by going into the root directory of each one and executing:

```console
docker-compose up -d
```
To stop and remove all containers of the sample application run:

```console
docker-compose down
```
To hit endpoints of the run postman:

```console
import task.postman_collection
import task.postman_environment
```
To run test:

```console
go test -v -coverprofile=coverage.out ./...
```
