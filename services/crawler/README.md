
# Crawler Service

A concurrent web crawler built in Go that discovers and fetches Wikipedia pages, extracting content and metadata for the search engine pipeline.

## Overview

The Crawler service implements a breadth-first search (BFS) crawling strategy specifically designed for Wikipedia content. It starts from seed URLs across various academic and professional domains, discovers new pages through link extraction, and stores both raw HTML content and structured metadata for downstream processing.

## Architecture

### Core Components

- **Worker Pool**: Concurrent workers (default: number of CPU cores) that process crawl jobs in parallel
- **Job Queue**: Buffered channel system for distributing URLs to workers (default buffer: 10,000)
- **Visited Tracker**: Thread-safe deduplication system to avoid re-crawling pages
- **Storage Abstraction**: Pluggable storage interface supporting MongoDB + MinIO/R2

### Data Flow

```
Seed URLs → Job Queue → Worker Pool → {HTML Extraction, Link Discovery} → Storage
                ↑                                    ↓
            New URLs ←────────────── Link Filtering & Validation
```

## Key Features

### 1. Intelligent URL Filtering
- **Namespace filtering**: Skips Wikipedia special pages (`Special:`, `File:`, `Category:`, etc.)
- **Domain validation**: Only crawls `en.wikipedia.org` pages
- **URL normalization**: Handles relative links and removes fragments
- **Deduplication**: Prevents re-crawling of already visited URLs

### 2. Content Extraction
- **HTML parsing**: Extracts clean text content using goquery
- **Metadata extraction**: 
  - Page title from `<title>` tag
  - First paragraph for search snippets
  - Image URLs from `<img>` tags
  - Content length calculation
- **Link discovery**: Finds and validates outbound Wikipedia links

### 3. Concurrent Processing
- **Worker pool pattern**: Configurable number of concurrent workers
- **Rate limiting**: Respects server load with controlled concurrency
- **Batch processing**: Efficient handling of large URL queues
- **Graceful shutdown**: Proper cleanup and data flushing

### 4. Robust Storage
- **Dual storage**: Raw HTML in MinIO/R2, metadata in MongoDB
- **Content hashing**: SHA-256 hashes for deduplication and integrity
- **Batch writes**: Optimized database operations
- **Error handling**: Resilient to network and storage failures

## Implementation Details

### URL Processing Pipeline

1. **URL Validation**: Check domain, namespace, and format
2. **Deduplication**: Skip if already crawled or queued
3. **HTTP Fetch**: Get page content with timeout and retries
4. **Content Processing**: Extract text, metadata, and links
5. **Storage**: Persist HTML and metadata to respective stores
6. **Link Enqueueing**: Add discovered URLs to crawl queue

### Storage Schema

**MongoDB Document**:
```go
type DocMetadata struct {
    URL            string    // Original page URL
    Depth          int       // Crawl depth from seed
    Title          string    // Page title
    Hash           string    // SHA-256 content hash
    ContentLength  int       // Text content length
    CrawledAt      time.Time // Timestamp
    FirstParagraph string    // First paragraph for snippets
    Images         []string  // Image URLs found on page
}
```

**MinIO/R2 Storage**:
- **Key**: `{hash}.html`
- **Content**: Raw HTML content
- **Bucket**: Automatically created if not exists

### Configuration

The crawler uses a configuration struct that defines:

```go
type Config struct {
    StartLinks   []string      // Seed URLs for crawling
    MaxDepth     int          // Maximum crawl depth (default: 1)
    JobsBuffer   int          // Job queue buffer size (default: 10,000)
    MaxRounds    int          // Maximum crawl rounds (default: 1,000)
    NumWorkers   int          // Concurrent workers (default: CPU cores)
    MongoUri     string       // MongoDB connection string
    MinioClient  *minio.Client // MinIO/R2 client
}
```

### Seed URLs

The crawler starts from carefully selected Wikipedia pages across major academic domains:

- **STEM**: Mathematics, Computer Science, Physics, Chemistry, Biology, Astronomy
- **Social Sciences**: Philosophy, Psychology, Economics, Business
- **Humanities**: Literature, History, Art, Music, Religion
- **Professional**: Medicine, Engineering, Law, Finance

## Performance Characteristics

### Scalability
- **Horizontal**: Stateless design allows multiple crawler instances
- **Vertical**: Worker pool scales with available CPU cores
- **Memory**: Efficient visited set using hash-based deduplication

### Throughput
- **Target**: ~50,000 Wikipedia pages in 10-30 minutes
- **Concurrency**: CPU-core based worker pool (typically 4-16 workers)
- **Rate limiting**: Respectful crawling to avoid overwhelming servers

### Error Handling
- **Network failures**: Automatic retries with exponential backoff
- **Storage errors**: Graceful degradation and error logging
- **Invalid content**: Skip and continue processing
- **Resource limits**: Memory and connection pool management

## Monitoring & Observability

### Key Metrics
- **Pages crawled**: Total successful page fetches
- **Errors encountered**: HTTP errors, parsing failures, storage issues
- **Queue depth**: Pending URLs in job queue
- **Processing rate**: Pages per second throughput
- **Storage efficiency**: Deduplication ratio

### Logging
- **Structured logging**: JSON format with contextual information
- **Error tracking**: Detailed error messages with stack traces
- **Progress reporting**: Regular status updates during crawling
- **Performance metrics**: Timing information for optimization

## Integration Points

### Upstream Dependencies
- **MongoDB**: Document metadata storage
- **MinIO/Cloudflare R2**: Raw HTML content storage
- **Network**: HTTP client for page fetching

### Downstream Consumers
- **Indexer Service**: Reads crawled content for search index building
- **Search Service**: Uses metadata for result ranking and display

## Security Considerations

- **User Agent**: Identifies crawler to web servers
- **Rate limiting**: Prevents overwhelming target servers
- **Content validation**: Sanitizes and validates extracted content
- **Error isolation**: Prevents malformed content from affecting system
- **Credential management**: Secure storage of database and cloud credentials
