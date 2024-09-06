# alertmanager-notification-router

This is a Go application that receives webhooks from [Alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/) and routes them to the [alert-hub-mvp](https://github.com/zinrai/alert-hub-mvp) system. It acts as a bridge between Alertmanager and the alert management system, transforming and forwarding alerts as needed.

## Features

- Receives Alertmanager webhooks
    - https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
- Processes and transforms alert data
- Routes alerts to the [alert-hub-mvp](https://github.com/zinrai/alert-hub-mvp) system

## Installation

```
$ go build -o alertmanager-notification-router ./cmd/server
```

## Usage

Run the application:

```
./alertmanager-notification-router
```

The server will start and listen on port 8080 by default.

## Using the /ahm Endpoint with curl

To send a webhook to the `/ahm` endpoint using curl, use the following command:

```bash
curl -X POST http://localhost:8080/ahm \
     -H "Content-Type: application/json" \
     -d @path/to/your/alertmanager_webhook.json
```

Replace `path/to/your/alertmanager_webhook.json` with the actual path to your JSON file containing the Alertmanager webhook payload.

Example of an Alertmanager webhook JSON:

```json
{
  "version": "4",
  "groupKey": "{}:{alertname=\"high_cpu_usage\"}",
  "status": "firing",
  "receiver": "team-X-pager",
  "groupLabels": {
    "alertname": "high_cpu_usage"
  },
  "commonLabels": {
    "alertname": "high_cpu_usage",
    "severity": "warning",
    "instance": "server01.example.com"
  },
  "commonAnnotations": {
    "summary": "High CPU usage detected",
    "description": "CPU usage is above 75% on server01.example.com"
  },
  "externalURL": "http://alertmanager.example.com",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "high_cpu_usage",
        "severity": "warning",
        "instance": "server01.example.com"
      },
      "annotations": {
        "summary": "High CPU usage detected",
        "description": "CPU usage is above 75% on server01.example.com"
      },
      "startsAt": "2023-03-15T08:00:00Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://prometheus.example.com/graph?g0.expr=cpu_usage_percent+%3E+75",
      "fingerprint": "1234567890abcdef"
    }
  ]
}
```

## Domain-Driven Design Implementation

This project follows Domain-Driven Design (DDD) principles to organize and implement its functionality. Here's an overview of how the domains are separated and implemented:

1. **Domain Layer** (`internal/domain`):
   - Contains the core domain models (`Alert` and `AlertmanagerWebhook`).
   - Defines the essential business logic and rules.

2. **Use Case Layer** (`internal/usecase`):
   - Implements the application's use cases, orchestrating the flow of data to and from the domain objects.
   - Contains the `AlertUseCase` which processes the Alertmanager webhooks and transforms them into our domain's `Alert` model.

3. **Interface Layer** (`internal/interface/handler`):
   - Handles the HTTP requests and responses.
   - The `AHMHandler` is responsible for receiving the webhook, calling the appropriate use case, and returning the response.

4. **Infrastructure Layer** (`internal/infrastructure/repository`):
   - Implements the data persistence logic.
   - The `AlertRepository` is responsible for saving alerts to the Alert Hub MVP system.

5. **Application Layer** (`cmd/server`):
   - Contains the main application logic.
   - Wires up all the dependencies and starts the HTTP server.

This separation of concerns allows for:

- Clear boundaries between different parts of the system.
- Easier testing and maintenance.
- Flexibility to change implementation details without affecting the core domain logic.

The flow of data typically goes from the Interface layer, through the Use Case layer, interacting with the Domain layer, and finally to the Infrastructure layer for persistence.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
