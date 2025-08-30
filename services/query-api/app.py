from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.middleware.gzip import GZipMiddleware
from pydantic import BaseModel
from symspellpy import SymSpell
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

def diff_cost(a: str, b: str) -> int:
    """Small edit-distance (Levenshtein) for short strings."""
    if a == b: return 0
    if not a: return len(b)
    if not b: return len(a)
    prev = list(range(len(b)+1))
    for i, ca in enumerate(a, 1):
        cur = [i]
        for j, cb in enumerate(b, 1):
            cur.append(min(cur[-1]+1, prev[j]+1, prev[j-1] + (ca != cb)))
        prev = cur
    return prev[-1]

def simple_suggestion(q: str, threshold: int = 2) -> str | None:
    # quick guards to avoid noisy suggestions
    if len(q) < 3 or "http" in q or "/" in q:
        return None
    res = sym_spell.lookup_compound(q, max_edit_distance=2)
    if not res: 
        return None
    cand = res[0].term
    if cand == q:
        return None
    # only show if the difference is meaningful
    return cand if diff_cost(q.lower(), cand.lower()) > threshold else None

@app.post("/api/search")
async def search(request: SearchRequest):
    query = request.query
    page = request.page
    count = request.count 



    response = stub.SearchQuery(search_pb2.SearchRequest(query=query, page=page, count=count))
    if not response or len(response.results) == 0:
        return {"results": [], "total": 0, "suggestion": None}
    
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
    
    suggestion = simple_suggestion(query)
    return {"results": results, "total": len(response.results), "suggestion": suggestion}
    
if __name__ == "__main__":
    channel = grpc.insecure_channel("localhost:50051")
    stub = search_pb2_grpc.SearchStub(channel)
    uvicorn.run(app, host="0.0.0.0", port=8000)
