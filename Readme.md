# Load Balancer

This package provides a load balancer implementation in Go. The load balancer distributes incoming requests across multiple servers in a round-robin fashion.

## Installation

To use this package, you need to have Go installed. Then, you can install the package by running the following command:

```shell
go get github.com/ritankarsaha/loadbalancer
```

## Usage

To use the load balancer, follow these steps:

1. Import the package:

```go
import "github.com/your-username/loadbalancer"
```

2. Create an array of servers:

```go
servers := []loadbalancer.Server{
    loadbalancer.NewSimpleServer("http://localhost:3000"),
    loadbalancer.NewSimpleServer("http://localhost:3001"),
    loadbalancer.NewSimpleServer("http://localhost:3002"),
}
```

3. Create a load balancer instance:

```go
lb := loadbalancer.NewLoadBalancer(":8080", servers)
```

4. Define a request handler function:

```go
handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
    lb.ServeProxy(rw, req)
}
```

5. Register the request handler:

```go
http.HandleFunc("/", handleRedirect)
```

6. Start the server:

```go
fmt.Printf("Serving requests at port listening on port %s\n", lb.Port)
http.ListenAndServe(lb.Port, nil)
```

## Documentation

### Server Interface

The `Server` interface represents a server that can handle incoming requests. It has the following methods:

- `Address() string`: Returns the address of the server.
- `IsAlive() bool`: Checks if the server is alive.
- `Serve(rw http.ResponseWriter, req *http.Request)`: Serves the incoming request.

```go
type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}
```

### SimpleServer

The `SimpleServer` struct is a concrete implementation of the `Server` interface. It represents a simple server that can handle HTTP requests. It has the following methods:

- `Address() string`: Returns the address of the server.
- `IsAlive() bool`: Checks if the server is alive.
- `Serve(rw http.ResponseWriter, req *http.Request)`: Serves the incoming request.
- 	addr: The server’s address (URL).
-  proxy: A reverse proxy that forwards requests to the actual server.

```go
type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}
```

### LoadBalancer

The `LoadBalancer` struct represents a load balancer. It has the following methods:

- `NewLoadBalancer(port string, servers []Server) *LoadBalancer`: Initializes a new load balancer instance.
- `ServeProxy(rw http.ResponseWriter, req *http.Request)`: Forwards the incoming request to the next available server.

	•	servers := []Server{...}: Creates a list of simpleServer instances, each representing a different website (Facebook, Bing, DuckDuckGo).
	•	NewLoadBalancer("8080", servers): Creates a new LoadBalancer that listens on port 8000.
	•	handleRedirect: Defines a function to forward any incoming requests to one of the servers via the load balancer.
	•	http.HandleFunc("/", handleRedirect): Registers the handleRedirect function to handle all incoming requests ("/").
	•	http.ListenAndServe(":"+lb.port, nil): Starts the HTTP server and listens for requests on localhost:8000.

### Round Robin Algorithm

The load balancer uses a round-robin algorithm to select the next available server. The `getNextAvailableServer()` method iterates through the servers in a round-robin fashion and returns the first server that is alive.


