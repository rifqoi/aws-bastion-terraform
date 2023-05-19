resource "aws_dynamodb_table" "todo-table" {
  name           = "todo"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }

  tags = {
    Name        = "todo-table"
    Environment = "development"
  }

}
