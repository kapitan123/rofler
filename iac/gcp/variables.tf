variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "europe-central2"
}

variable "bot_token" {
  type      = string
  sensitive = true
}

variable "registry_id" {
  type    = string
  default = "eu.gcr.io"
}

variable "bot_name" {
  type = string
}

variable "port" {
  type    = number
  default = 9001
}

variable "video_convertion_topic_name" {
  type    = string
  default = "video_convertion"
}

variable "deletion_queue_name" {
  type    = string
  default = "message-deletion"
}