resource "google_cloud_run_service" "telegrofler" {
  name     = var.name
  location = var.region
  project  = var.project_id

  template {
    spec {
      containers {
        image = local.image_url

        ports {
          container_port = var.port
        }

        resources {
          limits = {
            cpu    = 2.0
            memory = "1GiB"
          }
        }

        env {
          name = "TELEGRAM_BOT_TOKEN"
        }

        env {
          name = "FIRESTORE_PROJECT_ID"
        }

        env {
          name = "GOOGLE_CLOUD_PROJECT"
        }

      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cr_googleapis_com]
}

# Allow unauthenticated users to invoke the service
resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = google_cloud_run_service.telegrofler.name
  location = google_cloud_run_service.telegrofler.location
  role     = "roles/run.invoker"
  member   = "allUsers"
} 