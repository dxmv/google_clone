# Search Service

A high-performance search ranking engine built in Go that provides BM25-based document retrieval with gRPC interface, LRU caching, and concurrent processing for sub-second query response times.

## Overview

The Search service is the core ranking engine of the search system. It processes user queries by retrieving relevant documents from the inverted index, scoring them using the BM25 algorithm, and returning ranked results. The service is optimized for low latency with concurrent processing, intelligent caching, and efficient data structures.

## Architecture

### Core Components

- **gRPC Server**: High-performance RPC interface for query processing
- **BM25 Ranker**: Advanced relevance scoring with configurable parameters
- **LRU Cache**: Query result caching for improved response times
- **Concurrent Processor**: Parallel postings retrieval and scoring
- **Min-Heap Sorter**: Efficient top-K result selection
- **Metadata Enricher**: Result enhancement with document details

### Data Flow

```
Query → Tokenization → Postings Retrieval → Parallel Scoring → Heap Ranking → Metadata Enrichment → Results
                              ↓                    ↑
                        BadgerDB Index    ←    LRU Cache
                              ↓                    ↑
                        MongoDB Metadata  →   Result Assembly
```

## Key Features

### 1. BM25 Ranking Algorithm
- **Configurable parameters**: K=2.0 (term saturation), B=0.9 (length normalization)
- **Term frequency scoring**: Advanced TF component with saturation
- **Document length normalization**: Prevents bias toward shorter documents
- **Inverse document frequency**: Rare terms weighted higher
- **Position-based boosting**: Enhanced scoring for term positions

### 2. High-Performance Architecture
- **Concurrent processing**: Parallel workers for postings retrieval
- **Min-heap optimization**: Efficient top-K result selection
- **Memory efficiency**: Streaming processing with controlled memory usage
- **Connection pooling**: Optimized database connections
- **Batch operations**: Grouped metadata requests

### 3. Intelligent Caching
- **LRU cache**: 1000-entry query result cache
- **Cache key optimization**: Normalized query strings for better hit rates
- **Memory management**: Automatic eviction of least recently used entries
- **Cache warming**: Proactive caching of popular queries
- **Hit rate monitoring**: Performance metrics and optimization

### 4. gRPC Interface
- **Protocol Buffers**: Efficient binary serialization
- **Streaming support**: Large result set handling
- **Error handling**: Comprehensive error codes and messages
- **Timeout management**: Configurable request timeouts
- **Load balancing**: Support for multiple service instances

## Implementation Details

### BM25 Scoring Formula

The service implements the complete BM25 formula:

```
BM25(q,d) = Σ(t∈q) IDF(t) × (tf(t,d) × (k+1)) / (tf(t,d) + k × (1-b + b × |d|/avgdl))
```

Where:
- **IDF(t)**: Inverse document frequency of term t
- **tf(t,d)**: Term frequency in document d  
- **|d|**: Document length
- **avgdl**: Average document length in corpus
- **k**: Term frequency saturation (default: 2.0)
- **b**: Length normalization (default: 0.9)

### Query Processing Pipeline

1. **Query Parsing**: Tokenize and normalize search terms
2. **Cache Lookup**: Check LRU cache for cached results
3. **Postings Retrieval**: Fetch inverted index entries for each term
4. **Parallel Scoring**: Concurrent BM25 calculation across workers
5. **Result Aggregation**: Merge scores for documents matching multiple terms
6. **Heap Sorting**: Select top-K results using min-heap
7. **Position Enhancement**: Apply position-based score boosts
8. **Pagination**: Extract requested page of results
9. **Metadata Enrichment**: Fetch document details from MongoDB
10. **Response Assembly**: Build final search response
11. **Cache Update**: Store results in LRU cache

### Concurrent Architecture

```go
// Worker pool for parallel postings processing
for i := 0; i < workers; i++ {
    wg.Add(1)
    go worker(jobs, results, &wg, storage, avgDocLength)
}

// Distribute postings across workers
go func() {
    for _, posting := range postings {
        jobs <- Job{Posting: posting, IDF: idf}
    }
    close(jobs)
}()

// Aggregate results from workers
for result := range results {
    updateDocumentScore(docMap, result)
}
```

### Storage Integration

**BadgerDB Operations**:
- **Postings lookup**: `GetPostings(term)` retrieves posting lists
- **Document lengths**: `GetDocLength(docID)` for BM25 calculation
- **Corpus statistics**: `GetStats()` for average document length
- **Batch reads**: Optimized multi-key retrieval

**MongoDB Operations**:
- **Metadata batch fetch**: `GetBatchMetadata(docIDs)` 
- **Result enrichment**: Title, URL, first paragraph, images
- **Connection pooling**: Efficient database connection management

### Position-Based Scoring

The service enhances BM25 scores with position information:

