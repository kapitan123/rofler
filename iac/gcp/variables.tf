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

variable "name" {
  type = string
}

variable "port" {
  type    = number
  default = 9001
}
