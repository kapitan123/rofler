variable "project_id" {
  type = string
}

variable "region" {
  type    = string
  default = "europe-central2"
}

variable "registry_id" {
  type    = string
  default = "eu.gcr.io"
}

variable "port" {
  type    = number
  default = 9001
}

variable "bot_name" {
  type    = string
  default = "bot"
}

variable "downloader_name" {
  type    = string
  default = "downloader"
}