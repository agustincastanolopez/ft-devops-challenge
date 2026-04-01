# This module represents the existing EKS cluster and core node groups managed by the platform team.
# Do NOT recreate this resource — consume the outputs below.
#
# Cluster: EKS 1.29, managed node groups (m5.xlarge), eu-west-1
# OIDC provider is enabled for IRSA (IAM Roles for Service Accounts).

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# The actual aws_eks_cluster and aws_eks_node_group resources are managed elsewhere.
# This module only exposes outputs for downstream consumers.
