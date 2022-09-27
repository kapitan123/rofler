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

variable "name" {
  type = string
}

variable "port" {
  type    = number
  default = 9001
}
