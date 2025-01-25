variable "db_connection_name" {
  type = string
}

variable "events_db_user" {
  type = string
}
variable "events_db_password" {
  type = string
}
variable "events_db_name" {
  type = string
}

variable "drivers_db_user" {
  type = string
}
variable "drivers_db_password" {
  type = string
}
variable "drivers_db_name" {
  type = string
}

variable "region" {
  type    = string
  default = "europe-west1"
}

variable "domain" {
  type = string
}
