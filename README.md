# Message Broker

# Introduction
In this project, I have implemented a message broker.

# Roadmap
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

