variable "project_id" {
  type    = string
  default = "neat-cycling-346311" // secret
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

variable "state_bucket" {
  type    = string
  default = "304373b586ff10b7-bucket-tfstate" // secret
}
