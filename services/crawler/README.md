
---

### Example cards filled for your stack

#### Crawler (Go)
**Purpose**  
Discover and fetch pages, extract text + links, enqueue newly found URLs, and persist raw HTML.

**Responsibilities**
- Respect domain rules you set (e.g., allowed hosts/namespaces).  
- Deduplicate URLs and avoid re-crawling too quickly.  
- Emit content to the indexing pipeline.

**Inputs / Outputs**
- **In:** `redis:list:crawl_queue` seed URLs.  
- **Out:** Raw HTML → MinIO (`s3://raw/{hash}.html`); metadata → Mongo (`metadata`); enqueue new URLs → `crawl_queue`.

**How it works (basics)**
1) Pop URL from Redis.  
2) Fetch with timeout & retry, normalize final URL.  
3) Parse `<title>`, `<meta>`, main text; extract in-domain links; filter unwanted namespaces (e.g., `Special:`, `File:`).  
4) Persist raw HTML to MinIO; upsert metadata in Mongo (url, title, length, hash, fetched_at).  
5) Enqueue novel links (seen-set/Bloom to dedup) with politeness delay per host.

**Interfaces**  
- HTTP: `GET /healthz`  
- Queues: `crawl_queue`, `seen_urls` (set)

**Data & Storage**  
- Mongo: `metadata(url, title, doc_length, hash, fetched_at)`  
- MinIO: raw blobs keyed by content hash

**Config**  
`START_URLS`, `ALLOWED_HOSTS`, `MAX_DEPTH`, `CRAWL_CONCURRENCY`, `REQUEST_TIMEOUT_MS`, `REDIS_URL`, `MINIO_*`

**Scaling & limits**  
Stateless; scale workers. Watch: robots politeness, per-host rate limiting, memory for seen-set.

**Observability**  
- Metrics: `crawl_q_depth`, `fetch_ok`, `fetch_4xx`, `fetch_5xx`, `bytes_fetched`, `p95_fetch_ms`.  

---

#### Indexer (Go)
**Purpose**  
Transform documents into postings and maintain the inverted index + BM25 stats.

**Responsibilities**
- Tokenize/normalize text.  
- Build postings lists (term → [(docID, tf)]).  
- Track `N`, `avgDocLen`, `docLen[docID]`.

**Inputs / Outputs**
- **In:** documents from Mongo/MinIO (or a doc stream).  
- **Out:** BadgerDB keyspace for postings; doc stats in Badger/Mongo.

**How it works**  
1) Load text → tokenize, lowercase, strip stop-words (if any).  
2) Count term frequencies; write postings batches to Badger.  
3) Update `docLen` and corpus stats.

**Interfaces**  
- CLI or gRPC “index this doc/batch”.

**Config**  
`BATCH_SIZE`, `BADGER_DIR`, `STOPWORDS_PATH`

**Observability**  
- Metrics: `docs_indexed`, `postings_written`, `index_batch_ms`.

---

#### Search/Ranker (Go)
**Purpose**  
Top-K retrieval with BM25 over BadgerDB; join metadata.

**Responsibilities**
- For each query term: load postings, compute score with BM25 (`k`, `b`), merge by doc, keep a heap of top results.  
- Fetch titles/urls from Mongo; return JSON.

**Inputs / Outputs**
- **In:** gRPC `SearchRequest{query, page, page_size}`  
- **Out:** `SearchResponse{hits:[{url,title,score}], total, took_ms}`

**How it works**  
Worker pool pulls postings per term; concurrent scoring; min-heap for top-K; pagination by (from/size) or search-after.

**Interfaces**  
- gRPC: `Search.Query`  
- HTTP (optional): `/v1/search?q=…`

**Config**  
`BM25_K`, `BM25_B`, `TOP_K`, `MAX_POSTINGS_SCAN`, `MONGO_URI`

**Observability**  
- Metrics: `qps`, `p95_latency_ms`, `postings_scanned`, `heap_size`.

---

#### Query API (Python FastAPI)
**Purpose**  
Edge HTTP gateway: spell-correct, rate-limit, cache, call search/ranker.

**Responsibilities**
- `/search?q=&page=&size=` endpoint.  
- Optional SymSpell correction + segmentation with threshold.  
- Redis cache for popular queries.

**Interfaces**
- HTTP: `/search`, `/healthz`  
- gRPC client → Search/Ranker

**Config**
`REDIS_URL`, `SEARCH_GRPC_ADDR`, `SPELL_MAX_EDIT_DIST`

**Observability**
- Metrics: `cache_hit_ratio`, `p95_http_ms`, `corrections_applied`.

---

#### Frontend (React)
**Purpose**  
Search UI, results list, pagination, “I’m Feeling Lucky”, basic analytics.

**Interfaces**
- Calls Query API `/search`.  
- Static assets served by Vite/Next/Static host.

---

### End-to-end flows

#### Indexing pipeline (high level)
```mermaid
flowchart LR
  Q[Seed URLs] --> R[Redis crawl_queue]
  R --> C[Crawler]
  C -->|HTML| M((MinIO/S3))
  C -->|Metadata| MG[(Mongo)]
  C -->|New URLs| R
  MG --> I[Indexing Job]
  M --> I
  I -->|Postings + Stats| B[(BadgerDB)]
