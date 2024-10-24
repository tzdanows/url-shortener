# url-shortener
a url shortener service built to learn about buf, cpng, & grpc integration

## Quickstart

Create a cluster & run the application:

```bash
minikube start
skaffold dev
```

Shorten a URL:

```bash
curl -X POST http://localhost:8080/urlshortener.v1.URLShortenerService/ShortenURL \
     -H "Content-Type: application/json" \
     -d '{"long_url": "https://example.com/very/long/url"}'
```

Get the original URL:

```bash
curl -X POST http://localhost:8080/urlshortener.v1.URLShortenerService/GetOriginalURL \
     -H "Content-Type: application/json" \
     -d '{"short_url": "abc123"}'
```

## Technologies used:

- Go
- Protobuf
- Buf
- Buf Schema Registry
- gRPC
- Docker
- Kubernetes
- Skaffold OR Tilt
- Minikube
- Postgres (cnpg and pgx)