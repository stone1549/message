# message
Example REST service for sending/receiving messages, written in Go using Chi.

## Configuration

### Environment Variables

##### MESSAGE_SERVICE_ENVIRONMENT

Controls log levels and configuration defaults.

* DEV
* PRE_PROD
* PROD

##### MESSAGE_SERVICE_REPO_TYPE

* IN_MEMORY
* POSTGRESQL
    * MESSAGE_SERVICE_PG_URL - Full connection string for PG

##### MESSAGE_SERVICE_TIMEOUT

Incoming request timeout value in seconds.

##### MESSAGE_SERVICE_PORT

Port to run service on.


## Run

```go run main.go```

