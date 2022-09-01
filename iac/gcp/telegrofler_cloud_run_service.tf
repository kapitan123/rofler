resource "google_cloud_run_service" "run_service" {
  name = "telegrofler"
  location = "europe-central2"

  // AK TODO check min config

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.run_api]
}

// Add config for firestore


