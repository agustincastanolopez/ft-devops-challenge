output "primary_endpoint" {
  description = "Primary endpoint for the Redis replication group"
  value       = "fasttrack-redis-${var.environment}.abc123.0001.euw1.cache.amazonaws.com"
}

output "port" {
  description = "Redis port (TLS)"
  value       = 6379
}

output "security_group_id" {
  description = "Security group ID attached to the ElastiCache cluster"
  value       = "sg-0redis${var.environment}abc123"
}

output "arn" {
  description = "ARN of the ElastiCache replication group"
  value       = "arn:aws:elasticache:eu-west-1:123456789012:replicationgroup:fasttrack-redis-${var.environment}"
}
