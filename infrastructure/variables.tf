variable "iracing_email" {
  type = string
}

variable "iracing_password" {
  type = string
}

variable "db_password" {
  type = string
}

variable "db_whitelist" {
  type    = list(string)
  default = []
}

variable "region" {
  type    = string
  default = "europe-west1"
}

variable "domain" {
  type = string
}

variable "project" {
  type = string
}

variable "project_number" {
  type = number
}
