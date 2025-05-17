
![CI status](https://github.com/aalbacetef/txnotify/actions/workflows/ci.yml/badge.svg) ![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg) 

# txnotify

`txnotify` is an Ethereum transaction observer that tracks incoming and outgoing transactions for subscribed addresses using Ethereum's JSON-RPC API. It pushes real-time updates to clients via WebSockets, making it suitable for integration with notification systems or wallet UIs.

Notes:
 - the core type is the `Watcher` 
 - there is a CLI tool as well as a WebSockets server and client

## Features

- Polls Ethereum blocks for new transactions
- Subscribes to and monitors any number of addresses
- Caches fetched blocks and transactions in memory
- Pushes notifications to clients over WebSockets
- Only external dependency is gorilla/websockets
- Extensible interfaces for RPC, storage, and notifier logic

## Architecture

### Components

- **Watcher**: Core engine that polls new blocks, fetches transactions, and notifies clients.
- **RPC Client**: Low-level JSON-RPC interface for Ethereum endpoints.
- **Cache**: In-memory store for blocks, transactions, and processing state. **Can be easily extended to any data storage backend.**
- **Notifier**: Interface for pushing updates to subscribers (WebSockets implementation included).
- **CLI / Server**: Commands to run the observer as a server or test client.

### Data Flow

1. `Watcher` polls latest block using JSON-RPC.
2. If a new block exists, it is fetched and stored in cache.
3. Transactions from the block are filtered and matched against subscribed addresses.
4. For matching transactions, the notifier sends data to connected clients.

### Usage

We'll use two addresses as examples:
- USDT: 0xdAC17F958D2ee523a2206206994597C13D831ec7
- USDC: 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48

You can see them on etherscan.io

#### WebSockets

Run the WebSocket server:

```bash
go run ./cmd/server

## if you prefer pretty printed output 

go run ./cmd/server | jq 
```

Run a test client:

```bash
go run ./cmd/client --addresses 0xdAC17F958D2ee523a2206206994597C13D831ec7,0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48
```

Note these are the addresses for USDT and USDC.

#### CLI tool for watching txs
Run the block watcher:

```bash
go run ./cmd/watch --address 0xdAC17F958D2ee523a2206206994597C13D831ec7

## if you prefer pretty printed output 

go run ./cmd/watch --address 0xdAC17F958D2ee523a2206206994597C13D831ec7 | jq
```

## Configuration

The RPC endpoint and polling interval can be set directly in the `main` or passed as arguments/environment variables in a production setup.

## Testing

Unit tests for core logic are included:

```bash
make test
```

Tests cover:
- JSON-RPC implementation 
- Cache behavior 
- Hex utilities

The `cmd/watch`, `cmd/server`, and `cmd/client` binaries serve as integration tests.

## Roadmap

- Add persistent storage backend (e.g., file or Redis)
- Improve reprocessing of failed or incomplete blocks
- Move a lot of the processing to a separate job system/message queue
- Enrich notifications with detailed transaction data
- Extend CLI for manual queries (e.g., list txs for address)

## Known Limitations

- Could potentially retry blocks continuously
- In-memory cache only; restarts clear state
- No batching or pipelining of RPC calls beyond transaction fetch concurrency
- Clients will timeout if no messages come in 


