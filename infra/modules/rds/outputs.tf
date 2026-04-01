output "cluster_endpoint" {
  description = "Writer endpoint for the Aurora MySQL cluster"
  value       = "fasttrack-aurora-${var.environment}.cluster-abc123def456.eu-west-1.rds.amazonaws.com"
}

output "reader_endpoint" {
  description = "Reader endpoint for the Aurora MySQL cluster"
  value       = "fasttrack-aurora-${var.environment}.cluster-ro-abc123def456.eu-west-1.rds.amazonaws.com"
}

output "port" {
  description = "MySQL port"
  value       = 3306
}

output "database_name" {
  description = "Default database name"
  value       = "fasttrack_players"
}

output "security_group_id" {
  description = "Security group ID attached to the RDS cluster"
  value       = "sg-0rds${var.environment}abc123"
}

output "cluster_arn" {
  description = "ARN of the Aurora cluster"
  value       = "arn:aws:rds:eu-west-1:123456789012:cluster:fasttrack-aurora-${var.environment}"
}
