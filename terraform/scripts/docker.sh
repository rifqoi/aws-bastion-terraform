#!/bin/bash
# Install Docker
sudo yum update
sudo yum install docker golang-go -y
sudo usermod -a -G docker ec2-user

# Reload a Linux user's group assignments to docker w/o logout
newgrp docker
sudo systemctl enable docker.service
sudo systemctl start docker.service

docker run -d --rm -p 80:8000 rifqoi/todo-backend
