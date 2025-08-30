# Google‑Clone Project Roadmap

This repository is a from-scratch mini search engine designed to mirror the core ideas behind a web search stack while staying small enough to understand. It has a Go crawler and indexer, a Go-based ranking pipeline (BM25), a Python FastAPI query layer, and a React front-end. Storage is split by concern: HTML in object storage (MinIO/S3), document metadata in MongoDB, and an inverted index in BadgerDB (for fast local lookups). Everything speaks gRPC under the hood (or is moving there).

### Why this exists
 - Practice real IR: tokenization, postings, tf-idf/BM25, snippets, pagination.
 - Systems thinking: concurrency in Go, backpressure, batching, caches.
 - Cloud-ish storage: object storage for the corpus, DB for metadata, KV for the index.
 - Pragmatism: ship an MVP first; make it faster and more robust later.

### What it does (today)
 - Crawls a bounded slice of Wikipedia (seeded to Math/Philosophy), storing raw HTML in MinIO and metadata in MongoDB.
 - Builds an inverted index in BadgerDB, computes collection stats (doc count, average doc length), and serves BM25 queries with k=1.2, b=0.75.
 - Exposes a simple FastAPI endpoint the React UI calls to show paginated results.

### Tech choices (and why)
 - Go for crawler/indexer: easy concurrency, strong stdlib, fast binaries.
 - BadgerDB for the index: blazing fast local KV; perfect for a single-writer, single-box MVP.
 - MinIO/S3 for raw HTML: cheap, durable, scalable; decouples storage from compute.
 - MongoDB for metadata: flexible schema; easy querying.
 - FastAPI: quick to ship a clean API the frontend can hit.
 - gRPC: stable, language-agnostic contracts between services.


---
## Phase 0 – **Baseline**
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
---


## Phase 1 – **Find things** (Crawler V1)

1. **Minimal crawler**
   - [x] BFS from a seed list (only work with wikipedia for now)
   - [x] Skip some links like '#...' and handle relative links 
   - [x] Store raw HTML 
2. **Extraction**
   - [x] Extract metadata from html (title,meta_description,outlinks,depth,url)
   - [x] Clean up the crawler code
   - [x] In index, use the crawleded pages & their metadata (remove the old docmeta logic)

**Milestone:** Crawl across ~1k pages & test search on the frontend
---

## Phase 2 – **Talk better** 

- [x] Make a simple gRPC ping-pong communication between the query-api and indexer, just to see how gRPC works
- [x] Draft the new search .proto 
- [x] gRPC between QueryApi -> Indexer
- [x] Add pagination to the query-api & indexer
- [x] Add a simple search results page, that uses pagination


**Milestone:** All services speak gRPC and we have a working frontend demo
---

## Phase 3 – **Rank smarter**
   - [x] Concurrency in indexer
   - [x] Concurrency in crawler
   - [x] Modify crawler to get the content length for each document
   - [x] Modify the indexer, to use k=1.2 b=0.75, and to calculate avg document length at the start
   - [x] BM25 in indexer, at least for now

**Milestone:** Have a working demo that use BM25 & crawl all 'Math' wikipedia under 2 mins
---

## Phase 4 – **Scale the crawl & index**

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

---

## Phase 5 – **Specific crawling & improvements**

1. **Crawler V3**
   - [x] Figure out what kind of search engine we want, and crawl those pages, like if we want a stocks serach engine or something more specific - Wikipedia
   - [x] Add a user agent
   - [x] Remove outlinks from metadata in crawler
   - [x] Add image links from the page to metadata
   - [x] Add first paragraph to metadata
   - [x] Batch write to both minio and mongodb
   

2. **Improvements**
   - [x] Use cloud storage
   - [ ] 50k-100k pages craweled

**Milestone:** Optimized crawler, and a large corpus.

---

## Phase 6 – **User experience polish**

1. **Better UI**
   - [ ] Make a new figma design 
   - [ ] Implement the full design

2. **Autocomplete**
   - [ ] Top query n‑grams in a Redis trie (<5 ms).

3. **Spell checking**
   - [ ] 'Did you mean' text for spell fixing

**Milestone:** Google user experience

---

## Phase 7 – **Hardening & Ops**

- [ ] Dockerise everything

---

## Phase 8 - **Final phase**

- [ ] Record a demo video
- [ ] Write a README.md for everything

---
## Improvement ideas
- [ ] More specialized crawler
- [ ] Bigger corpus (more pages craweled)
- [ ] In indexer save html tag for each word, in search use html tag to add a bonus to the score
- [ ] News tab
- [ ] Image indexing
- [ ] Actually host the app
- [ ] Position highlighting


