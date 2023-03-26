# message
Example REST service for sending/receiving messages, written in Go using Chi.

## Prereqs
You probably want to checkout the monorepo [yapyapyap](https://www.github.com/stone1549/yapyapyap) instead

## Configuration

### Environment Variables


| Variable                    | Description                                    | Possible Values                            |
|-----------------------------|------------------------------------------------|--------------------------------------------|
| MESSAGE_SERVICE_ENVIRONMENT | Controls log levels and configuration defaults | DEV, PRE_PROD, PROD                        |
| MESSAGE_SERVICE_REPO_TYPE   | Sets the type of storage to be used            | IN_MEMORY, POSTGRESQL, AUTH_SERVICE_PG_URL |
| MESSAGE_SERVICE_TIMEOUT     | Incoming request timeout value in seconds      | number                                     |  
| MESSAGE_SERVICE_PORT        | Port to run service on                         | number                                     |
| MESSAGE_SERVICE_PG_URL      | Full connection string for PG                  | string                                     |
| MESSAGE_SERVICE_TOKEN_PRIV  | private key for signing jwt tokens             | string                                     |
| MESSAGE_SERVICE_TOKEN_PUB   | public key for signing jwt tokens              | string                                     |

## Run

```go run main.go```

