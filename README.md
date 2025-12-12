# Mini Prometheus Clone (Golang)

A lightweight, educational implementation of a monitoring system inspired by Prometheus, written in Go.

## What is this project about?
This project is a minimalist attempt to recreate the core functionality of **Prometheus**, the industry-standard monitoring tool. It strips away the complexity of distributed storage, complex query languages (PromQL), and alerting to focus on the fundamental feedback loop of monitoring: **Scrape, Store, and Serve**.

## Why was it made?
I built this project to look under the hood of time-series databases and monitoring systems. I wanted to answer questions like:
- How does a monitoring system collect data from hundreds of services?
- How do you handle reading and writing data at the same time?
- How can you build a simple yet effective "Time Series Database" (TSDB) in memory?

## What I Learnt
Building this project taught me several key concepts in **Software Engineering** and **System Design**:
1.  **Concurrency in Go**: Using `Goroutines` to run multiple scrapers in parallel without blocking the main server.
2.  **Thread Safety**: Implementing `sync.RWMutex` to ensure safe concurrent access to the in-memory storage (solving the "Readers-Writers" problem).
3.  **Modular Architecture**: Separating concerns into distinct packages (`api`, `scraper`, `storage`, `config`) for maintainability.
4.  **Resource Management**: The importance of **Data Retention**. Without the pruning mechanism I added, the application would eventually crash from running out of RAM.
5.  **HTTP Networking**: Building robust clients that can tolerate failures when scraping external targets.

## What does it do?
1.  **Scrapes Metrics**: It periodically connects to defined HTTP endpoints (targets) to fetch metrics in the standard Prometheus text format.
2.  **In-Memory Storage**: It stores time-stamped data points (Metric Name, Time, Value) in memory.
3.  **Data Retention**: It automatically cleans up data older than the configured retention period (e.g., 10 minutes) to manage memory usage.
4.  **API Access**: It exposes a REST API allowing users to query the collected metrics in JSON format.

## How to Run It

### Quick Start (Windows)
I've included a batch script to launch both the **Server** (collector) and the **Target** (dummy data generator) simultaneously.

1.  Double-click `start_all.bat`
    *   *Or run `.\start_all.bat` in your terminal.*

### Manual Startup
If you prefer to run the components separately, you will need two terminal windows:

**Terminal 1 (The Dummy Target)**
This mimics a real application generating metrics.
```bash
go run target/main.go
# Output: Dummy exporter running at :8081/metrics
```

**Terminal 2 (The Mini Prometheus Server)**
This is the main application.
```bash
go run main.go
# Output: Server starting on port 9090...
```

## How it Works (Architecture)

The system is built on four pillars:

1.  **Config (`/config`)**:
    Reads `config.yaml` to determine which URLs to scrape, how often (interval), and how long to keep data (retention).

2.  **Storage (`/storage`)**:
    The purely in-memory database. It uses a `Appender` method to add data and a `Prune` method to remove old data. It uses locking (`Lock()` vs `RLock()`) to allow multiple users to read metrics while the scrapers are writing new ones.

3.  **Scraper (`/scraper`)**:
    For every target defined in the config, a separate Goroutine is spawned. It wakes up every X seconds, fetches the URL, parses the simple text format (`metric_name value`), and pushes it to Storage.

4.  **API (`/api`)**:
    A simple HTTP server that acts as the read-only interface to the Storage.
    - `GET /metrics`: Returns everything.
    - `GET /metrics/{name}`: Returns history for a specific metric.

## What can others learn from this?
This codebase is a great starting point for developers who want to understand:
*   **Go Fundamentals**: Structs, Slices, Maps, and Goroutines.
*   **Race Conditions**: Why you can't just write to a map from multiple threads and how `Mutexes` solve this.
*   **Operational Patterns**: Why "Silent Failures" are bad (we added logging!) and why "Unbounded Growth" is dangerous (we added retention!).

## Configuration
Adjust `config.yaml` to change behavior:
```yaml
server:
  port: 9090      # Change the API port via config

storage:
  retention: 10m  # Keep data for 10 minutes

targets:
  - url: http://localhost:8081/metrics
    interval: 5s  # Scrape every 5 seconds
```
