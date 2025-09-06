# Indexer Service

A high-performance inverted index builder written in Go that transforms crawled web pages into searchable data structures optimized for fast retrieval and BM25 ranking.

## Overview

The Indexer service processes raw HTML content and metadata from the crawler to build a comprehensive inverted index. It tokenizes text content, calculates term frequencies and positions, computes corpus statistics, and stores everything in BadgerDB for ultra-fast search operations.

## Architecture

### Core Components

- **Document Processor**: Concurrent workers that parse HTML and extract searchable text
- **Tokenization Engine**: Advanced text processing with normalization and filtering
- **Inverted Index Builder**: Creates term → document postings mappings
- **Statistics Calculator**: Computes BM25 ranking parameters (avgDocLen, N, docLen)
- **Storage Manager**: Efficient BadgerDB operations with batch writes

### Data Flow

```
MongoDB Metadata → Document Queue → Worker Pool → {Tokenize, Extract} → Postings
MinIO HTML Content ↗                                                      ↓
                                                            BadgerDB ← Merge & Store
```

## Key Features

### 1. Concurrent Processing
- **Worker pool**: Parallel document processing (CPU-core based scaling)
- **Pipeline architecture**: Overlapped I/O and computation
- **Memory management**: Controlled memory usage during large batch processing
- **Graceful scaling**: Adapts to available system resources

### 2. Advanced Text Processing
- **HTML parsing**: Clean text extraction from raw HTML using goquery
- **Tokenization**: Unicode-aware word boundary detection
- **Normalization**: Lowercase conversion, punctuation handling
- **Position tracking**: Word positions for phrase queries and snippet generation
- **Content filtering**: Removes boilerplate and navigation elements

### 3. Inverted Index Construction
- **Term-document mapping**: Efficient postings list construction
- **Frequency calculation**: Term frequency (TF) computation per document
- **Position indexing**: Word positions for advanced query features
- **Batch processing**: Memory-efficient handling of large document collections
- **Incremental updates**: Support for index updates without full rebuilds

### 4. BM25 Statistics
- **Corpus statistics**: Total document count (N) and average document length
- **Document lengths**: Individual document length storage for ranking
- **Term statistics**: Document frequency (DF) for IDF calculations
- **Real-time updates**: Statistics maintained during index construction

## Implementation Details

### Document Processing Pipeline

1. **Metadata Retrieval**: Fetch document metadata from MongoDB
2. **Content Loading**: Load raw HTML from MinIO/R2 storage
3. **HTML Parsing**: Extract clean text content using goquery
4. **Tokenization**: Split text into normalized terms with positions
5. **Posting Generation**: Create term → (docID, tf, positions) mappings
6. **Statistics Update**: Update corpus and document-level statistics
7. **Storage**: Persist postings and statistics to BadgerDB

### Storage Schema

**BadgerDB Key-Value Structure**:

```go
// Postings: term -> []Posting
type Posting struct {
    DocID     []byte    // Document hash (SHA-256)
    TF        uint32    // Term frequency in document
    Positions []uint32  // Word positions in document
}

// Document lengths: "doclen:" + docID -> length
// Corpus stats: "stats" -> Stats
type Stats struct {
    AvgDocLength float64  // Average document length
    TotalDocs    uint32   // Total number of documents
}
```

**Index Organization**:
- **Postings**: `term` → serialized posting list
- **Doc lengths**: `doclen:{hash}` → document length
- **Corpus stats**: `stats` → global statistics
- **Metadata cache**: Document metadata for result enrichment

### Tokenization Algorithm

The indexer implements sophisticated text processing:

```go
func Tokenize(text string) (map[string]uint32, []TokenPosition) {
    // 1. Unicode normalization
    // 2. Word boundary detection
    // 3. Case normalization
    // 4. Punctuation filtering
    // 5. Position tracking
    // 6. Frequency counting
}
```

**Features**:
- **Unicode support**: Proper handling of international characters
- **Word boundaries**: Accurate token separation
- **Position preservation**: Maintains word order for phrase queries
- **Frequency aggregation**: Efficient term frequency calculation

