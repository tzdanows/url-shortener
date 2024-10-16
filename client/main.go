package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"buf.build/gen/go/tzdanows/url-shortener/connectrpc/go/urlshortener/v1/urlshortenerv1connect"
	urlshortenerv1 "buf.build/gen/go/tzdanows/url-shortener/protocolbuffers/go/urlshortener/v1"
	"connectrpc.com/connect"
)

func main() {
	client := urlshortenerv1connect.NewURLShortenerServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	// Test ShortenURL
	shortenResp, err := client.ShortenURL(
		context.Background(),
		connect.NewRequest(&urlshortenerv1.ShortenURLRequest{LongUrl: "https://example.com/very/long/url"}),
	)
	if err != nil {
		log.Fatalf("ShortenURL error: %v", err)
	}
	fmt.Printf("Shortened URL: %s\n", shortenResp.Msg.ShortUrl)

	// Test GetOriginalURL
	getOriginalResp, err := client.GetOriginalURL(
		context.Background(),
		connect.NewRequest(&urlshortenerv1.GetOriginalURLRequest{ShortUrl: "abc123"}),
	)
	if err != nil {
		log.Fatalf("GetOriginalURL error: %v", err)
	}
	fmt.Printf("Original URL: %s\n", getOriginalResp.Msg.LongUrl)
}

