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
	// Use the more efficient paginated search function
	results, totalResults := searchPaginated(req.Query, s.storage, s.avgDocLength, s.collectionSize, s.cache, req.Page, req.Count)
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
				Url:            docMetadata.URL,
				FirstParagraph: docMetadata.FirstParagraph,
				Depth:          int32(docMetadata.Depth),
				Title:          docMetadata.Title,
				Hash:           docMetadata.Hash,
				Images:         docMetadata.Images,
			},
			Score:     result.Score,
			TermCount: int32(result.CountTerm),
		}
	}

	return &pb.SearchResponse{Results: finalResults, Total: totalResults}, nil
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
