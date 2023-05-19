// AK TODO remame resource to bot
// Callback processing service
resource "google_cloud_run_service" "telegrofler" {
  name     = var.bot_name
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
        image = local.bot_image_url

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
          name  = "MESSAGE_DELETION_QUEUE_NAME"
          value = var.message_deletion_queue_name
        }

        env {
          name  = "VIDEO_CONVERTION_TOPIC_NAME"
          value = var.video_convertion_topic_name
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

// Convertor service
resource "google_cloud_run_service" "convertor" {
  name     = var.convertor_name
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
        image = local.bot_image_url

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
          name  = "VIDEO_FILES_BUCKET"
          value = local.video_files_bucket
        }

        env {
          name  = "VIDEO_CONVERTION_QUEUE_NAME"
          value = local.video_convertion_queue_name
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