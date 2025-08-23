package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/dxmv/google_clone/search/pb"
	shared "github.com/dxmv/google_clone/shared"
	"google.golang.org/grpc"
)

const PORT = ":50051"

type searchServer struct {
	pb.UnimplementedSearchServer
	storage        *shared.Storage
	avgDocLength   float64
	collectionSize int64
	cache          *LRUCache[string, []SearchResult]
}

func NewSearchServer(storage *shared.Storage, avgDocLength float64, collectionSize int64, cache *LRUCache[string, []SearchResult]) *searchServer {
	return &searchServer{storage: storage, avgDocLength: avgDocLength, collectionSize: collectionSize,
		cache: cache,
	}
}

func (s *searchServer) SearchQuery(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	results := search(req.Query, s.storage, s.avgDocLength, s.collectionSize, s.cache)

	// pagination here
	offset := (req.Page - 1) * req.Count // start index
	// check if offset is out of bounds
	if offset >= int32(len(results)) {
		return &pb.SearchResponse{Results: []*pb.SearchResult{}}, nil
	}
	// check if offset + count is out of bounds
	if offset+req.Count > int32(len(results)) {
		results = results[offset:]
	} else {
		results = results[offset : offset+req.Count]
	}
	docIDs := make([]string, len(results))
	for i, result := range results {
		docIDs[i] = result.Hash
	}
	docs, err := s.storage.Corpus.GetBatchMetadata(ctx, docIDs)
	if err != nil {
		log.Fatalf("failed to get metadata: %v", err)
	}
	finalResults := make([]*pb.SearchResult, len(results))

	for i, result := range results {
		docMetadata := docs[i]
		finalResults[i] = &pb.SearchResult{
			Doc: &pb.DocMetadata{
				Url:   docMetadata.URL,
				Depth: int32(docMetadata.Depth),
				Title: docMetadata.Title,
				Hash:  docMetadata.Hash,
				Links: docMetadata.Links[:3],
			},
			Score:     result.Score,
			TermCount: int32(result.CountTerm),
		}
	}
	fmt.Println("Final results: ", finalResults)

	return &pb.SearchResponse{Results: finalResults}, nil
}

func startServer(storage *shared.Storage) {

	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server is running on port", PORT)
	grpcServer := grpc.NewServer()
	stats, err := storage.GetStats()
	if err != nil {
		log.Fatalf("failed to get stats: %v", err)
	}
	collectionSize := int64(stats.TotalDocs)
	fmt.Println("Avg doc length:", stats.AvgDocLength)
	cache := NewLRUCache[string, []SearchResult](1000)
	pb.RegisterSearchServer(grpcServer, NewSearchServer(storage, stats.AvgDocLength, collectionSize, &cache))
	grpcServer.Serve(lis)
}
