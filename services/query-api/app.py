from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from pydantic import BaseModel
import httpx  # or requests
import uvicorn
import grpc
from pb import search_pb2_grpc, search_pb2

app = FastAPI()

# Add gzip compression middleware for better performance
app.add_middleware(GZipMiddleware, minimum_size=1000)

app.add_middleware(
  CORSMiddleware,
  allow_origins = ["*"],
  allow_methods = ["*"],
  allow_headers = ["*"]
)

INDEXER_URL = "http://localhost:8080"

class SearchRequest(BaseModel):
    query: str
    page: int = 1
    count: int = 24


@app.post("/api/search")
async def search(request: SearchRequest):
    query = request.query
    page = request.page
    count = request.count 

    response = stub.SearchQuery(search_pb2.SearchRequest(query=query, page=page, count=count))
    if not response or len(response.results) == 0:
        return {"results": []}
    
    # Convert protobuf objects to JSON-serializable dictionaries
    results = []
    for result in response.results:
        result_dict = {
            "doc": {
                "url": result.Doc.url,
                "depth": result.Doc.depth,
                "title": result.Doc.title,
                "hash": result.Doc.hash,
            },
            "score": result.Score,
            "term_count": result.TermCount
        }
        results.append(result_dict)
    
    return {"results": results, "total": len(response.results)}
    
if __name__ == "__main__":
    channel = grpc.insecure_channel("localhost:50051")
    stub = search_pb2_grpc.SearchStub(channel)
    uvicorn.run(app, host="0.0.0.0", port=8000)
