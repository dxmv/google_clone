# Google‑Clone Project Roadmap

## Phase 0 – Baseline (✅ completed)

- [x] Go indexer (in‑memory postings)
- [x] Python query API that calls Go search
- [x] React mini‑frontend hitting the Python API

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

## Phase 4 – **Scale the crawl**

1. **Crawler V2**
   - [ ] Extract meta description too
   - [ ] Store the metadata and html somewhere else, like redis
   - [ ] Extract the imagelinks and backlinks
   - [ ] Figure out what kind of search engine we want, and crawl those pages
   - [ ] Make the delay better
2. **Incremental indexing**
   - [ ] Make the indexer use that new storage method
   - [ ] The indexer should store it's results in mongodb database
   - [ ] Optimize indexer + search
   - [ ] Write new segments; background merge.
   - [ ] TTL old pages and re‑crawl on expiry.

**Milestone:** Live dashboard of URLs/sec and zero‑downtime index updates.

---

## Phase 5 – **User experience polish**

1. **Full results page**

   - [ ] Snippet highlighting, favicons, domain breadcrumbs.
2. **Spell‑correction**

   - [ ] SymSpell / edit‑distance over query logs.
3. **Autocomplete**

   - [ ] Top query n‑grams in a Redis trie (<5 ms).
4. **Analytics dashboard**

   - [ ] Grafana/Loki: top queries, zero‑result rate, latency histograms.

**Milestone:** Search feels “Google‑ish”; real usage graphs live.

---

## Phase 6 – **Verticals & advanced toys**

1. **News tab**
   - [ ] RSS fetcher → extractor/indexer; tag docs `type=news`.
   - [ ] Freshness boost in Ranker.
2. **Images tab**
   - [ ] Download `<img>` sources; compute pHash; index via LSH.

**Milestone:** Instant tab switching; image similarity search demo.

---

## Phase 7 – **Hardening & Ops**

- [ ] Docker‑compose → Kubernetes (kind/k3s).
- [ ] Prometheus + Alertmanager (latency, error budget, disk).
- [ ] Chaos testing: kill a Ranker pod; verify graceful degradation.


