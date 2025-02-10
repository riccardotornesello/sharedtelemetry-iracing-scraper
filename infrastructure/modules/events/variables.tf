variable "db_instance_name" {
  type = string
}

variable "db_connection_name" {
  type = string
}

variable "db_password" {
  type = string
}

variable "iracing_email" {
  type = string
}

variable "iracing_password" {
  type = string
}

variable "region" {
  type    = string
  default = "europe-west1"
}

variable "project" {
  type = string
}

variable "project_number" {
  type = string
}
