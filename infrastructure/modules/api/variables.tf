variable "db_connection_name" {
  type = string
}

variable "db_user" {
  type = string
}
variable "db_password" {
  type = string
}
variable "db_name" {
  type = string
}

variable "region" {
  type    = string
  default = "europe-west1"
}

variable "domain" {
  type = string
}
