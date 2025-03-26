package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	pb "github.com/khzs/connectivity-listener/proto/ping"
	"google.golang.org/grpc"
)

// Environment variables for configuration
const (
	defaultHTTPPort = "8080"
	defaultGRPCPort = "50051"
	envHTTPPort     = "HTTP_PORT"
	envGRPCPort     = "GRPC_PORT"
)

type server struct {
	pb.UnimplementedPingServiceServer
}

// Ping implements the gRPC Ping method
func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	sum := in.A + in.B
	hash := sha1.Sum([]byte(strconv.Itoa(int(sum))))
	return &pb.PingResponse{Hash: hex.EncodeToString(hash[:])}, nil
}

// pingHandler handles HTTP ping requests
func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/ping" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Parse timestamp from request
	var timestamp time.Time
	err := timestamp.UnmarshalText([]byte(r.Header.Get("X-Timestamp")))
	if err != nil {
		http.Error(w, "Invalid timestamp format", http.StatusBadRequest)
		return
	}

	// Check if timestamp is in the future
	if timestamp.Before(time.Now()) {
		http.Error(w, "Timestamp must be in the future", http.StatusBadRequest)
		return
	}

	// Return 200 OK if everything is fine
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Get port configuration from environment variables
	httpPort := os.Getenv(envHTTPPort)
	if httpPort == "" {
		httpPort = defaultHTTPPort
	}

	grpcPort := os.Getenv(envGRPCPort)
	if grpcPort == "" {
		grpcPort = defaultGRPCPort
	}

	// Start HTTP server
	go func() {
		http.HandleFunc("/ping", pingHandler)
		log.Printf("Starting HTTP server on port %s", httpPort)
		if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPingServiceServer(s, &server{})
	log.Printf("Starting gRPC server on port %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
