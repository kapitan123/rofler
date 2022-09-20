locals {
  telegrofler_url = google_cloud_run_service.telegrofler.status[0].url
  image_url       = "${var.registry_id}/${var.project_id}/${var.name}:latest"
}