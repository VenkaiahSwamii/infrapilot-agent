output "cluster_node_ips" {
  description = "Public IP addresses of InfraPilot cluster nodes"
  value       = aws_instance.infrapilot_server[*].public_ip
}

output "vpc_id" {
  description = "VPC ID of InfraPilot deployment"
  value       = aws_vpc.infrapilot_vpc.id
}
