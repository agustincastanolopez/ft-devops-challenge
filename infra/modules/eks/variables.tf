variable "environment" {
  description = "Deployment environment (staging, production)"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID where the EKS cluster is deployed"
  type        = string
}

variable "subnet_ids" {
  description = "Subnet IDs for EKS node groups"
  type        = list(string)
}
