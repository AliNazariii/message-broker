# Message Broker

## Introduction

The message broker project aims to implement a message broker that can be used to manage and route messages between various clients. This project includes implementing the broker.Broker interface, using Postgresql and Cassandra for data persistence, adding basic logs and Prometheus metrics, implementing a gRPC API for the broker and main functionalities, creating Dockerfile and docker-compose files for deployment, and deploying the application on a remote machine and Kubernetes (K8).

## Roadmap

- [X] Implement `broker.Broker` interface and pass all tests
- [X] Use postgresql for persisting data
- [X] Use cassandra for persisting data
- [X] Add basic logs and prometheus metrics
  - Metrics for each RPCs:
    - `method_count` to show count of failed/successful RPC calls
    - `method_duration` for latency of each call, in 99, 95, 50 quantiles
    - `active_subscribers` to display total active subscriptions
  - Env metrics:
    - Metrics for application memory, cpu load, cpu utilization, GCs
- [X] Implement gRPC API for the broker and main functionalities
- [X] Create *dockerfile* and *docker-compose* files for deployment
- [X] Deploy app with the previous `docker-compose` on a remote machine
- [X] Deploy app on K8

This project is designed to be scalable, reliable, and highly available. The use of Postgresql and Cassandra for data persistence ensures that data is stored safely and can be retrieved quickly when needed. The addition of logs and Prometheus metrics provides insight into how the system is performing, allowing for easy identification of issues and bottlenecks. The implementation of a gRPC API provides a simple and efficient way to communicate with the broker, while the use of Docker and Kubernetes simplifies deployment and management of the application.

Overall, this project is a comprehensive implementation of a message broker that can be used in various applications and environments. It provides a reliable and scalable solution that can handle high loads of data and provide real-time insights into system performance.

## Contributing

Contributions to this project are always welcome! Please feel free to open issues or submit pull requests to help improve the
project.
