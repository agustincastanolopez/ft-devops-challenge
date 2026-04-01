output "amqp_endpoint" {
  description = "AMQP+TLS endpoint for the RabbitMQ broker"
  value       = "amqps://b-abc12345-6789-0abc-def0-123456789abc.mq.eu-west-1.amazonaws.com:5671"
}

output "web_console_url" {
  description = "RabbitMQ management console URL"
  value       = "https://b-abc12345-6789-0abc-def0-123456789abc.mq.eu-west-1.amazonaws.com:443"
}

output "broker_arn" {
  description = "ARN of the Amazon MQ broker"
  value       = "arn:aws:mq:eu-west-1:123456789012:broker:fasttrack-rabbitmq-${var.environment}:b-abc12345-6789-0abc-def0-123456789abc"
}

output "security_group_id" {
  description = "Security group ID attached to the Amazon MQ broker"
  value       = "sg-0rmq${var.environment}abc123"
}
