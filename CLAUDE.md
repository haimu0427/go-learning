# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This is a Go language learning repository containing multiple independent sub-projects. Each directory is a standalone Go module with its own `go.mod` file. The repository demonstrates various Go technologies and frameworks through practical examples.

**Go version:** 1.24.5

## Project Structure

Each subdirectory is an independent Go module:

- **go-ai/** - OpenAI API integration example using `github.com/sashabaranov/go-openai`
- **go-crawler/** - Web scraping examples
  - `go-colly/colly/` - Colly framework crawler with CSV export, YAML config, and rate limiting
  - `py-scrapy/` - Lightweight Python scraper using requests+BeautifulSoup
- **go-docker/** - Docker learning notes and examples
- **go-kafka/** - Kafka learning progression in 4 steps
  - `step2_producer/` - Producer and consumer basics
  - `step3_consumergroup/` - Consumer group implementation
  - `step4_batchingAndCompression/` - Performance optimization with batching and compression
- **go-kratos/helloworld/** - Kratos microservice framework example with proto definitions
- **go-redis/learn01/** - Redis operations including cache-aside pattern
- **go-resty/** - HTTP client examples using `github.com/go-resty/resty/v2`
- **go-wire/** - Dependency injection using Google Wire
- **go-zero/** - go-zero framework examples
  - `gemini_1_uer/` - User API service with MySQL and Redis
  - `gemini_1_mall-rpc/` - RPC service example
- **grpc-learn/** - gRPC server/client examples with streaming
- **lua/** - Lua scripting examples
- **minio/** - MinIO object storage examples
- **py-Vs-go/** - Python vs Go comparison examples

## Common Commands

Since each subdirectory is an independent module, you must cd into the specific directory before running Go commands:

```bash
# Build any module
cd <module-directory>
go build

# Run any module
cd <module-directory>
go run main.go

# Tidy dependencies
cd <module-directory>
go mod tidy
```

### Specific Module Commands

**go-zero user service:**
```bash
cd go-zero/gemini_1_uer
go run user.go -f etc/user-api.yaml
```

**Kratos helloworld:**
```bash
cd go-kratos/helloworld
make build        # Build to bin/
make api          # Generate API proto
make config       # Generate internal proto
make generate     # Run wire generation
make all          # Run all generation
```

**Python crawler:**
```bash
cd go-crawler/py-scrapy
python -m venv .venv
.venv/Scripts/activate  # Windows: .venv/Scripts/Activate.ps1
pip install -r requirements.txt
python main.py
pytest  # Run tests
```

**gRPC examples:**
```bash
cd grpc-learn
# Server
cd server && go run main.go
# Client (separate terminal)
cd client && go run main.go
```

## Architecture Notes

### go-zero Services
Follows the standard go-zero architecture:
- `*.api` files define API routes and types
- `etc/*.yaml` contains configuration
- `internal/handler/` - HTTP request handlers
- `internal/logic/` - Business logic layer
- `internal/model/` - Data models
- `internal/svc/` - Service context (dependencies)
- `internal/types/` - Generated request/response types

### Kafka Learning Progression
The kafka modules follow a 4-stage learning path:
1. Topic creation and basic setup
2. Simple producer/consumer
3. Consumer groups for parallel processing
4. Batching and compression for performance

### Redis Pattern
The `go-redis/learn01/cacheAside/` directory implements the Cache-Aside pattern with Redis as a caching layer.

### gRPC Streaming
`grpc-learn/server/main.go` demonstrates server-side streaming with the `GenerateUserReport` method which sends multiple responses over a single stream.

## Dependencies

Key Go packages used across modules:
- `github.com/IBM/sarama` - Kafka client library
- `github.com/gocolly/colly/v2` - Web scraping framework
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/go-resty/resty/v2` - HTTP client
- `github.com/google/wire` - Dependency injection
- `google.golang.org/grpc` - gRPC framework
- `github.com/go-kratos/kratos/v2` - Microservice framework
- `github.com/zeromicro/go-zero` - Go-Zero framework

## Environment

Some modules require external services:
- **go-zero**: MySQL and Redis
- **go-kafka**: Apache Kafka broker
- **go-redis**: Redis server

Check individual module READMEs for connection details.
