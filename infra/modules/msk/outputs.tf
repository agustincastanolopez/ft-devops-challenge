output "bootstrap_brokers_tls" {
  description = "TLS bootstrap broker connection string"
  value       = "b-1.fasttrack-msk-${var.environment}.abc123.c2.kafka.eu-west-1.amazonaws.com:9094,b-2.fasttrack-msk-${var.environment}.abc123.c2.kafka.eu-west-1.amazonaws.com:9094,b-3.fasttrack-msk-${var.environment}.abc123.c2.kafka.eu-west-1.amazonaws.com:9094"
}

output "cluster_arn" {
  description = "ARN of the MSK cluster"
  value       = "arn:aws:kafka:eu-west-1:123456789012:cluster/fasttrack-msk-${var.environment}/abc12345-6789-0abc-def0-123456789abc-1"
}

output "cluster_name" {
  description = "Name of the MSK cluster"
  value       = "fasttrack-msk-${var.environment}"
}

output "zookeeper_connect" {
  description = "ZooKeeper connection string (for topic management)"
  value       = "z-1.fasttrack-msk-${var.environment}.abc123.c2.kafka.eu-west-1.amazonaws.com:2181,z-2.fasttrack-msk-${var.environment}.abc123.c2.kafka.eu-west-1.amazonaws.com:2181,z-3.fasttrack-msk-${var.environment}.abc123.c2.kafka.eu-west-1.amazonaws.com:2181"
}

output "security_group_id" {
  description = "Security group ID attached to MSK brokers"
  value       = "sg-0msk${var.environment}abc123"
}
