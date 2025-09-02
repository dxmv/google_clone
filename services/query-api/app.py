from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from pydantic import BaseModel
from symspellpy import SymSpell, Verbosity
import uvicorn
import grpc
from pb import search_pb2_grpc, search_pb2
from redis_worker import enqueue_query, suggest as suggest_from_worker


app = FastAPI()

# Add gzip compression middleware for better performance
app.add_middleware(GZipMiddleware, minimum_size=1000)

app.add_middleware(
  CORSMiddleware,
  allow_origins = ["*"],
  allow_methods = ["*"],
  allow_headers = ["*"]
)

sym_spell = SymSpell(max_dictionary_edit_distance=2, prefix_length=7)

# Load unigram dictionary
sym_spell.load_dictionary("frequency_dictionary_en_82_765.txt", term_index=0, count_index=1)

# Load bigram dictionary (optional, for multi-word context)
sym_spell.load_bigram_dictionary("frequency_bigramdictionary_en_243_342.txt", term_index=0, count_index=2)

INDEXER_URL = "http://localhost:8080"

class SearchRequest(BaseModel):
    query: str
    page: int = 1
    count: int = 24


def simple_suggestion(q: str, threshold: int = 2) -> str | None:
    # quick guards to avoid noisy suggestions
    if len(q) < 3 or "http" in q or "/" in q:
        return None
    
    words = q.split()
    suggestion = []
    for word in words:
        res = sym_spell.lookup(word, Verbosity.CLOSEST, max_edit_distance=threshold)
        if res:
            suggestion.append(res[0].term)
        else:
            suggestion.append(word)
    
    cand = " ".join(suggestion)
    if cand.lower() == q.lower():
        return None

    return cand

@app.post("/api/search")
async def search(request: SearchRequest):
    query = request.query
    page = request.page
    count = request.count 

    # enqueue query to redis queue
    enqueue_query(query)

    # get resullts from search service
    response = stub.SearchQuery(search_pb2.SearchRequest(query=query, page=page, count=count))
    print("Results: ", len(response.results))
    suggestion = simple_suggestion(query)

    # return results
    if not response or len(response.results) == 0:
        return {"results": [], "total": 0, "suggestion": suggestion}
    
    # Convert protobuf objects to JSON-serializable dictionaries
    results = []
    for result in response.results:
        result_dict = {
            "doc": {
                "url": result.Doc.url,
                "depth": result.Doc.depth,
                "title": result.Doc.title,
                "hash": result.Doc.hash,
                "images": list(result.Doc.images),
            },
            "score": result.Score,
            "term_count": result.TermCount
        }
        results.append(result_dict)
    
    print("Suggestion: ", suggestion)
    print("Query: ", query)
    print("Results: ", len(results))
    return {"results": results, "total": len(response.results), "suggestion": suggestion}
    

@app.get("/api/suggest")
async def suggest(prefix: str):
    return suggest_from_worker(prefix)

if __name__ == "__main__":
    channel = grpc.insecure_channel("localhost:50051")
    stub = search_pb2_grpc.SearchStub(channel)
    uvicorn.run(app, host="0.0.0.0", port=8000)
