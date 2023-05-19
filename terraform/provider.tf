terraform {
  backend "s3" {
    bucket = "aws-proj-remote-state2"
    region = "us-east-1"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>4.0"
    }
  }
}

provider "aws" {
  region = var.region
}
