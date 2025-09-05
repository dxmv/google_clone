

FROM golang:1.24.5

WORKDIR /app

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the crawler service source code
COPY . .

# Build the application
RUN go build -o crawler .

# Set environment variables
ENV MONGO_CONNECTION='mongodb+srv://dima:admin@cluster0.w5clpwy.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0'
ENV MINIO_ENDPOINT=0aef907fef2d92736fdc45fe8c4cded7.r2.cloudflarestorage.com
ENV MINIO_ACCESS_KEY=40b2b4803b6384f3a003b147ef769185
ENV MINIO_SECRET_KEY=f16623682b342145ad23326ba83e6a2048480d3c2dc67b0f4cdab1b1a7cb3dfb

# Run the application
CMD ["./crawler"]
