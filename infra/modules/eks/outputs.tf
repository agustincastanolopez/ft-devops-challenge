output "cluster_name" {
  description = "Name of the EKS cluster"
  value       = "fasttrack-eks-${var.environment}"
}

output "cluster_endpoint" {
  description = "API server endpoint"
  value       = "https://ABC123DEF456.gr7.eu-west-1.eks.amazonaws.com"
}

output "cluster_certificate_authority" {
  description = "Base64 encoded CA certificate for the cluster"
  value       = "LS0tLS1CRUdJTi..."
}

output "oidc_provider_arn" {
  description = "ARN of the OIDC provider for IRSA"
  value       = "arn:aws:iam::123456789012:oidc-provider/oidc.eks.eu-west-1.amazonaws.com/id/ABC123DEF456"
}

output "oidc_provider_url" {
  description = "URL of the OIDC provider (without https://)"
  value       = "oidc.eks.eu-west-1.amazonaws.com/id/ABC123DEF456"
}

output "node_security_group_id" {
  description = "Security group ID attached to EKS worker nodes"
  value       = "sg-0eks${var.environment}nodes123"
}

output "cluster_security_group_id" {
  description = "Security group ID for the EKS control plane"
  value       = "sg-0eks${var.environment}ctrl123"
}
