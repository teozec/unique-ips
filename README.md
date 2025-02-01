# Project github.com/teozec/unique-ips

This microservice calculates the number of unique IPs in the los.
Logs are sent to the `/logs` endpoint, and the number of unique IPs is returned by the `/visitors` endpoint.
By default, it listens on port 5000. It can be changed using the `PORT` environment variable.

## Endpoints

### `/logs`

Accepts POST requests in the following format:
`{ "timestamp": "2020-06-24T15:27:00.123456Z", "ip": "83.150.59.250", "url": ... }`

Adds the posted IP address to logs.

### `/visitors`

Accepts GET requests. Responds with the number of unique logged IP addresses in the following format:
`{ "count": 5 }`


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
