package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dxmv/google_clone/pb"
	"google.golang.org/grpc"
)

type DocMetadata struct {
	URL   string
	Depth int
	Title string
	Hash  string
	Links []string
}

const PAGES_DIR = "../crawler/pages"
const METADATA_DIR = "../crawler/metadata"

// Check for errors and exit if they occur
func error_check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// Global inverted index structure
var postings map[string][]Posting

func main1() {
	// Initialize the postings map
	postings = make(map[string][]Posting)

	// Read the metadata directory
	files, err := os.ReadDir(METADATA_DIR)
	error_check(err)

	// Open the Badger database
	db, err := openDB()
	error_check(err)
	defer db.Close()

	// Index each file
	for _, file := range files {
		// open the metadata file
		metadataFilePath := fmt.Sprintf("%s/%s", METADATA_DIR, file.Name())
		metadataFile, err := os.ReadFile(metadataFilePath)
		error_check(err)
		var metadata DocMetadata
		err = json.Unmarshal(metadataFile, &metadata)
		error_check(err)

		// index the html file
		hash := metadata.Hash
		htmlFilePath := fmt.Sprintf("%s/%s", PAGES_DIR, hash+".html")
		index_file(htmlFilePath, []byte(hash), postings)
		fmt.Println("Indexed: ", metadata.Title, "with hash: ", hash)
		// save the metadata to the Badger database
		err = saveMetadata(db, []byte(hash), metadata)
		error_check(err)

	}

	fmt.Printf("Total terms indexed: %d\n", len(postings))

	for term, postingsList := range postings {
		singleTermMap := map[string][]Posting{term: postingsList}
		err = savePostings(db, singleTermMap)
		error_check(err)
	}

	if err != nil {
		fmt.Println("Error saving postings: ", err)
	} else {
		fmt.Println("Saved postings...")
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		results := search(query, db)
		json.NewEncoder(w).Encode(results)
		fmt.Println("Results: ", results)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type searchServer struct {
	pb.UnimplementedSearchServer
}

func (s *searchServer) SearchQuery(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	db, err := openDB()
	error_check(err)
	defer db.Close()

	results := search(req.Query, db)
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
	return &pb.SearchResponse{Results: finalResults}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server is running on port 50051")
	grpcServer := grpc.NewServer()
	pb.RegisterSearchServer(grpcServer, &searchServer{})
	grpcServer.Serve(lis)
}
