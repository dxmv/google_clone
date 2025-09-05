

FROM golang:1.24.5

WORKDIR /app

# Copy the shared module (assuming build context is from services/ directory)
COPY shared/ ./shared/

# Copy the Badger database from indexer service
COPY indexer/tmp/badger/ ./badger/

# Copy the search service source code
COPY indexer/ .

# Fix the replace directive in go.mod to point to the correct location
RUN sed -i 's|=> ../shared|=> ./shared|' go.mod

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o indexer .

# Set environment variables
ENV MONGO_CONNECTION=MONGO_CONNECTION
ENV MINIO_ENDPOINT=MINIO_ENDPOINT
ENV MINIO_ACCESS_KEY=MINIO_ACCESS_KEY
ENV MINIO_SECRET_KEY=MINIO_SECRET_KEY

# Run the application
CMD ["./indexer"]
