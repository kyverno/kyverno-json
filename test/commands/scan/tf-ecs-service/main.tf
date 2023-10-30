provider "aws" {
  region                      = "eu-west-1"
  skip_credentials_validation = true
  skip_requesting_account_id  = true
  skip_metadata_api_check     = true
  access_key                  = "mock_access_key"
  secret_key                  = "mock_secret_key"
}

resource "aws_ecs_service" "mongo" {
  name            = "mongodb"
  cluster         = "id"
  task_definition = "aws_arn"
  desired_count   = 3

  ordered_placement_strategy {
    type  = "binpack"
    field = "cpu"
  }

  load_balancer {
    container_name   = "mongo"
    container_port   = 8080
  }

  network_configuration {
    subnets = ["subnet-abcde012", "subnet-bcde012a", "subnet-fghi345a"]
    assign_public_ip = true
  }

  launch_type = "FARGATE"
  platform_version = "LATEST"

  placement_constraints {
    type       = "memberOf"
    expression = "attribute:ecs.availability-zone in [us-west-2a, us-west-2b]"
  }
}