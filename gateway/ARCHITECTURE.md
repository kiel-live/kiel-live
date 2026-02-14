# Gateway Architecture

## Requirements

- viewport-based filtering
- search backend (not frontend anymore)
- on-demand data fetching
- realtime updates via websocket
- support aggregation (bus + train + bike as one stop)
- prevent data staleness
- request deduplication (exp only query stop details once and broadcast to all clients)
- cheap to host, performant
- KISS - simple, few failure points

## Current Limitations (NATS)

- no spatial filtering => all clients get ALL data
- search in frontend => should be backend
- hard to aggregate data from multiple sources

## New Gateway Architecture

### Core Concept

```
Collectors --HTTP POST--> Gateway <--WebSocket--> Clients (Web/PWA)
                          |
                          +-- In-memory state
                          +-- S2 spatial index
                          +-- Search index
```

### Components

Gateway:
- HTTP endpoints for data ingestion
- WebSocket for clients (viewport subs, stop details, vehicle details, search)
- WebSocket for collectors (on-demand fetch requests)
- In-memory storage with S2 spatial index
- In-memory search

### Data Flow

1. Updates: Collector → HTTP POST → Gateway → Broadcast to relevant clients
2. On-demand: Client subscribes → Gateway checks staleness → WS request to collector → Collector fetches → POST to gateway → Broadcast
3. Search: Client WS → Gateway searches in-memory → Returns results

### Protocol

#### Client => Gateway (WebSocket)
```json
{"type": "subscribe", "topic": "map.stops.<s2-cell-id>" }
{"type": "unsubscribe", "topic": "map.stops.<s2-cell-id>" }
{"type": "search", "query": "hauptbahnhof"}
```

#### Gateway => Client (WebSocket)

```json
{"type": "update", "id": "vehicle:123", "data": {...}}
{"type": "delete", "entity": "vehicle:123" }
{"type": "search_results", "data": [...]}
```

#### Collector => Gateway (HTTP)

```bash
POST /vehicles/{id}
DELETE /vehicles/{id}
```

#### Gateway => Collector (WebSocket)

```json
{"type": "fetch", "entity": "stop", "id": "..."}
```

#### Collector => Gateway (after fetch)

```bash
POST /stops/{id}
```

### Tech Stack

- Go
- gorilla/websocket
- golang/geo/s2
- golang.org/x/text
- no external dependencies

### Future

- redis for persistence?
- proxy to regional gateways
- aggregate data from multiple collectors
- metrics/monitoring
- rate limiting
