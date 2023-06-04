## VPC
resource "aws_vpc" "rifqoi-vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "rifqoi-vpc"
  }
}
#################################

## Public subnet 1
resource "aws_subnet" "public-subnet-1" {
  vpc_id            = aws_vpc.rifqoi-vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = var.zone

  tags = {
    Name = "Public subnet 1"
  }
}

## Internet gateway
resource "aws_internet_gateway" "igw-rifqoi" {
  vpc_id = aws_vpc.rifqoi-vpc.id

  tags = {
    Name = "igw-rifqoi"
  }
}

## Route table for public subnet
resource "aws_route_table" "public-rtb-1" {
  vpc_id = aws_vpc.rifqoi-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw-rifqoi.id
  }

  tags = {
    Name = "public-rtb-1"
  }
}

resource "aws_route_table_association" "public-rtb-1-assoc" {
  route_table_id = aws_route_table.public-rtb-1.id
  subnet_id      = aws_subnet.public-subnet-1.id
}

## NAT Gateway
resource "aws_eip" "ngw-eip-1" {
  vpc = true

  tags = {
    Name = "ngw-eip-1"
  }
}

resource "aws_nat_gateway" "ngw-vpc-rifqoi" {
  subnet_id     = aws_subnet.public-subnet-1.id
  allocation_id = aws_eip.ngw-eip-1.id

  tags = {
    Name = "ngw-public-1"
  }

  depends_on = [aws_internet_gateway.igw-rifqoi]
}

#################################

## Public subnet 1
resource "aws_subnet" "public-subnet-2" {
  vpc_id            = aws_vpc.rifqoi-vpc.id
  cidr_block        = "10.0.4.0/24"
  availability_zone = var.zone2

  tags = {
    Name = "Public subnet 2"
  }
}

## Route table for public subnet
resource "aws_route_table" "public-rtb-2" {
  vpc_id = aws_vpc.rifqoi-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw-rifqoi.id
  }

  tags = {
    Name = "public-rtb-2"
  }
}

resource "aws_route_table_association" "public-rtb-2-assoc" {
  route_table_id = aws_route_table.public-rtb-2.id
  subnet_id      = aws_subnet.public-subnet-2.id
}

## NAT Gateway
resource "aws_eip" "ngw-eip-2" {
  vpc = true

  tags = {
    Name = "ngw-eip-2"
  }
}

resource "aws_nat_gateway" "ngw-vpc-rifqoi-2" {
  subnet_id     = aws_subnet.public-subnet-2.id
  allocation_id = aws_eip.ngw-eip-2.id

  tags = {
    Name = "ngw-public-2"
  }

  depends_on = [aws_internet_gateway.igw-rifqoi]
}

#################################

## Private subnet 1
resource "aws_subnet" "private-subnet-1" {
  vpc_id            = aws_vpc.rifqoi-vpc.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = var.zone

  tags = {
    Name = "Private subnet 1"
  }
}

resource "aws_route_table" "private-rtb-1" {
  vpc_id = aws_vpc.rifqoi-vpc.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.ngw-vpc-rifqoi.id
  }

  tags = {
    Name = "private-rtb-1"
  }
}

resource "aws_route_table_association" "private-rtb-1-assoc" {
  subnet_id      = aws_subnet.private-subnet-1.id
  route_table_id = aws_route_table.private-rtb-1.id
}

#################################

## Private subnet 2
resource "aws_subnet" "private-subnet-2" {
  vpc_id            = aws_vpc.rifqoi-vpc.id
  cidr_block        = "10.0.3.0/24"
  availability_zone = var.zone

  tags = {
    Name = "Private subnet 2"
  }
}

resource "aws_route_table" "private-rtb-2" {
  vpc_id = aws_vpc.rifqoi-vpc.id

  tags = {
    Name = "private-rtb-1"
  }
}

resource "aws_route_table_association" "private-rtb-2-assoc" {
  subnet_id      = aws_subnet.private-subnet-2.id
  route_table_id = aws_route_table.private-rtb-2.id
}
