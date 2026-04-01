# This module represents the existing ElastiCache Redis cluster managed by the platform team.
# Do NOT recreate this resource — consume the outputs below.
#
# Engine: Redis 7.x, cluster mode disabled, r6g.large, encryption in transit, eu-west-1

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# The actual aws_elasticache_replication_group resource is managed elsewhere.
# This module only exposes outputs for downstream consumers.
