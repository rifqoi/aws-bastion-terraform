resource "aws_security_group" "launch-template-sg" {
  name        = "Launch template security group"
  description = "Allow connection from public subnet to asg"
  vpc_id      = aws_vpc.rifqoi-vpc.id

  # Allow to connect to ssh
  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.bastion-sg.id]
    cidr_blocks     = [aws_subnet.private-subnet-1.cidr_block]
    /* cidr_blocks = [aws_security_group.bastion-sg] */
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    /* cidr_blocks = [aws_security_group.bastion-sg] */
  }

  ingress {
    from_port       = 8000
    to_port         = 8000
    protocol        = "tcp"
    security_groups = [aws_security_group.bastion-sg.id]
    cidr_blocks     = [aws_subnet.private-subnet-1.cidr_block]
    /* cidr_blocks = [aws_security_group.bastion-sg] */
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
    Name = "launch-template-sg"
  }
}

resource "aws_launch_template" "backend-launch-template" {
  name                   = "backend-launch-template"
  description            = "Launch template for backend services"
  instance_type          = var.instance-type
  image_id               = data.aws_ami.amazon-linux-2-ami.id
  key_name               = "vockey"
  update_default_version = true
  user_data              = filebase64("./scripts/docker.sh")

  network_interfaces {
    security_groups = [aws_security_group.launch-template-sg.id]
  }

  block_device_mappings {
    device_name = "/dev/sda1"
    ebs {
      volume_size           = 10
      delete_on_termination = true
      volume_type           = "gp2"
    }
  }

  tag_specifications {
    resource_type = "instance"
    tags = {
      "Name" : "backend-launch-template"
    }
  }
}

resource "aws_autoscaling_group" "backend-asg" {
  name_prefix         = "backend-asg-"
  min_size            = 2
  max_size            = 3
  desired_capacity    = 2
  vpc_zone_identifier = [aws_subnet.private-subnet-1.id]
  health_check_type   = "ELB"
  target_group_arns   = [aws_alb_target_group.backend-alb-target-group.arn]

  depends_on = [aws_alb.backend-alb]

  launch_template {
    id      = aws_launch_template.backend-launch-template.id
    version = aws_launch_template.backend-launch-template.latest_version
  }

  /* instance_refresh { */
  /*   strategy = "Rolling" */
  /*   triggers = ["launch_template"] */
  /* } */

  tag {
    key                 = "Service"
    value               = "Backend"
    propagate_at_launch = true
  }
}

# Autoscaling Notifications
## SNS - Topic
resource "aws_sns_topic" "asg_sns_topic" {
  name = "asg-sns-topic"
}

## SNS - Subscription
resource "aws_sns_topic_subscription" "asg_sns_topic_subscription" {
  topic_arn = aws_sns_topic.asg_sns_topic.arn
  protocol  = "email"
  endpoint  = "alfurqon.rifqi@gmail.com"
}

## Create Autoscaling Notification Resource
resource "aws_autoscaling_notification" "asg_notifications" {
  group_names = [aws_autoscaling_group.backend-asg.id]
  notifications = [
    "autoscaling:EC2_INSTANCE_LAUNCH",
    "autoscaling:EC2_INSTANCE_TERMINATE",
    "autoscaling:EC2_INSTANCE_LAUNCH_ERROR",
    "autoscaling:EC2_INSTANCE_TERMINATE_ERROR",
  ]
  topic_arn = aws_sns_topic.asg_sns_topic.arn
}

###### Target Tracking Scaling Policies ######
# TTS - Scaling Policy-1: Based on CPU Utilization
# Define Autoscaling Policies and Associate them to Autoscaling Group
resource "aws_autoscaling_policy" "avg_cpu_policy_greater_than_xx" {
  name                      = "avg-cpu-policy-greater-than-xx"
  policy_type               = "TargetTrackingScaling" # Important Note: The policy type, either "SimpleScaling", "StepScaling" or "TargetTrackingScaling". If this value isn't provided, AWS will default to "SimpleScaling."    
  autoscaling_group_name    = aws_autoscaling_group.backend-asg.id
  estimated_instance_warmup = 180 # defaults to ASG default cooldown 300 seconds if not set
  # CPU Utilization is above 50
  target_tracking_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ASGAverageCPUUtilization"
    }
    target_value = 50.0
  }

}
