# Google‑Clone Project Roadmap

## Phase 0 – Baseline (✅ completed)

- [x] Go indexer (in‑memory postings)
- [x] Python query API that calls Go search
- [x] React mini‑frontend hitting the Python API

---


## Phase 1 – **Find things** (Crawler V1)

1. **Minimal crawler**
   - [ ] BFS from a seed list
   - [ ] Only work with wikipedia for now
   - [ ] Skip some links like 'mailto:', and handle relative links 
   - [ ] Do a simple delay between requests
   - [ ] Store raw HTML + discovery metadata (`url`, status, fetch time) in object storage (local FS/S3/MinIO).
2. **Extraction**
   - [ ] Extract meta_title, meta_description and outlinks and store them in the metadata
   - [ ] In index, use the crawleded pages & their metadata (remove the old docmeta logic)

**Milestone:** Crawl across ~10k pages & test search on the frontend
---

## Phase 2 – **Talk better** (Service Mesh & Protocols)

- [ ] Use gRPC instead of HTTP for service communication
- [ ] Maintain shared **proto** definitions; generate stubs for Go & Python (Buf or `protoc`).
- [ ] Add pagination to the query-api
- [ ] Add a simple search results page, that uses pagination

**Milestone:** All services speak gRPC and we have a working demo
---

## Phase 3 – **Rank smarter**

1. **Ranking service (C++)**
   - [ ] BM25 + field boosts.
   - [ ] gRPC: `RankDocuments(req{query_terms, candidate_ids}) → ranked_ids`.
2. **Indexer changes**
   - [ ] Return top‑*N* candidates quickly and call Ranker for final ordering.
3. **Evaluation harness**
   - [ ] YAML gold‑set; compute NDCG / MAP in CI; fail on regressions.

**Milestone:** Side‑by‑side relevance improvement with metrics.
---

## Phase 4 – **Scale the crawl**

1. **Crawler V2**
   - [ ] Make the delay better
   - [ ] Work with the whole world wide web
   - [ ] Extract the imagelinks and backlinks
   - [ ] Redis frontier (URL, depth, priority).
   - [ ] Stateless workers; content‑hash deduplication.
2. **Incremental indexing**
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


