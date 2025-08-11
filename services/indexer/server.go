package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/dxmv/google_clone/pb"
	"google.golang.org/grpc"
)

const PORT = ":50051"

type searchServer struct {
	pb.UnimplementedSearchServer
	storage        *Storage
	avgDocLength   float64
	collectionSize int64
}

func NewSearchServer(storage *Storage, avgDocLength float64, collectionSize int64) *searchServer {
	return &searchServer{storage: storage, avgDocLength: avgDocLength, collectionSize: collectionSize}
}

func (s *searchServer) SearchQuery(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	results := search(req.Query, s.storage, s.avgDocLength, s.collectionSize)
	finalResults := make([]*pb.SearchResult, len(results))
	for i, result := range results {
		finalResults[i] = &pb.SearchResult{
			Doc: &pb.DocMetadata{
				Url:   result.DocMetadata.URL,
				Depth: int32(result.DocMetadata.Depth),
				Title: result.DocMetadata.Title,
				Hash:  result.DocMetadata.Hash,
				Links: result.DocMetadata.Links,
			},
			Score:     result.Score,
			TermCount: int32(result.CountTerm),
		}
	}

	// pagination here
	offset := (req.Page - 1) * req.Count // start index
	// check if offset is out of bounds
	if offset >= int32(len(finalResults)) {
		return &pb.SearchResponse{Results: []*pb.SearchResult{}}, nil
	}
	// check if offset + count is out of bounds
	if offset+req.Count > int32(len(finalResults)) {
		finalResults = finalResults[offset:]
	} else {
		finalResults = finalResults[offset : offset+req.Count]
	}

	fmt.Println("Final results: ", finalResults)

	return &pb.SearchResponse{Results: finalResults}, nil
}

func startServer(storage *Storage) {

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server is running on port", PORT)
	grpcServer := grpc.NewServer()
	avgDocLength := storage.getStats().AvgDocLength
	collectionSize := int64(storage.getStats().TotalDocs)
	fmt.Println("Avg doc length:", avgDocLength)
	pb.RegisterSearchServer(grpcServer, NewSearchServer(storage, avgDocLength, collectionSize))
	grpcServer.Serve(lis)
}
