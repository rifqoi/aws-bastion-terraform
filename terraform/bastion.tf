resource "aws_security_group" "bastion-sg" {
  name        = "Bastion Host SG"
  description = "Allow ssh to connect to instances in the private subnet"
  vpc_id      = aws_vpc.rifqoi-vpc.id

  # Allow to connect to ssh
  ingress {
    from_port   = 22
    to_port     = 22
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
  subnet_id   = aws_subnet.public-subnet-1.id
  private_ips = ["10.0.1.12"]

  tags = {
    Name = "bastion_host_eni"
  }
}


resource "aws_instance" "bastion-host-ec2" {
  ami               = data.aws_ami.amazon-linux-2-ami.id
  instance_type     = var.instance-type
  availability_zone = var.zone
  subnet_id         = aws_subnet.public-subnet-1.id

  network_interface {
    network_interface_id = aws_network_interface.bastion-eni.id
    device_index         = 0
  }

  vpc_security_group_ids = [aws_security_group.bastion-sg.id]
  key_name               = "vockey"
}
