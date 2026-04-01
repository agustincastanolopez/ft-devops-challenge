# This module represents the existing MSK cluster managed by the platform team.
# Do NOT recreate this resource — consume the outputs below.
#
# Cluster: 3 brokers, kafka.m5.large, TLS enabled, eu-west-1

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# The actual aws_msk_cluster resource is managed elsewhere.
# This module only exposes outputs for downstream consumers.
