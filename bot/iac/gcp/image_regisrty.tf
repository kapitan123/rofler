resource "google_artifact_registry_repository" "eu_gcr_io" {
  format        = "DOCKER"
  location      = "europe"
  project       = var.project_id
  repository_id = var.registry_id
}