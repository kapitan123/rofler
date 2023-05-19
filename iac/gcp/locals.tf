locals {
  telegrofler_url = google_cloud_run_service.telegrofler.status[0].url
  bot_image_url       = "${var.registry_id}/${var.project_id}/${var.bot_name}:latest"
  convertor_image_url       = "${var.registry_id}/${var.project_id}/${var.convertor_name}:latest"
}