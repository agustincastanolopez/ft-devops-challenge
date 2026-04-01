output "http_endpoint" {
  description = "ClickHouse HTTP interface endpoint (load balancer)"
  value       = "http://clickhouse-internal-${var.environment}.fasttrack.internal:8123"
}

output "native_endpoint" {
  description = "ClickHouse native protocol endpoint"
  value       = "clickhouse-internal-${var.environment}.fasttrack.internal:9000"
}

output "database_name" {
  description = "Target database for analytics events"
  value       = "analytics"
}

output "security_group_id" {
  description = "Security group ID for the ClickHouse cluster"
  value       = "sg-0ch${var.environment}abc123"
}
