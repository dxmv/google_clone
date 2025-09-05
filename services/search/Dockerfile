

FROM golang:1.24.5

WORKDIR /app

# Copy the shared module (assuming build context is from services/ directory)
COPY shared/ ./shared/

# Copy the Badger database from indexer service
COPY indexer/tmp/badger/ ./app/badger/

# Copy the search service source code
COPY search/ .

# Fix the replace directive in go.mod to point to the correct location
RUN sed -i 's|=> ../shared|=> ./shared|' go.mod

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o search .

# Set environment variables
ENV MONGO_CONNECTION='mongodb+srv://dima:admin@cluster0.w5clpwy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0'
ENV MINIO_ENDPOINT=127.0.0.1:9000
ENV MINIO_ACCESS_KEY=X2999CPVW7XSBRTRTYLV
ENV MINIO_SECRET_KEY=QC0iGquthtxFAjdKaUbBe2i642w5raadRQabX2D5

# Run the application
CMD ["./search"]
