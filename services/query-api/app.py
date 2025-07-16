from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def read_root():
    return {"message": "Hello World"}

@app.post("/api/search")
async def search(query: str):
    return {"message": "Hello World"}