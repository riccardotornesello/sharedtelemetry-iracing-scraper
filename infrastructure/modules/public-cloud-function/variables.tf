variable "region" {
  description = "The region to deploy the resources"
  type        = string
}

variable "project" {
  description = "The project name"
  type        = string
}

variable "source_archive" {
  description = "The source archive file path"
  type        = string
}

variable "name" {
  description = "The name of the function"
  type        = string
}

variable "short_name" {
  description = "The short name of the function"
  type        = string
}

variable "entry_point" {
  description = "The entry point of the function"
  type        = string
}

variable "runtime" {
  description = "The runtime of the function"
  type        = string
}

variable "environment_variables" {
  description = "The environment variables of the function"
  type        = map(string)
  default     = {}
}

variable "domain" {
  description = "The domain name"
  type        = string
}

variable "roles" {
  description = "The roles of the function"
  type        = list(string)
  default     = []
}