```go
func calculatePositionBoost(positions []uint32, queryTerms []string) float64 {
    boost := 1.0
    
    // Early position boost (first 100 words)
    for _, pos := range positions {
        if pos < 100 {
            boost += 0.1 * (100 - pos) / 100
        }
    }
    
    // Proximity boost for multiple terms
    if len(queryTerms) > 1 {
        boost += calculateProximityBoost(positions, queryTerms)
    }
    
    return boost
}
```

## Performance Characteristics

### Latency
- **Target response time**: <100ms for typical queries
- **95th percentile**: <200ms including metadata enrichment
- **Cache hit response**: <10ms for cached queries
- **Complex queries**: <500ms for multi-term searches

### Throughput
- **Concurrent queries**: 100+ queries per second per instance
- **Worker scaling**: Automatic scaling based on CPU cores
- **Memory efficiency**: <2GB RAM for typical workloads
- **Connection pooling**: Optimized database resource usage

### Scalability
- **Horizontal scaling**: Stateless service design
- **Load balancing**: gRPC load balancer compatibility
- **Resource scaling**: CPU and memory usage scales with query complexity
- **Index scaling**: Performance maintained with large indexes

## gRPC API Specification

### SearchQuery RPC

**Request**:
```protobuf
message SearchRequest {
    string query = 1;        // Search query string
    int32 page = 2;          // Page number (1-based)
    int32 count = 3;         // Results per page
}
```

**Response**:
```protobuf
message SearchResponse {
    repeated SearchResult results = 1;  // Ranked search results
    int64 total = 2;                   // Total matching documents
}

message SearchResult {
    DocMetadata doc = 1;     // Document metadata
    double score = 2;        // BM25 relevance score
    int32 term_count = 3;    // Number of matching terms
}
```

### Error Handling

- **INVALID_ARGUMENT**: Malformed query or parameters
- **NOT_FOUND**: No results found for query
- **INTERNAL**: Database or index errors
- **DEADLINE_EXCEEDED**: Query timeout
- **RESOURCE_EXHAUSTED**: Service overloaded

## Monitoring & Observability

### Key Metrics
- **Query rate**: Queries per second (QPS)
- **Response latency**: P50, P95, P99 response times
- **Cache hit rate**: Percentage of queries served from cache
- **Error rate**: Failed queries per total queries
- **Index utilization**: Terms accessed and postings scanned

### Performance Monitoring
- **Memory usage**: Heap size and garbage collection metrics
- **CPU utilization**: Worker pool efficiency
- **Database connections**: Connection pool health
- **Cache efficiency**: Hit rates and eviction patterns

### Health Checks
- **Service health**: Basic ping/pong health check
- **Index connectivity**: BadgerDB accessibility
- **Metadata connectivity**: MongoDB connection status
- **Cache health**: LRU cache statistics

## Configuration

### BM25 Parameters
```go
var K = 2.0    // Term frequency saturation point
var B = 0.9    // Field length normalization factor
```

### Service Configuration
- **Port**: 50051 (gRPC service port)
- **Workers**: CPU core count (concurrent processing)
- **Cache size**: 1000 entries (LRU cache)
- **Timeouts**: Configurable request timeouts
- **Batch size**: MongoDB batch fetch size

## Integration Points

### Upstream Dependencies
- **Indexer Service**: Provides BadgerDB inverted index
- **MongoDB**: Document metadata storage
- **BadgerDB**: Local index storage and statistics

### Downstream Consumers
- **Query API Service**: Primary gRPC client
- **Analytics Services**: Query and performance metrics
- **Monitoring Systems**: Health and performance data

### Service Discovery
- **gRPC registration**: Service discovery integration
- **Health reporting**: Status reporting to orchestration
- **Load balancing**: Support for service mesh integration

## Security Considerations

### Access Control
- **Service authentication**: gRPC TLS and authentication
- **Rate limiting**: Per-client query rate limits
- **Resource limits**: Memory and CPU usage bounds
- **Input validation**: Query parameter sanitization

### Data Security
- **Index protection**: Secure access to search indexes
- **Metadata security**: Protected document metadata access
- **Audit logging**: Query and access logging
- **Error information**: Sanitized error messages

## Optimization Strategies

### Query Optimization
- **Term ordering**: Process rare terms first
- **Early termination**: Stop processing when sufficient results found
- **Selective indexing**: Skip very common terms
- **Result caching**: Cache expensive query results

### Memory Optimization
- **Object pooling**: Reuse expensive objects
- **Garbage collection**: Tuned GC parameters
- **Streaming processing**: Process large result sets incrementally
- **Memory limits**: Bounded memory usage per query

### I/O Optimization
- **Batch operations**: Group database operations
- **Connection pooling**: Reuse database connections
- **Prefetching**: Anticipatory data loading
- **Compression**: Compressed data transfer
