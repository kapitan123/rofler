resource "google_cloud_run_service" "telegrofler" {
  name     = var.name
  location = var.region
  project  = var.project_id

  lifecycle {
    ignore_changes = [
      template[0].spec[5],
    ]
  }

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
            memory = "1 Gi"
          }
        }

        env {
          name  = "TELEGRAM_BOT_TOKEN"
          value = var.bot_token
        }

        env {
          name  = "DELETION_QUEUE_NAME"
          value = local.deletionQueueName
        }

        env {
          name  = "PROJECT_ID"
          value = var.project_id
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.cloud_run_googleapis_com]
}

# Allow unauthenticated users to invoke the service
resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = google_cloud_run_service.telegrofler.name
  location = google_cloud_run_service.telegrofler.location
  role     = "roles/run.invoker"
  member   = "allUsers"
} 