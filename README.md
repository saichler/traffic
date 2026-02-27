# traffic

UDP/TCP traffic generator app & service for testing Kubernetes Network Policies.

## Overview

A Go-based tool that runs in two modes:

- **Service mode (`Start`)** — Binds UDP and TCP listeners, waits for incoming commands and requests.
- **Client mode (`Do`)** — Sends a command via UDP to a running service, which then generates UDP or TCP traffic to a target destination and reports results.

This allows deploying listener pods across a Kubernetes cluster (as a DaemonSet) and then issuing traffic generation commands to verify that network policies permit or block the expected traffic.

## Project Structure

```
traffic/
  go/
    generator/
      Generator.go             # Entry point
      cmd/                     # CLI command framework (reflection-based arg parsing)
      message/                 # Message protocol, handler dispatch, timeout
      tcp/                     # TCP HTTP client (worker pool) and server
      udp/                     # UDP socket implementation
    tests/                     # UDP and TCP tests (1 to 10K packets)
    test.sh                    # Run tests with coverage
  docker/                      # Dockerfile, build script, run scripts
  k8s/
    deploy.yaml                # DaemonSet + ClusterIP + NodePort services
```

## Usage

### CLI Arguments

Arguments are passed as `Name=Value` pairs after the command name.

### Start a Listener

```bash
./generator Start Udp_port=15000 Tcp_port=16000
```

Starts both a UDP listener on port 15000 and a TCP (HTTP) listener on port 16000.

### Generate Traffic

```bash
# Send 1000 UDP packets to a remote listener
./generator Do Udp_port=15000 Destination=10.0.0.5 Port=16000 Quantity=1000

# Send 1000 TCP (HTTP) requests to a remote listener
./generator Do Udp_port=15000 Tcp_port=15010 Destination=10.0.0.5 Port=16000 Quantity=1000
```

| Argument      | Description                                      |
|---------------|--------------------------------------------------|
| `Udp_port`    | Local UDP port of the service to send command to  |
| `Tcp_port`    | If set, sends TCP traffic instead of UDP          |
| `Destination` | Target host IP                                    |
| `Port`        | Target port                                       |
| `Quantity`    | Number of packets/requests to send                |
| `Timeout`     | Seconds to wait for replies (auto-calculated if 0)|

### Report Output

```
Total UDP Sent:1000 OK:1000 Err:0 Timeout:false Took:2 seconds
```

## Docker

### Build and Push

```bash
cd docker && ./build.sh
```

Cross-compiles for `linux/amd64`, builds the image `saichler/traffic-generator:latest`, and pushes it.

### Run a Listener Container

```bash
docker run -e CMD=Start -e Udp_port='15000' -e Tcp_port='16000' saichler/traffic-generator:latest
```

## Kubernetes Deployment

```bash
kubectl apply -f k8s/deploy.yaml
```

Deploys into namespace `traffic-01`:
- **DaemonSet** — runs a listener pod on every node
- **ClusterIP Service** — internal access on ports 15000 (UDP) and 16000 (TCP)
- **NodePort Service** — external access on ports 30000 (UDP) and 30001 (TCP)

## Tests

```bash
cd go && ./test.sh
```

Tests cover:
- Unknown/invalid command handling
- UDP: single packet, 1K, 10K packets, timeout scenarios
- TCP: single packet, 100, 1K, 10K packets, error scenarios (unreachable destination)

## License

See [LICENSE](LICENSE).