### Concurrent Architecture

The indexer uses a producer-consumer pattern:

```go
// Producer: Enqueue documents from storage
go func() {
    for _, doc := range documents {
        jobs <- doc
    }
    close(jobs)
}()

// Consumers: Process documents in parallel
for i := 0; i < numWorkers; i++ {
    go worker(jobs, results, &wg)
}

// Aggregator: Merge results into final index
for result := range results {
    mergePostings(result.Postings)
    updateStats(result.DocLength)
}
```

**Benefits**:
- **Scalability**: Utilizes all available CPU cores
- **Memory efficiency**: Controlled memory usage per worker
- **Error isolation**: Worker failures don't affect other workers
- **Progress tracking**: Real-time processing status

## Performance Characteristics

### Throughput
- **Target**: Index ~50,000 documents in 5-15 minutes
- **Concurrency**: CPU-core based worker scaling
- **Memory usage**: ~1-2GB for typical Wikipedia corpus
- **Storage efficiency**: Compressed postings with delta encoding

### Scalability
- **Horizontal**: Stateless workers allow easy scaling
- **Vertical**: Memory and CPU usage scales with corpus size
- **Storage**: BadgerDB provides excellent read/write performance
- **Index size**: Typically 10-20% of original content size

### Optimization Techniques
- **Batch operations**: Grouped database writes for efficiency
- **Memory pooling**: Reused buffers to reduce GC pressure
- **Delta compression**: Compressed posting lists save storage
- **Lazy loading**: On-demand document loading reduces memory usage

## Index Quality Metrics

### Completeness
- **Term coverage**: Percentage of unique terms indexed
- **Document coverage**: Successfully processed documents
- **Position accuracy**: Correct word position tracking
- **Metadata preservation**: Complete document metadata retention

### Efficiency
- **Index compression ratio**: Storage efficiency vs. raw text
- **Query performance**: Average term lookup time
- **Update performance**: Incremental index modification speed
- **Memory footprint**: RAM usage during indexing and querying

## Integration Points

### Upstream Dependencies
- **Crawler Service**: Provides raw HTML content and metadata
- **MongoDB**: Document metadata and corpus information
- **MinIO/R2**: Raw HTML content storage

### Downstream Consumers  
- **Search Service**: Queries the inverted index for document retrieval
- **Analytics**: Index statistics for corpus analysis
- **Monitoring**: Performance metrics and health checks

### Storage Dependencies
- **BadgerDB**: High-performance key-value store for index data
- **File system**: Local storage for BadgerDB files
- **Memory**: RAM for in-memory processing and caching

## Monitoring & Observability

### Key Metrics
- **Documents processed**: Total and rate of document indexing
- **Index size**: Number of unique terms and postings
- **Processing speed**: Documents per second throughput
- **Memory usage**: Peak and average memory consumption
- **Error rates**: Failed document processing percentage

### Health Checks
- **Storage connectivity**: BadgerDB and external storage health
- **Memory pressure**: Available memory and GC statistics
- **Processing queue**: Backlog and processing rate
- **Index integrity**: Validation of stored data structures

## Error Handling & Recovery

### Fault Tolerance
- **Document failures**: Skip malformed documents, continue processing
- **Storage errors**: Retry with exponential backoff
- **Memory pressure**: Graceful degradation and cleanup
- **Corruption detection**: Validate index integrity

### Recovery Mechanisms
- **Checkpoint system**: Periodic progress saves
- **Incremental processing**: Resume from last successful state
- **Validation**: Post-processing integrity checks
- **Rollback**: Revert to last known good state

## Security & Compliance

### Data Protection
- **Content sanitization**: Remove potentially harmful content
- **Access control**: Secure storage access credentials
- **Audit logging**: Track all index modifications
- **Data retention**: Configurable content lifecycle management

### Performance Security
- **Resource limits**: Prevent resource exhaustion attacks
- **Input validation**: Sanitize all external input
- **Error isolation**: Prevent cascading failures
- **Monitoring**: Detect unusual processing patterns