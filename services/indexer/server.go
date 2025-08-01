package main

import (
	"context"
	"fmt"
	"log"
	"net"

	badger "github.com/dgraph-io/badger/v4"
	pb "github.com/dxmv/google_clone/pb"
	"google.golang.org/grpc"
)

const PORT = ":50051"

type searchServer struct {
	pb.UnimplementedSearchServer
	db *badger.DB
}

func NewSearchServer(db *badger.DB) *searchServer {
	return &searchServer{db: db}
}

func (s *searchServer) SearchQuery(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	results := search(req.Query, s.db)
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

	return &pb.SearchResponse{Results: finalResults}, nil
}

func startServer(db *badger.DB) {

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server is running on port", PORT)
	grpcServer := grpc.NewServer()
	pb.RegisterSearchServer(grpcServer, NewSearchServer(db))
	grpcServer.Serve(lis)
}
