resource "aws_security_group" "backend-alb-sg" {
  name        = "Backend ALB SG"
  description = "Allow http and https connection"
  vpc_id      = aws_vpc.rifqoi-vpc.id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
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

resource "aws_alb" "backend-alb" {
  name                       = "backend-alb"
  load_balancer_type         = "application"
  internal                   = false
  security_groups            = [aws_security_group.backend-alb-sg.id]
  subnets                    = [aws_subnet.public-subnet-1.id, aws_subnet.public-subnet-2.id]
  enable_deletion_protection = true
}

resource "aws_alb_target_group" "backend-alb-target-group" {
  name     = "backend-alb-tg-http"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.rifqoi-vpc.id
}

/* resource "aws_alb_target_group" "backend-alb-target-group2" { */
/*   name        = "backend-alb-tg-https" */
/*   target_type = "instance" */
/*   port        = 443 */
/*   protocol    = "TCP" */
/*   vpc_id      = aws_vpc.rifqoi-vpc.id */
/* } */

resource "aws_lb_listener" "backend-alb-listener" {
  port              = 80
  load_balancer_arn = aws_alb.backend-alb.arn
  protocol          = "HTTP"

  default_action {
    target_group_arn = aws_alb_target_group.backend-alb-target-group.arn
    type             = "forward"
  }
}
