# Google‑Clone Project Roadmap

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
   - [ ] Write new segments
   - [ ] Background merge
3. **Search service**
   - [ ] Seperate the search stuff into a seperate go service
   - [ ] Use gRPC for communication between query-api
   - [ ] Make search concurrent
   - [ ] Create a LRU for search
   - [ ] Figure out how to make the search even faster


**Milestone:** No more storing files on my disk, more optimal everything, only query-api and search communicate directly

---

## Phase 5 – **User experience polish**

1. **Better UI**
   - [ ] Make a new figma design 
   - [ ] Implement the full design
2. **Autocomplete**
   - [ ] Top query n‑grams in a Redis trie (<5 ms).
3. **Analytics dashboard**
   - [ ] Grafana/Loki: top queries, zero‑result rate, latency histograms.

**Milestone:** Search feels “Google‑ish”; real usage graphs live.

---

## Phase 6 – **Specific crawling & improvements**

1. **Crawler V3**
   - [ ] Figure out what kind of search engine we want, and crawl those pages, like if we want a stocks serach engine or something more specific
   - [ ] Make the delay better
   - [ ] TTL old pages and re‑crawl on expiry

2. **News tab**
   - [ ] Crawler tags the pages in metadata with 'news' tag

3. **Images tab**
   - [ ] Crawler fetches the images

**Milestone:** Image and news tab when searching, and now usefull

---

## Phase 7 – **Hardening & Ops**

- [ ] Docker‑compose → Kubernetes (kind/k3s).
- [ ] Prometheus + Alertmanager (latency, error budget, disk).
- [ ] Chaos testing: kill a Ranker pod; verify graceful degradation.

---

## Phase 8 - **Final phase**

- [ ] Figure out how to host everything
- [ ] Write a README.md for everything


