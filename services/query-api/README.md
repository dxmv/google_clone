# Query API Service

A Python FastAPI-based HTTP gateway that provides search endpoints with autocomplete, spell correction, and intelligent query enhancement. Acts as the primary interface between the React frontend and the Go search service.

## Overview

The Query API service serves as the HTTP gateway for the search engine, providing RESTful endpoints for search queries and autocomplete suggestions. It enhances user queries with spell correction, manages autocomplete suggestions via Redis, and communicates with the search service through gRPC for optimal performance.

## Architecture

### Core Components

- **FastAPI Server**: High-performance async HTTP server
- **gRPC Client**: Efficient communication with search service
- **SymSpell Engine**: Advanced spell correction and suggestion
- **Redis Integration**: N-gram storage for autocomplete functionality
- **Query Processor**: Query normalization and enhancement
- **Response Handler**: Result formatting and error management

### Data Flow

```
HTTP Request → Query Processing → Spell Check → gRPC Call → Search Service
                    ↓                                          ↓
              Redis N-grams                              Search Results
                    ↓                                          ↓
            Autocomplete Suggestions ←                Response Assembly → HTTP Response
```

## Key Features

### 1. RESTful Search API
- **Search endpoint**: `/api/search` with pagination support
- **Autocomplete endpoint**: `/api/suggest` for real-time suggestions
- **CORS support**: Cross-origin requests enabled for frontend
- **Gzip compression**: Automatic response compression
- **Error handling**: Comprehensive error responses with proper HTTP codes

### 2. Advanced Spell Correction
- **SymSpell algorithm**: Fast and accurate spell correction
- **Bigram support**: Context-aware multi-word corrections
- **Edit distance**: Configurable maximum edit distance (default: 2)
- **Frequency-based**: Uses English word frequency dictionaries
- **Suggestion threshold**: Smart "Did you mean?" suggestions

### 3. Autocomplete System
- **N-gram indexing**: 2-gram and 3-gram phrase suggestions
- **Redis storage**: Fast in-memory suggestion retrieval
- **Query enrichment**: Automatic n-gram generation from search queries
- **Real-time updates**: Dynamic suggestion database updates
- **Prefix matching**: Efficient prefix-based suggestion lookup

### 4. Query Enhancement
- **Query normalization**: Lowercase conversion and trimming
- **Stop word handling**: Intelligent stop word processing
- **Multi-term support**: Complex query parsing and processing
- **Pagination support**: Configurable result pages and counts
- **Performance tracking**: Query timing and metrics collection

## Implementation Details

### API Endpoints

#### Search Endpoint
```python
@app.post("/api/search")
async def search(request: SearchRequest):
    # Input: query, page, count
    # Output: results, total, suggestion, query_time
```

**Request Schema**:
```python
class SearchRequest(BaseModel):
    query: str          # Search query string
    page: int = 1       # Page number (default: 1)
    count: int = 24     # Results per page (default: 24)
```

**Response Schema**:
```python
{
    "results": [SearchResult],    # Array of search results
    "total": int,                 # Total matching documents
    "suggestion": str | None,     # "Did you mean?" suggestion
    "query_time": float          # Query processing time in seconds
}
```

#### Autocomplete Endpoint
```python
@app.get("/api/suggest")
async def suggest(prefix: str):
    # Input: search prefix
    # Output: suggestion array
```

### Spell Correction Engine

The service uses SymSpell for advanced spell correction:

```python
# Initialize SymSpell with dictionaries
sym_spell = SymSpell(max_dictionary_edit_distance=2, prefix_length=7)
sym_spell.load_dictionary("frequency_dictionary_en_82_765.txt")
sym_spell.load_bigram_dictionary("frequency_bigramdictionary_en_243_342.txt")

def simple_suggestion(query: str, threshold: int = 2) -> str | None:
    # Generate spell correction suggestions
    suggestions = sym_spell.lookup_compound(query, max_edit_distance=2)
    
    # Apply intelligent filtering
    if suggestions and should_suggest(query, suggestions[0].term):
        return suggestions[0].term
    return None
```

**Features**:
- **Compound correction**: Multi-word spell correction
- **Context awareness**: Bigram dictionary for better suggestions
- **Frequency weighting**: Prefer common word corrections
- **Threshold filtering**: Avoid noisy suggestions for valid queries

### Redis Autocomplete System

The autocomplete system uses Redis for fast suggestion retrieval:

```python
# N-gram generation and storage
def enqueue_query(query: str):
    # Generate 2-grams and 3-grams from query
    ngrams = generate_ngrams(query, [2, 3])
    
    # Store in Redis sorted sets with frequency scoring
    for ngram in ngrams:
        r.zincrby(SORTED_SET_NAME, 1, ngram)
        r.zadd(LEX_SET_NAME, {ngram: 0})

def suggest(prefix: str) -> List[str]:
    # Retrieve suggestions using Redis lexicographical range
    suggestions = r.zrangebylex(LEX_SET_NAME, f"[{prefix}", f"[{prefix}\xff")
    
    # Score and rank suggestions
    scored_suggestions = []
    for suggestion in suggestions:
        score = r.zscore(SORTED_SET_NAME, suggestion)
        scored_suggestions.append((suggestion, score))
    
    # Return top suggestions sorted by frequency
    return [s[0] for s in sorted(scored_suggestions, key=lambda x: x[1], reverse=True)[:10]]
```

### gRPC Integration

