resource "random_id" "bucket_prefix" {
  byte_length = 8
}

resource "google_storage_bucket" "tfstate" {
  name          = "${random_id.bucket_prefix.hex}-bucket-tfstate"
  force_destroy = false
  location      = var.region
  storage_class = "STANDARD"

  versioning {
    enabled = true
  }
}

resource "google_storage_bucket" "converted_videos" {
  name          = "${random_id.bucket_prefix.hex}-converted-videos"
  force_destroy = true
  location      = var.region
  storage_class = "STANDARD"

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }
}