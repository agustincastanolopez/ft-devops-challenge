# Fast Track — Senior DevOps Engineer Take-Home Challenge

Starter repository for the event-enricher microservice. Contains a compiling Go stub with health endpoints, shared Terraform module stubs for existing infrastructure, and empty directories for your deliverables.

## Getting Started

```bash
cd services/event-enricher
go test ./...
go build -o event-enricher .
./event-enricher
# curl http://localhost:8080/healthz  -> {"status":"ok"}
# curl http://localhost:8080/readyz   -> {"status":"ready"}
```

Full challenge instructions will be provided separately.
