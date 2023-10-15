**What is this project?**
--------------

This is a basic chat backend project. It has only fundamental features, such as authentication and sending and receiving messages.

## **Features implemented into the backend:**

* Authentication.
* Send and receive message.

## **Architecture:**

* Clean arquitecture.

## **API:**

<!-- * REST API. -->
* GraphQL API.

## **Tech used for build the backend:**

* **Database:** PostgreSQL.
* **Cache:** Redis.
* **Pub/Sub:** Redis (to send and receive messages in real time).
* **Real time connection with the client:** GraphQL Subscription.

## **Deploy**

To deploy the backend you can use docker or kubernetes.

**Docker:**


Clone the github project of the chat

`git clone https://github.com/daniarmas/chat.git`


Run the docker compose for deploy the needed services

`docker-compose -f docker-compose-dev.yaml up -d`


Run the docker image with the migrations

`docker run --rm --network="chat_default" ghcr.io/daniarmas/chat_migrations:latest --url="jdbc:postgresql://postgres:5432/chat?currentSchema=public" --changelog-file="v1_changelog.sql" --username="postgres" --password="postgres" update --log-level error`

***Github project for the migrations. [Click](https://github.com/daniarmas/chat_migrations)***


Seed the database with users for test the chat

`go run main.go database seed`

*In cmd/database/seed are the users for seed, you can see the password here.*


And finally run the server

`go run main.go server run`


**Kubernetes:**


In the kubernetes directory there are the manifests to deploy in kubernetes, run this command with the terminal open in the directory.

`kubectl apply -f .`

## **Run database migrations**

`docker run --rm --network="chat_default" -v ./internal/sql:/liquibase/changelog/ liquibase/liquibase:4.11 --url="jdbc:postgresql://postgres:5432/chat?currentSchema=public" --changelog-file="V1.0__changelog.sql" --username="postgres" --password="postgres" update --log-level error`

***This is for testing and development purposes only. It is recommended to change the passwords and environment variables.***