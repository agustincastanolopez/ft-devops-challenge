# This module represents the existing Aurora MySQL cluster managed by the platform team.
# Do NOT recreate this resource — consume the outputs below.
#
# Engine: Aurora MySQL 8.0, Multi-AZ, r6g.large instances, eu-west-1

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# The actual aws_rds_cluster resource is managed elsewhere.
# This module only exposes outputs for downstream consumers.
