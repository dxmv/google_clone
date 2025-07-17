from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

app = FastAPI()
app.add_middleware(
  CORSMiddleware,
  allow_origins = ["*"],
  allow_methods = ["*"],
  allow_headers = ["*"]
)

class SearchRequest(BaseModel):
    query: str

@app.get("/")
async def read_root():
    return {"message": "Hello World"}

@app.post("/api/search")
async def search(request: SearchRequest):
    return {"message": f"Received search query: {request.query}"}
