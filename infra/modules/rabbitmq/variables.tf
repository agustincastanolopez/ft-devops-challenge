variable "environment" {
  description = "Deployment environment (staging, production)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where Amazon MQ is deployed"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for the Amazon MQ broker"
  type        = list(string)
}
