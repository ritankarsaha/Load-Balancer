package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

// handling error gracefully by the handleErr function by the user
func handleErr(err error) {
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

// creating an interface for the server
type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

// a concrete implementation of the server interface
type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

// simple server methods to be used by the user
func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	resp, err := http.Get(s.addr)
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// newSimpleServer(addr string): This function initializes a simpleServer.
// 	•	url.Parse(addr): Parses the server’s address (addr) into a URL structure.
// 	•	httputil.NewSingleHostReverseProxy(serverUrl): Creates a reverse proxy that forwards requests to the specified URL.
func newSimpleServer(addr string) *simpleServer {
	serveURL, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serveURL),
	}
}

// load balancer struct
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

// the basic idea of the getNextAvailableServer() function is to iterate through the servers in a round-robin fashion and return the first server that is alive.
func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

// round robin algorithm to select the next server by the load balancer
func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++

	return server
}

// forwarding the request to the next available server
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()
	log.Printf("proxying request to %s\n", targetServer.Address())
	targetServer.Serve(rw, req)
}

// main function to start the server at the respective URLs to be used by the user.
func main() {
	servers := []Server{
		newSimpleServer("http://localhost:3000"),
		newSimpleServer("http://localhost:3001"),
		newSimpleServer("http://localhost:3002"),
	}

	lb := NewLoadBalancer(":8080", servers)

	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Serving requests at port %s\n", lb.port)
		if err := http.ListenAndServe(lb.port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")

	// Perform any cleanup or additional shutdown logic here

	log.Println("Server gracefully stopped")
}
