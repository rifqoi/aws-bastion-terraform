locals {
  SSH_PORT = 2225
}

resource "aws_security_group" "bastion-sg" {
  name        = "Bastion Host SG"
  description = "Allow ssh to connect to instances in the private subnet"
  vpc_id      = aws_vpc.rifqoi-vpc.id

  # Allow to connect to ssh
  # Change ssh port to other than 22 for security reason
  # In this case I changed it to 2222
  ingress {
    from_port   = local.SSH_PORT
    to_port     = local.SSH_PORT
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # This is default outbound rules to ALLOW ALL connection
  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "bastion-sg"
  }
}

resource "aws_eip" "bastion-eip" {
  instance = aws_instance.bastion-host-ec2.id
  vpc      = true
}

data "aws_ami" "amazon-linux-2-ami" {
  owners      = ["amazon"]
  most_recent = true

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

}

resource "aws_network_interface" "bastion-eni" {
  subnet_id       = aws_subnet.public-subnet-1.id
  private_ips     = ["10.0.1.12"]
  security_groups = [aws_security_group.bastion-sg.id]

  tags = {
    Name = "bastion_host_eni"
  }
}


resource "aws_instance" "bastion-host-ec2" {
  ami               = data.aws_ami.amazon-linux-2-ami.id
  instance_type     = var.instance-type
  availability_zone = var.zone
  /* subnet_id         = aws_subnet.public-subnet-1.id */

  // This is the only instance profile 
  // that can be used in AWS Academy Sandbox...
  // I was wrong...
  /* iam_instance_profile = "EMR_EC2_DefaultRole" */


  network_interface {
    network_interface_id = aws_network_interface.bastion-eni.id
    device_index         = 0
  }

  user_data = <<EOF
    #!/bin/bash
    # Change SSH PORT from 22 to 2225
    sudo sed -i '/Port 22/a Port ${local.SSH_PORT}' /etc/ssh/sshd_config

    # Restart sshd daemon
    sudo service sshd restart

    # Install Docker
    sudo yum update
    sudo yum install docker
    sudo usermod -a -G docker ec2-user
    id ec2-user
    # Reload a Linux user's group assignments to docker w/o logout
    newgrp docker
    sudo systemctl enable docker.service
    sudo systemctl start docker.service
    EOF

  key_name = "vockey"
  tags = {
    "Name" = "bastion-instance"
  }
}
