variable "function_name" {
  description = "The name of the function"
  type        = string
}

variable "location" {
  description = "The location of the function"
  type        = string
}

variable "runtime" {
  description = "The runtime of the function"
  type        = string
  default     = "go123"
}

variable "entrypoint" {
  description = "The entry point of the function"
  type        = string
  default     = "Handler"
}

variable "environment_variables" {
  description = "Environment variables for the function"
  type        = map(string)
  default     = {}
}

variable "source_dir" {
  description = "The source directory of the function"
  type        = string
}

variable "pubsub_topic_id" {
  description = "The ID of the Pub/Sub topic"
  type        = string
}

variable "cron_schedule" {
  description = "The cron schedule for the function"
  type = list(object({
    schedule = string
    payload  = string
  }))
  default = []
} 
