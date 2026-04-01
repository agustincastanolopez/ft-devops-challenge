# This module represents the existing Amazon MQ (RabbitMQ) broker managed by the platform team.
# Do NOT recreate this resource — consume the outputs below.
#
# Engine: RabbitMQ 3.13, mq.m5.large, single-instance (staging) / active-standby (production), eu-west-1

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# The actual aws_mq_broker resource is managed elsewhere.
# This module only exposes outputs for downstream consumers.
