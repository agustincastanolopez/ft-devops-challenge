variable "environment" {
  description = "Deployment environment (staging, production)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where ClickHouse is deployed"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for the ClickHouse instances"
  type        = list(string)
}
