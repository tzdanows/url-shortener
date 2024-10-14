package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"buf.build/gen/go/tzdanows/url-shortener/connectrpc/go/urlshortener/v1/urlshortenerv1connect"
	urlshortenerv1 "buf.build/gen/go/tzdanows/url-shortener/protocolbuffers/go/urlshortener/v1"
	"connectrpc.com/connect"
)

// URLShortenerServiceServer implements the URLShortenerService interface.
type URLShortenerServiceServer struct{}

// NewURLShortenerServiceServer creates a new instance of URLShortenerServiceServer.
func NewURLShortenerServiceServer() *URLShortenerServiceServer {
	return &URLShortenerServiceServer{}
}

// ShortenURL handles the ShortenURL RPC.
func (s *URLShortenerServiceServer) ShortenURL(ctx context.Context, req *connect.Request[urlshortenerv1.ShortenURLRequest]) (*connect.Response[urlshortenerv1.ShortenURLResponse], error) {
	return connect.NewResponse(&urlshortenerv1.ShortenURLResponse{ShortUrl: "bar"}), nil
}

// GetOriginalURL handles the GetOriginalURL RPC.
func (s *URLShortenerServiceServer) GetOriginalURL(ctx context.Context, req *connect.Request[urlshortenerv1.GetOriginalURLRequest]) (*connect.Response[urlshortenerv1.GetOriginalURLResponse], error) {
	return connect.NewResponse(&urlshortenerv1.GetOriginalURLResponse{LongUrl: "foo"}), nil
}

func (s *URLShortenerServiceServer) generateShortURL() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func main() {
	mux := http.NewServeMux()
	server := NewURLShortenerServiceServer()

	// Register the URLShortenerService with the Connect server.
	path, handler := urlshortenerv1connect.NewURLShortenerServiceHandler(server)
	mux.Handle(path, handler)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
