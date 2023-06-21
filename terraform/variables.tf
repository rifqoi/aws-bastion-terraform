variable "region" {
  default     = "us-east-1"
  type        = string
  description = "Default region"
}

variable "zone" {
  default     = "us-east-1a"
  type        = string
  description = "Default AZ"
}

variable "zone2" {
  default     = "us-east-1b"
  type        = string
  description = "Default AZ"
}

variable "instance-type" {
  default     = "t3.micro"
  type        = string
  description = "Default instance type of this project"
}

variable "amazon-linux-2023" {
  default     = "ami-0715c1897453cabd1"
  type        = string
  description = "Amazon Linux 2023 AMI 2023.0.20230517.1 x86_64 HVM kernel-6.1"
}

variable "ubuntu-2004-ami" {
  default     = "ami-0261755bbcb8c4a84"
  type        = string
  description = "Ubuntu 20.04 AMI"
}
