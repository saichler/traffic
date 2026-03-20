# l8traffic

UDP/TCP traffic generator for testing Kubernetes Network Policies.

## Overview

A Go-based tool (`github.com/saichler/traffic`) that runs in two modes:

- **Service mode (`Start`)** — Binds UDP and TCP listeners, waits for incoming commands and echo requests.
- **Client mode (`Do`)** — Sends a command via UDP to a running service, which then generates UDP or TCP traffic to a target destination and reports results.

Deploy listener pods across a Kubernetes cluster as a DaemonSet, then issue traffic generation commands from any pod to verify that network policies permit or block the expected traffic.

## How It Works

1. **Service pods** run on every node via a DaemonSet, each listening on UDP and TCP ports.
2. A **client** sends a command (via UDP) to any service pod, specifying a target destination, port, protocol, and packet count.
3. The service pod **generates traffic** to the target:
   - **UDP mode** — Allocates ephemeral UDP sockets, sends request messages, and tracks responses. Throttles at 5K packets (250ms pause) to avoid port exhaustion.
   - **TCP mode** — Uses a pool of 100 concurrent HTTP workers, each opening a TCP connection, sending a GET request, and reading the response.
4. The service pod **reports results** back to the client via UDP with a summary: sent count, OK count, error count, timeout status, and elapsed time.

## Project Structure

```
l8traffic/
  go/
    generator/
      Generator.go             # Entry point (main)
      cmd/                     # CLI command framework (reflection-based arg parsing)
        Command.go             # Argument parser (Name=Value pairs via reflection)
        Commands.go            # Command registry (Start, Do)
        Start.go               # Service listener mode
        Do.go                  # Traffic generation client mode
      message/                 # Message protocol, handler dispatch, timeout
        Message.go             # Wire format: pipe-delimited (id|action|dest|port|qty|text|timeout)
        MessageProtocol.go     # Action dispatch (sendUdp, sendTcp, request, response)
        Timeout.go             # Timeout management with condition variables
      tcp/                     # TCP implementation
        server.go              # HTTP server (echo listener)
        client.go              # HTTP client with 100-worker pool
      udp/                     # UDP implementation
        UDP.go                 # UDP socket handler
    tests/                     # Integration tests (UDP and TCP, 1 to 10K packets)
    test.sh                    # Run tests with coverage
  docker/                      # Dockerfile (Alpine), build script, entrypoint
  k8s/
    deploy.yaml                # Namespace + DaemonSet + ClusterIP + NodePort services
```

## Usage

### CLI Arguments

Arguments are passed as `Name=Value` pairs after the command name.

### Start a Listener

```bash
./generator Start Udp_port=15000 Tcp_port=16000
```

Starts both a UDP listener on port 15000 and a TCP (HTTP) echo server on port 16000.

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
OK Sample:
 - 127.0.0.1:15001
Err Sample:
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
- **DaemonSet** — runs a listener pod on every node (UDP 15000, TCP 16000)
- **ClusterIP Service** — internal access on ports 15000 (UDP) and 16000 (TCP)
- **NodePort Service** — external access on ports 30000 (UDP) and 30001 (TCP)

### Testing Network Policies

Once deployed, exec into any pod in the cluster and use the client mode to test connectivity:

```bash
# From any pod, test if traffic reaches the traffic-01 service
./generator Do Udp_port=15000 Destination=traffic-01-01.traffic-01.svc.cluster.local Port=15000 Quantity=10

# Expected output if NetworkPolicy ALLOWS traffic:
# Total UDP Sent:10 OK:10 Err:0 Timeout:false Took:0 seconds

# Expected output if NetworkPolicy BLOCKS traffic:
# Total UDP Sent:10 OK:0 Err:10 Timeout:true Took:5 seconds
```

## Tests

```bash
cd go && ./test.sh
```

Fetches dependencies, runs tests with coverage, and opens the coverage report. Tests cover:
- Argument validation (invalid format, unknown args, missing required args, port bounds)
- Command help output
- Unknown command handling
- UDP: single packet, 1K, 10K packets, timeout scenarios
- TCP: single packet, 100, 1K, 10K packets, error scenarios (unreachable destination)

## License

See [LICENSE](LICENSE).
