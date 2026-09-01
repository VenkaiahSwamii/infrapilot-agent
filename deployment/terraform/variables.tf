variable "aws_region" {
  description = "AWS region for InfraPilot deployment"
  default     = "us-east-1"
}

variable "environment" {
  description = "Target deployment environment"
  default     = "production"
}

variable "instance_type" {
  description = "EC2 instance size"
  default     = "t3.medium"
}

variable "ami_id" {
  description = "AMI ID for Linux server"
  default     = "ami-0c7217cdde317cfec"
}
