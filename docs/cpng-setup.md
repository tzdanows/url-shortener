# CloudNative PostgreSQL (CNPG) URL Shortener Setup Guide

## Introduction
This guide details setting up a highly available PostgreSQL cluster on Kubernetes using CloudNative PostgreSQL (CNPG) for a URL shortening service. The database will store mappings between original URLs and their shortened versions.

## Table of Contents
1. Install Helm
2. Add CNPG Helm Chart Repository
3. Install the CNPG Operator
4. Deploy a CNPG Cluster
5. Connect to the CNPG Cluster
6. Create the URLs Table
7. Implement the API logic

## Step-by-Step Guide

### 0. Ensure
```bash
minikube status
# ELSE
minikube start
kubectl get pods
```

### 1. Install Helm
```bash
brew install helm
```

### 2. Add CNPG Helm Chart Repository
```bash
helm repo add cnpg https://cloudnative-pg.github.io/charts
helm repo update
```

### 3. Install the CNPG Operator
```bash
helm install cnpg cnpg/cloudnative-pg --version 0.22.0
```

Verify the operator deployment:
```bash
kubectl get pods -l app.kubernetes.io/name=cloudnative-pg
```

### 4. Deploy a CNPG Cluster
Create & apply `postgres-cluster.yaml`:
  
Apply and verify:
```bash
kubectl apply -f postgres-cluster.yaml

# in a new tab, watch the pods get ready
kubectl get pods -w

# return to the former tab
kubectl get clusters.postgresql.cnpg.io
kubectl get pods -l postgresql.cnpg.io/cluster=urlshortener-db
```

### 5. Connect to the CNPG Cluster
CNPG creates several services:
- `urlshortener-db-rw`: Primary node (read-write)
- `urlshortener-db-ro`: Replicas (read-only)
- `urlshortener-db-r`: All replicas

Connect to primary node:
```bash
# Get the password from CNPG-managed secret
kubectl get secret urlshortener-db-app -o jsonpath="{.data.password}" | base64 -d

# Port forward the read-write service
kubectl port-forward svc/urlshortener-db-rw 5432:5432

# Connect (in new terminal)
PGPASSWORD=$(kubectl get secret urlshortener-db-app -o jsonpath="{.data.password}" | base64 -d) \
psql -h localhost -p 5432 -U app urlshortenerdb
```

### 6. Create the URLs Table
```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_url VARCHAR(50) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMPTZ,
    access_count INTEGER DEFAULT 0
);

CREATE INDEX idx_short_url ON urls(short_url);
CREATE INDEX idx_created_at ON urls(created_at);
```

### 7. Implement a non-mock url-shortener client/server

Currently, this implementation uses simplified placeholder values:
- Long URL: Always set to "foo" (simulating an input URL)
- Short URL: Always set to "bar" (simulating a generated short code)

These hardcoded values will be replaced in future iterations with:
- Long URL: Actual URLs provided by users via gRPC requests
- Short URL: Generated unique short codes using a proper URL shortening algorithm

The database structure is already set up to handle the full implementation when ready.

Test queries (must be run on primary/-rw service):
```sql
-- Insert test URL
INSERT INTO urls (short_url, long_url) 
VALUES ('test123', 'https://example.com/long/path');

-- Query from primary or replicas
SELECT * FROM urls WHERE short_url = 'test123';

-- Updates must go to primary
UPDATE urls 
SET last_accessed = CURRENT_TIMESTAMP, 
    access_count = access_count + 1 
WHERE short_url = 'test123';
```

### cnpg + kubernetes features
- 3 PostgreSQL instances running
- Automatic failover if primary fails
- Read scalability using -ro service
- Consistent reads from primary using -rw service

### Common Operations
```bash
# Check cluster status
kubectl get clusters.postgresql.cnpg.io

# View pods in cluster
kubectl get pods -l postgresql.cnpg.io/cluster=urlshortener-db

# Check services
kubectl get svc -l postgresql.cnpg.io/cluster=urlshortener-db

# View logs
kubectl logs -l postgresql.cnpg.io/cluster=urlshortener-db
```