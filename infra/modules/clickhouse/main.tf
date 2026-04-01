# This module represents the existing ClickHouse cluster managed by the data team.
# Do NOT recreate this resource — consume the outputs below.
#
# Deployment: ClickHouse on EC2 (3-node cluster), behind an internal NLB, eu-west-1
# The event-enricher is expected to write via the HTTP interface (port 8123).

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# The actual ClickHouse infrastructure is managed elsewhere.
# This module only exposes outputs for downstream consumers.
