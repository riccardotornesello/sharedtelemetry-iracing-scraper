variable "name" {
  type = string
}

variable "short_name" {
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

variable "pubsub_client" {
  type    = bool
  default = false
}

variable "ack_deadline_seconds" {
  type    = number
  default = 300
}

variable "timeout" {
  type    = string
  default = "600s"
}

variable "cron" {
  type    = string
  default = null
}
