# Go Circuit Breaker

A Go implementation of the Circuit Breaker pattern using the `gobreaker` library. This project demonstrates how to implement fault tolerance and resilience in microservices by protecting against cascading failures.

## 🎯 What is a Circuit Breaker?

A Circuit Breaker is a design pattern used in software architecture to prevent cascading failures. It works like an electrical circuit breaker:

- **Closed State**: Normal operation, requests pass through
- **Open State**: Circuit is open, requests are blocked immediately
- **Half-Open State**: Limited requests are allowed to test if the service has recovered

## 🚀 Quick Start

### Prerequisites

- Go 1.24.2 or higher
- Git

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd go-circuit-breaker
```

2. Install dependencies:

```bash
go mod download
```

### Running the Application

1. **Start the unstable API server** (in one terminal):

```bash
make unstable-api
# or
go run ./unstable/main.go
```

2. **Start the circuit breaker client** (in another terminal):

```bash
make circuit
# or
go run ./breaker/main.go
```

## 🔧 Configuration

The circuit breaker is configured in `breaker/main.go`:

```go
settings := gobreaker.Settings{
    Name:        "HTTP API",
    MaxRequests: 3,                    // Max requests in half-open state
    Interval:    0,                    // Reset interval for failure counts
    Timeout:     5 * time.Second,      // Time circuit stays open
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures >= 3  // Failures needed to open
    },
    OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
        // Log state changes
    },
}
```

### Configuration Parameters

- **Name**: Identifier for the circuit breaker
- **MaxRequests**: Maximum requests allowed in half-open state
- **Interval**: Time window for counting failures (0 = no window)
- **Timeout**: How long the circuit stays open before testing recovery
- **ReadyToTrip**: Function that determines when to open the circuit
- **OnStateChange**: Callback for state transitions

## 🧪 Testing the Circuit Breaker

The unstable API server (`unstable/main.go`) randomly fails with 10% probability:

```go
if rand.Float32() < 0.1 {
    // Return error response
    w.WriteHeader(500)
    response = pkg.Response{
        Message: "Error on request :(",
        Status:  500,
    }
} else {
    // Return success response
    w.WriteHeader(201)
    response = pkg.Response{
        Message: "Success",
        Status:  201,
    }
}
```

### Expected Behavior

1. **Normal Operation**: Requests succeed and fail randomly
2. **Circuit Opens**: After 3 consecutive failures, circuit opens
3. **Fast Failures**: When open, requests fail immediately without calling the API
4. **Recovery Test**: After 5 seconds, circuit goes to half-open state
5. **Recovery**: If requests succeed in half-open state, circuit closes

## 📊 Monitoring

The application logs:

- Successful requests with status codes
- Failed requests with error messages
- Circuit breaker state changes
- API response details

Example output:

```
2025/06/26 19:36:58 Resposta: Success - 201
2025/06/26 19:36:58 Erro na requisição: Error on request :(
2025/06/26 19:36:58 Circuit breaker mudou de Closed para Open
2025/06/26 19:36:58 Erro na requisição: circuit breaker is open
```

## 🏗️ Architecture

### Components

1. **Circuit Breaker Client** (`breaker/main.go`)

   - Implements the circuit breaker pattern
   - Makes HTTP requests to the unstable API
   - Handles failures and state transitions

2. **Unstable API Server** (`unstable/main.go`)

   - Simple HTTP server that randomly fails
   - Returns JSON responses with status codes
   - Simulates real-world service failures

3. **Shared Types** (`pkg/types.go`)
   - Common data structures used across the application
   - Response format for API communication

### Data Flow

```
Client Request → Circuit Breaker → HTTP Client → Unstable API
                ↓
            State Machine
                ↓
            Response/Error → Client
```

## 🔍 Circuit Breaker States

### Closed State

- **Behavior**: Normal operation
- **Requests**: Pass through to the service
- **Transitions**: Opens when failure threshold is reached

### Open State

- **Behavior**: Fast failure mode
- **Requests**: Immediately rejected
- **Transitions**: Goes to half-open after timeout

### Half-Open State

- **Behavior**: Limited testing
- **Requests**: Limited number allowed through
- **Transitions**: Closes on success, opens on failure

## 🛠️ Customization

### Adjusting Failure Threshold

To make the circuit breaker more or less sensitive:

```go
ReadyToTrip: func(counts gobreaker.Counts) bool {
    return counts.ConsecutiveFailures >= 5  // More failures needed
}
```

### Changing Timeout Duration

To adjust how long the circuit stays open:

```go
Timeout: 10 * time.Second  // Longer timeout
```

## 🔗 References

- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
- [gobreaker Library](https://github.com/sony/gobreaker)
- [Go Programming Language](https://golang.org/)
