variable "name" {
  type = string
}

variable "location" {
  type = string
}

variable "project" {
  type = string
}

variable "max_instance_count" {
  type    = number
  default = 1
}

variable "image" {
  type = string
}

variable "env" {
  type = map(string)
}

variable "db_connection_name" {
  type    = string
  default = null
}
