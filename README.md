# Woogle - A Mini Wikipedia Google Clone

A from-scratch search engine that crawls, indexes, and ranks ~50k Wikipedia pages with Google-like features including autocomplete, spell correction, and BM25 ranking.

## Demo

🎬 **[Watch the Demo Video](demo.mov)**



---

## Table of Contents
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Services](#services)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Roadmap](#roadmap)


## Features

- **Crawler (Go, concurrent)**
  - BFS crawl of ~50k Wikipedia pages (Math, Philosophy, CS, …)
  - URL normalization, namespace filters, dedupe, per-host rate limiting
  - Persists **raw HTML** to Cloudflare R2 / MinIO and **metadata** (title, url, doc length, first paragraph, timestamps) to MongoDB

- **Indexer (Go)**
  - Builds a concurrent **inverted index** in BadgerDB
  - Stores term frequencies (and positions) and computes corpus stats (N, avgDocLen) for ranking

- **Ranker (Go, gRPC)**
  - BM25 scoring (configurable; default **k=2.0, b=0.9**)
  - Parallel postings scan + top-K min-heap
  - **LRU query cache** to speed up repeated/热门 queries

- **Query API (Python FastAPI)**
  - Acts as the HTTP **gateway** to the ranker (gRPC client)
  - **Autocomplete** via Redis n-grams (2- and 3-grams)
  - **Spell correction** via SymSpell with thresholded “Did you mean”
  - Pagination, basic rate limiting, and response hydration (titles/urls)

- **Frontend (React)**
  - Clean search UI with pagination, “I’m Feeling Lucky”, query time
  - Autocomplete suggestions & “Did you mean”
  - Results view for **web pages** and **images** (when available)

- **Packaging & Infra**
  - Per-service Dockerfiles + Docker Compose
  - Environment-based configuration (Mongo, Redis, MinIO/R2, Badger paths)

## Architecture

![Architecture Diagram](architecture.drawio.png)

The system follows a microservices architecture with the following components:

- **Client (React Frontend)** - User interface for search queries and results
- **Query API Service (Python FastAPI)** - HTTP gateway with autocomplete and spell correction
- **Search Service (Go)** - Core ranking engine using BM25 with gRPC interface
- **Index Service (Go)** - Builds and maintains the inverted index
- **Crawler (Go)** - Web crawler for content discovery and extraction

**Data Flow:**
1. User queries flow from the frontend → Query API → Search Service
2. Search Service queries the inverted index (BadgerDB) for fast retrieval
3. Results are enriched with metadata from MongoDB and returned to the user
4. Crawler discovers content, stores raw HTML in Cloudflare R2, and metadata in MongoDB
5. Indexer processes the stored content to build searchable indexes

## Tech Stack

**Languages & Frameworks**
- Go 1.22+ — crawler, indexer, ranker
- Python 3.11+ — FastAPI (Query API)
- React (Vite) — frontend

**Inter-service Contracts**
- gRPC + Protocol Buffers (`/proto`) between Query API ↔ Ranker
- HTTP for other communication

**Storage**
- BadgerDB — inverted index & doc stats
- MongoDB — document metadata (title, url, lengths, first paragraph, timestamps)
- Cloudflare R2 — raw HTML pages

**Caching & Queues**
- Redis — for autocomplete suggestions

**Search**
- BM25 ranking (configurable `k`, `b`; default k=2.0, b=0.9)
- SymSpell — spell correction (“Did you mean”)

**Packaging & Dev**
- Dockerfiles per service 
- Docker compose

## Services

Woogle is built as a microservices architecture with five main services, each handling a specific aspect of the search engine pipeline:

### 🕷️ [Crawler Service](services/crawler/README.md)
**Language**: Go  
**Purpose**: Discovers and fetches Wikipedia pages, extracting content and metadata

- BFS crawling strategy with ~50k page target
- Concurrent worker pool processing
- Dual storage: raw HTML in Cloudflare R2, metadata in MongoDB
- Intelligent URL filtering and deduplication
- Rate limiting and respectful crawling

### 📚 [Indexer Service](services/indexer/README.md)
**Language**: Go  
**Purpose**: Builds inverted index from crawled content for fast search retrieval

- Concurrent document processing and tokenization
- BadgerDB storage for high-performance key-value operations
- BM25 statistics calculation (avgDocLen, corpus size)
- Position-aware indexing for phrase queries
- Memory-efficient batch processing

### 🔍 [Search Service](services/search/README.md)
**Language**: Go + gRPC  
**Purpose**: Core ranking engine providing BM25-based document retrieval

- BM25 algorithm with configurable parameters (k=2.0, b=0.9)
- Concurrent postings retrieval and scoring
- Min-heap optimization for top-K results
- LRU query cache for performance (1000 entries)
- Position-based scoring enhancements

### 🌐 [Query API Service](services/query-api/README.md)
**Language**: Python + FastAPI  
**Purpose**: HTTP gateway with autocomplete and spell correction

- RESTful search and suggestion endpoints
- SymSpell-powered "Did you mean?" functionality
- Redis-based n-gram autocomplete system
- gRPC client for search service communication
- CORS support and response compression

### 🎨 [Frontend Webapp](webapp/README.md)
**Language**: React + TypeScript  
**Purpose**: Modern search interface with Google-like user experience

- Real-time autocomplete with debounced input
- "I'm Feeling Lucky" functionality
- Responsive design with Tailwind CSS
- Result pagination and image search
- Performance optimized with Vite

Each service is containerized with Docker and can be scaled independently. The services communicate through well-defined APIs (HTTP/gRPC) and share data through MongoDB, Redis, and Cloudflare R2 storage.

## Quick Start

### Prerequisites

- **Go 1.22+** for crawler and indexer services
- **Docker & Docker Compose** for running the full stack
- **MongoDB** and **Cloudflare R2** credentials (for data storage)

### Step 1: Set Up Environment

First, clone the repository:
```bash
git clone https://github.com/dxmv/google_clone
cd google_clone
```

### Step 2: Crawl Pages

The crawler discovers and stores web pages. You'll need to configure your storage credentials first.

1. **Configure environment variables** in `services/crawler/start.sh`:
   ```bash
   export MONGO_CONNECTION='your-mongodb-connection-string'
   export MINIO_ENDPOINT='your-minio-or-r2-endpoint'
   export MINIO_ACCESS_KEY='your-access-key'
   export MINIO_SECRET_KEY='your-secret-key'
   ```

2. **Run the crawler**:
   ```bash
   cd services/crawler
   chmod +x start.sh
   ./start.sh
   ```

   The crawler will:
   - Start from Wikipedia seed URLs
   - Crawl ~50k pages using BFS
   - Store raw HTML in MinIO/R2
   - Save metadata (title, URL, content length) in MongoDB
   - Take 10-30 minutes depending on your connection

### Step 3: Build Search Index

After crawling, build the inverted index for fast search:

1. **Configure the same environment variables** in `services/indexer/start.sh`

2. **Run the indexer**:
   ```bash
   cd services/indexer
   chmod +x start.sh
   ./start.sh
   ```

   The indexer will:
   - Read crawled pages from storage
   - Build inverted index with term frequencies and positions
   - Store index in BadgerDB for fast retrieval
   - Calculate corpus statistics for BM25 ranking
   - Take 5-15 minutes depending on corpus size

### Step 4: Launch the Search Engine

With your corpus crawled and indexed, start the full search engine:

```bash
# From the project root
docker-compose up --build
```

This will start:
- **Redis** (port 6379) - for autocomplete caching
- **Search Service** (port 50051) - gRPC ranking engine
- **Query API** (port 8000) - HTTP gateway with spell correction
- **Frontend** (port 5173) - React search interface

### Step 5: Search!

Open your browser and navigate to:
```
http://localhost:5173
```

You now have a fully functional search engine! Try searching for:
- `mathematics` - broad topic with many results
- `philosophy` - another well-covered area
- `computer science` - technical content
- `logic` - interdisciplinary topic

### Quick Development Setup

For development, you can also run services individually:

```bash
# Terminal 1 - Search Service
cd services/search && go run .

# Terminal 2 - Query API
cd services/query-api && python -m uvicorn app:app --reload --port 8000

# Terminal 3 - Frontend
cd webapp && npm run dev
```

**Note:** Make sure Redis is running (`docker run -d -p 6379:6379 redis:7-alpine`) for autocomplete to work.

## Configuration

The system requires several environment variables for proper operation. Here's a comprehensive guide to configuring each service:

### Core Storage Configuration

These credentials are needed by both the **Crawler** and **Indexer** services:

```bash
# MongoDB connection (for document metadata)
MONGO_CONNECTION='mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority'

# Cloudflare R2 / MinIO configuration (for raw HTML storage)
MINIO_ENDPOINT='your-r2-endpoint.r2.cloudflarestorage.com'
MINIO_ACCESS_KEY='your-access-key'
MINIO_SECRET_KEY='your-secret-key'
```

### Service-Specific Configuration

#### 1. Crawler Service
**Location**: `services/crawler/start.sh`

```bash
export MONGO_CONNECTION='your-mongodb-connection-string'
export MINIO_ENDPOINT='your-r2-endpoint.r2.cloudflarestorage.com'
export MINIO_ACCESS_KEY='your-access-key'
export MINIO_SECRET_KEY='your-secret-key'
```

**What it does**:
- `MONGO_CONNECTION`: Stores document metadata (title, URL, content length, first paragraph, images)
- `MINIO_*`: Stores raw HTML content using content hash as filename
- Crawler will create buckets automatically if they don't exist

#### 2. Indexer Service
**Location**: `services/indexer/start.sh`

```bash
export MONGO_CONNECTION='your-mongodb-connection-string'
export MINIO_ENDPOINT='your-r2-endpoint.r2.cloudflarestorage.com'
export MINIO_ACCESS_KEY='your-access-key'
export MINIO_SECRET_KEY='your-secret-key'
```

**What it does**:
- Reads crawled HTML from MinIO/R2 storage
- Fetches document metadata from MongoDB
- Builds inverted index and stores it locally in BadgerDB (`./tmp/badger/`)

#### 3. Query API Service (Python FastAPI)
**Environment variables** (set in docker-compose or runtime):

```bash
REDIS_HOST=localhost          # Redis host for autocomplete cache
REDIS_PORT=6379              # Redis port
SEARCH_HOST=localhost        # Search service gRPC host
SEARCH_PORT=50051           # Search service gRPC port
```

**What it does**:
- `REDIS_*`: Connects to Redis for autocomplete n-gram storage and query caching
- `SEARCH_*`: Connects to the Go search service via gRPC for ranking

#### 4. Search Service (Go gRPC)
**No additional environment variables needed**

- Reads from local BadgerDB index (built by indexer)
- Serves gRPC requests on port `50051`
- Uses LRU cache for query results (1000 entries by default)

#### 5. Frontend (React)
**Environment variables**:

```bash
VITE_QUERY_API_URL=http://localhost:8000    # Query API endpoint
```

**What it does**:
- `VITE_QUERY_API_URL`: Points to the Query API service for search requests

### Docker Compose Configuration

The `docker-compose.yml` automatically sets up:

- **Redis**: Port 6379, with persistent volume
- **Service networking**: All services can communicate using service names
- **Volume mounting**: BadgerDB data persisted in `search_data` volume
- **Environment injection**: Automatically sets service discovery variables

### Configuration Files

#### BM25 Ranking Parameters
**Location**: `services/search/search.go`

```go
var K = 2.0    // Term frequency saturation point
var B = 0.9    // Field length normalization (0=no normalization, 1=full)
```

#### Crawler Seed URLs
**Location**: `services/crawler/config.go`

⚠️ **Important**: The example configuration files contain placeholder credentials.

## Roadmap
### Phase 0 – **Baseline**
1. **Simple frontend**
   - [x] Just a simple react page with an input

2. **Query api**
   - [x] Create a fast-api app
   - [x] Expose search endpoint
   - [x] Call the endpoint on frontend

3. **Indexer**
   - [x] Find a corpus
   - [x] Index, only in memory for now
   - [x] Use badgerDB to store the index stuff
   - [x] Query-api calls the search endpoint

**Milestone:** Basic search


### Phase 1 – **Find things** (Crawler V1)

1. **Minimal crawler**
   - [x] BFS from a seed list (only work with wikipedia for now)
   - [x] Skip some links like '#...' and handle relative links 
   - [x] Store raw HTML 
2. **Extraction**
   - [x] Extract metadata from html (title,meta_description,outlinks,depth,url)
   - [x] Clean up the crawler code
   - [x] In index, use the crawleded pages & their metadata (remove the old docmeta logic)

**Milestone:** Crawl across ~1k pages & test search on the frontend

### Phase 2 – **Talk better** 

- [x] Make a simple gRPC ping-pong communication between the query-api and indexer, just to see how gRPC works
- [x] Draft the new search .proto 
- [x] gRPC between QueryApi -> Indexer
- [x] Add pagination to the query-api & indexer
- [x] Add a simple search results page, that uses pagination


**Milestone:** All services speak gRPC and we have a working frontend demo


### Phase 3 – **Rank smarter**
   - [x] Concurrency in indexer
   - [x] Concurrency in crawler
   - [x] Modify crawler to get the content length for each document
   - [x] Modify the indexer, to use k=1.2 b=0.75, and to calculate avg document length at the start
   - [x] BM25 in indexer, at least for now

**Milestone:** Have a working demo that use BM25 & crawl all 'Math' wikipedia under 2 mins

### Phase 4 – **Scale the crawl & index**

1. **Crawler V2**
   - [x] Abstract the storage mechanism, so we can just plug in something else if we want
   - [x] Clean up the code a bit
   - [x] Save crawled time also
   - [x] Store the html in minIO
   - [x] Store doc metadata in mongodb
2. **Indexer V2**
   - [x] Abstract the reading files
   - [x] Make the indexer use that new storage method when indexing
   - [x] Put the corpus type inside of db, and don't save metadata in badger anymore, use mongo's metadata
   - [x] Save position of each word

3. **Search service**
   - [x] Seperate the search stuff into a seperate go service
   - [x] Use gRPC for communication between query-api
   - [x] Make search concurrent
   - [x] Figure out how to make the search even faster
      - [x] Save doc length for each file in badger db
      - [x] Use that database instead of mongo for search
      - [x] Only fetch from mongodb in the end of search (to get paginated results) 
   - [x] Create a LRU for search, where we'll store results for a query
   - [x] Batch metadata request
   - [x] Use heap to sort results
   - [x] Utilize position in search


**Milestone:** No more storing files on my disk, more optimal everything, only query-api and search communicate directly. <1s for queries with a lot of results like 'logic' or 'math'


### Phase 5 – **Specific crawling & improvements**

1. **Crawler V3**
   - [x] Figure out what kind of search engine we want, and crawl those pages, like if we want a stocks serach engine or something more specific - Wikipedia
   - [x] Add a user agent
   - [x] Remove outlinks from metadata in crawler
   - [x] Add image links from the page to metadata
   - [x] Add first paragraph to metadata
   - [x] Batch write to both minio and mongodb
   

2. **Improvements**
   - [x] Use cloud storage
   - [x] 50k-100k pages craweled
   - [x] Index those pages

**Milestone:** Optimized crawler, and a large corpus.

### Phase 6 – **User experience polish**

1. **Better UI**
   - [x] Make a new figma design 
   - [x] Implement the full design
   - [x] Display query time
   - [x] I'm feeling lucky

2. **Autocomplete**
   - [x] Top query n‑grams in redis
   - [x] Display the redis suggest on frontend

3. **Spell checking**
   - [x] 'Did you mean' text for spell fixing

**Milestone:** Google user experience

### Phase 7 – **Dockerization**

- [x] Dockerise everything
   - [x] Docker for crawler
   - [x] Docker for indexer
   - [x] Docker for query-api
   - [x] Docker for search
   - [x] Docker for frontend
- [x] Docker compose

### Phase 8 - **Final phase**

- [x] Record a demo video
- [x] Write a README.md for everything

---
### Next Ideas
- [ ] More specialized crawler
- [ ] Bigger corpus (more pages craweled)
- [ ] Tag-aware scoring (boost <h1>, <title>)
- [ ] News tab
- [ ] Actually host the app
- [ ] Snippet highlighting by positions