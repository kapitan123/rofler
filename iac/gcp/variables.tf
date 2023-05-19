variable "project_id" {
  type    = string
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

variable "port" {
  type    = number
  default = 9001
}

variable "message_deletion_queue_name" {
  type    = string
  default = "tg-message-deletion"
}

variable "bot_name" {
  type    = string
  default = "bot"
}

variable "convertor_name" {
  type    = string
  default = "convertor"
}