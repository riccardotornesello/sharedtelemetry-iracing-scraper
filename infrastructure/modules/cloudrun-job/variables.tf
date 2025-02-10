variable "env" {
  type = map(string)
}

variable "region" {
  type = string
}

variable "project" {
  type = string
}

variable "project_number" {
  type = number
}

variable "db_connection_name" {
  type = string
}

variable "name" {
  type = string
}

variable "short_name" {
  type = string
}

variable "image" {
  type = string
}

variable "run_after_deploy" {
  type    = bool
  default = false
}

variable "args" {
  type    = list(string)
  default = null
}
