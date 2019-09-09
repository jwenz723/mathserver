# mathserver

This repo contains various implementations of a basic server built in go that can execute the following
operations:

* Divide
* Max
* Min
* Multiply
* Pow
* Subtract
* Sum

# Purpose

The purpose of the various implementations provided in this repository is to give an example of how/when different 
implementation styles should be followed. The implementations compared are go-kit, standard gRPC, and standard HTTP.
The goal is to simply compare coding style and readability. Performance and efficiency of execution will not be compared.

# Comparison

The implementations are as follows:

### Client Implementations:

* Supports both gRPC and HTTP
    * [gokit](/client/cmd/gokit):
        * Total number of lines written: 103
        * gRPC is achieved via go-kit gRPC client
        * HTTP is achieved via go-kit HTTP client
    * [grpc_and_http](/client/cmd/grpc_and_http):
        * Total number of lines written: 223
        * gRPC is achieved via standard [gRPC library]("google.golang.org/grpc")
        * HTTP is achieved via standard [net/http library](https://golang.org/pkg/net/http/)
* Supports gRPC only
    * [grpcnative](/client/cmd/grpcnative):
        * Total number of lines written: 96
        * gRPC is achieved via standard [gRPC library]("google.golang.org/grpc")
        
### Server Implementations:

* Supports both gRPC and HTTP
    * [gokit](/grpc_and_http/gokit):
        * Total number of lines written: 1046
        * gRPC is achieved via go-kit gRPC transport layer
        * HTTP is achieved via go-kit HTTP transport layer
        * Logging is achieved via middleware wrapping business logic (go-kit service) layer
        * Prometheus instrumentation is achieved via middleware wrapping the business logic (go-kit service) layer
    * [std](/grpc_and_http/std):
        * Total number of lines written: 483
        * gRPC is achieved via standard [gRPC library]("google.golang.org/grpc")
        * HTTP is achieved via standard [net/http library](https://golang.org/pkg/net/http/) plus [gorilla router]("github.com/gorilla/mux")
        * Logging is achieved via middleware wrapping business logic layer
        * Prometheus instrumentation is achieved via middleware wrapping business logic layer
* Supports gRPC only
    * [gokit](/grpc_only/gokit)
        * Total number of lines written: 791
        * gRPC is achieved via go-kit gRPC transport layer
        * Logging is achieved via middleware wrapping business logic (go-kit service) layer
        * Prometheus instrumentation is achieved via middleware wrapping the business logic (go-kit service) layer
    * [grpcnative](/grpc_only/grpcnative)
        * Total number of lines written: 292
        * gRPC is achieved via standard [gRPC library]("google.golang.org/grpc")
        * Logging is achieved via [zap grpc interceptor](https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/logging/zap)
        * Prometheus instrumentation is achieved via [prometheus grpc interceptor](https://github.com/grpc-ecosystem/go-grpc-prometheus)
    * [std](/grpc_only/std)
        * Total number of lines written: 371
        * gRPC is achieved via standard [gRPC library]("google.golang.org/grpc")
        * Logging is achieved via middleware wrapping business logic layer
        * Prometheus instrumentation is achieved via middleware wrapping business logic layer