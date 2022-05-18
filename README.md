# dp-frontend-interactives-controller

Controller for handling interactive visualisation pages on ONS website

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                   | Description                                                                                                        |
|------------------------------|---------------------------|--------------------------------------------------------------------------------------------------------------------|
| BIND_ADDR                    | :27300                    | Host and port to bind to                                                                                           |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                        | Graceful shutdown timeout in seconds (`time.Duration` format)                                                      |
| HEALTHCHECK_INTERVAL         | 30s                       | Time between self-healthchecks (`time.Duration` format)                                                            |
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| SERVE_FROM_EMBEDDED_CONTENT  | false                     | To serve content from embedded, static FS (storage/localfs) - testing/development only                             |
| API_ROUTER_URL               | http://localhost:23200/v1 | URL of the dp-api-router                                                                                           |

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
