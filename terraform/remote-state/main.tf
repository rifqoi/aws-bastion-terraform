terraform {
  backend "local" {
    path = "terraform.tfstate"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>4.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "aws-proj-remote-state" {
  bucket              = "aws-proj-remote-state"
  object_lock_enabled = false
  tags = {
    Name        = "TF State Bucket"
    Environment = "DEV"
  }
}

resource "aws_s3_bucket_versioning" "remote-state-versioning" {
  bucket = aws_s3_bucket.aws-proj-remote-state.id
  versioning_configuration {
    status = "Enabled"
  }
}
