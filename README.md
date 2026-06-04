# Mosaic

Mosaic is a highly available, fault-tolerant distributed key-value store written purely in Go. It is designed to prioritize strict consistency and durability under heavy concurrent workloads, utilizing a custom implementation of the Raft consensus algorithm over a high-performance gRPC network layer.

## Architecture Overview

Mosaic operates as a coordinated 5-node cluster. It avoids the eventual consistency models of Dynamo-style databases, opting instead for state machine replication via Raft to ensure strict linearizability across all cluster operations.

### Core Components
* **Consensus Engine (Raft):** Handles leader election, log replication, and dynamic membership. Ensures that data remains consistent and the cluster remains available even in the event of partial network partitions or up to 2 simultaneous node failures.
* **Storage Layer:** A heavily concurrent, thread-safe in-memory map (`sync.RWMutex`) designed to prevent writer starvation, backed by a persistent Write-Ahead Log (WAL) for absolute data durability.
* **Transport Layer:** Internal node-to-node communication and external client gateways are powered entirely by gRPC and Protocol Buffers, drastically reducing serialization overhead and network bandwidth compared to standard REST/JSON architectures.

## Performance & Guarantees

* **Throughput:** Sustains ~12,000 read/write operations per second (ops/sec).
* **Latency:** Maintains a p99 latency of <40ms during benchmark load testing.
* **Linearizable Reads:** Implements Raft ReadIndex to validate leadership before serving reads, guaranteeing that clients never receive stale data, even under injected partition faults or leader-crash scenarios.
* **Bounded Crash-Recovery:** State restoration is bounded via an append-only WAL paired with a 10k-entry snapshot compaction mechanism, preventing infinite log growth and ensuring rapid node reboot cycles.

## Technical Implementation Details

* **Language:** Go 1.21+
* **Concurrency:** Extensive use of goroutines and select-statement multiplexing over channels to handle concurrent network RPCs and heartbeat timeouts without lock contention.
* **Fault Tolerance:** Validated against injected network partitions, message reordering, and arbitrary process crashes with zero observed stale reads or log divergence upon cluster recovery.

## Getting Started

### Prerequisites
* Go 1.21 or higher
* Protocol Buffer Compiler (`protoc`)
* gRPC Go plugins

### Installation
```bash
git clone [https://github.com/jelb30/Mosaic.git](https://github.com/jelb30/Mosaic.git)
cd mosaic
go mod tidy
