from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import httpx  # or requests
import uvicorn
import grpc
from pb import ping_pb2_grpc,ping_pb2

app = FastAPI()
app.add_middleware(
  CORSMiddleware,
  allow_origins = ["*"],
  allow_methods = ["*"],
  allow_headers = ["*"]
)

INDEXER_URL = "http://localhost:8080"

class SearchRequest(BaseModel):
    query: str
    count: int = 24


@app.get("/")
async def read_root():
    return {"message": "Hello World"}

@app.post("/api/search")
async def search(request: SearchRequest):
    query = request.query
    count = request.count
    
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{INDEXER_URL}/search",
            params={"q": query}
        )
        results = response.json()
        if len(results) > count:
            results = results[:count]
    return results
    
if __name__ == "__main__":
    channel = grpc.insecure_channel("localhost:50051")
    stub = ping_pb2_grpc.HealthStub(channel)
    response = stub.Ping(ping_pb2.PingRequest())
    print(response)
    # uvicorn.run(app, host="0.0.0.0", port=8000)
