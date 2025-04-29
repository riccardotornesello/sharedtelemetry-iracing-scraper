variable "region" {
  description = "The Google Cloud region where the function will be deployed."
  type        = string
  default     = "europe-west1"
}

variable "iracing_email" {
  description = "The email address used to authenticate with the iRacing API."
  type        = string
}

variable "iracing_password" {
  description = "The password used to authenticate with the iRacing API."
  type        = string
}