Efficient communication with the search service:

```python
# gRPC client setup
search_host = os.getenv('SEARCH_HOST', 'localhost')
search_port = os.getenv('SEARCH_PORT', '50051')
channel = grpc.insecure_channel(f"{search_host}:{search_port}")
stub = search_pb2_grpc.SearchStub(channel)

# Search query execution
response = stub.SearchQuery(
    search_pb2.SearchRequest(
        query=query,
        page=page,
        count=count
    )
)
```

**Benefits**:
- **High performance**: Binary protocol with minimal overhead
- **Type safety**: Protocol buffer schema validation
- **Error handling**: Comprehensive gRPC error codes
- **Streaming support**: Large result set handling capability

## Performance Characteristics

### Response Times
- **Cached queries**: <50ms average response time
- **New queries**: <200ms including spell check and search
- **Autocomplete**: <10ms for suggestion retrieval
- **Complex queries**: <500ms for multi-term searches with corrections

### Throughput
- **Concurrent requests**: 200+ requests per second
- **Async processing**: Non-blocking I/O operations
- **Connection pooling**: Efficient resource utilization
- **Memory efficiency**: Minimal memory footprint per request

### Caching Strategy
- **Redis caching**: Fast suggestion and n-gram storage
- **Query result caching**: Temporary result caching in Redis
- **Dictionary caching**: In-memory spell correction dictionaries
- **Connection caching**: Persistent gRPC connections

## Configuration

### Environment Variables
```bash
REDIS_HOST=localhost          # Redis server hostname
REDIS_PORT=6379              # Redis server port
SEARCH_HOST=localhost        # Search service hostname  
SEARCH_PORT=50051           # Search service gRPC port
```

### Dictionary Configuration
- **Unigram dictionary**: `frequency_dictionary_en_82_765.txt`
- **Bigram dictionary**: `frequency_bigramdictionary_en_243_342.txt`
- **Max edit distance**: 2 (configurable)
- **Prefix length**: 7 (SymSpell optimization)

### Service Configuration
- **Port**: 8000 (HTTP service port)
- **CORS**: Enabled for all origins (configurable)
- **Compression**: Gzip enabled for responses >1KB
- **Timeout**: Configurable request timeouts

## Integration Points

### Upstream Dependencies
- **Search Service**: gRPC communication for search queries
- **Redis**: Autocomplete suggestions and query caching
- **Dictionary files**: Spell correction word lists

### Downstream Consumers
- **React Frontend**: Primary HTTP client for search interface
- **Mobile Apps**: RESTful API for mobile search applications
- **Third-party integrations**: External API consumers

### Service Discovery
- **Environment-based**: Service discovery via environment variables
- **Health checks**: Built-in health check endpoints
- **Load balancing**: Support for reverse proxy load balancing

## Error Handling

### HTTP Error Codes
- **200 OK**: Successful search with results
- **400 Bad Request**: Invalid query parameters
- **404 Not Found**: No results found (with suggestions)
- **500 Internal Server Error**: Service or dependency failures
- **503 Service Unavailable**: Search service unavailable

### Error Response Format
```python
{
    "error": {
        "code": "SEARCH_SERVICE_UNAVAILABLE",
        "message": "Search service is temporarily unavailable",
        "details": "gRPC connection failed"
    },
    "suggestion": "Try again in a few moments"
}
```

## Monitoring & Observability

### Key Metrics
- **Request rate**: HTTP requests per second
- **Response latency**: P50, P95, P99 response times
- **Error rate**: Failed requests percentage
- **Search service latency**: gRPC call performance
- **Cache hit rate**: Redis cache effectiveness

### Health Monitoring
- **Service health**: Basic health check endpoint
- **Dependency health**: Search service and Redis connectivity
- **Dictionary status**: Spell correction dictionary loading
- **Memory usage**: Process memory consumption

### Logging
- **Structured logging**: JSON format with contextual information
- **Query logging**: Search queries and response times
- **Error logging**: Detailed error information and stack traces
- **Access logging**: Request/response logging for analysis

## Security Considerations

### Input Validation
- **Query sanitization**: Clean and validate search queries
- **Parameter validation**: Type checking and range validation
- **XSS prevention**: Output encoding and sanitization
- **Rate limiting**: Per-IP request rate limiting

### API Security
- **CORS configuration**: Controlled cross-origin access
- **Content-Type validation**: Strict content type checking
- **Request size limits**: Maximum request payload size
- **Error information**: Sanitized error responses

## Development Features

### Development Tools
- **Auto-reload**: Automatic service restart on code changes
- **Debug mode**: Detailed error information and stack traces
- **API documentation**: Automatic OpenAPI/Swagger documentation
- **Testing support**: Built-in test client and fixtures

### Code Quality
- **Type hints**: Complete Python type annotations
- **Async/await**: Modern asynchronous programming patterns
- **Error handling**: Comprehensive exception handling
- **Code structure**: Modular and maintainable code organization

## Deployment Considerations

### Production Deployment
- **WSGI server**: Uvicorn with Gunicorn for production
- **Process management**: Multi-worker process configuration
- **Resource limits**: Memory and CPU usage limits
- **Health checks**: Container health check endpoints

### Scaling Strategies
- **Horizontal scaling**: Multiple service instances
- **Load balancing**: Nginx or cloud load balancer integration
- **Caching layers**: Redis cluster for high availability
- **Connection pooling**: Optimized database and service connections
