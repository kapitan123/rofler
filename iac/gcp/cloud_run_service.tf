resource "google_cloud_run_service" "bot" {
  name     = var.bot_name
  location = var.region
  project  = var.project_id

  autogenerate_revision_name = true

  template {
    spec {
      containers {
        image = local.bot_image_url

        ports {
          container_port = var.port
        }

        resources {
          limits = {
            cpu    = "2000m"
            memory = "1Gi"
          }
        }

        env {
          name  = "TELEGRAM_BOT_TOKEN"
          value = var.bot_token
        }

        env {
          name  = "VIDEO_URL_PUBLISHED_TOPIC"
          value = google_pubsub_topic.video_url_published.name
        }

        env {
          name  = "VIDEO_FILES_BUCKET"
          value = google_storage_bucket.videos.name
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

resource "google_cloud_run_service" "downloader" {
  name     = var.downloader_name
  location = var.region
  project  = var.project_id

  autogenerate_revision_name = true

  template {
    spec {
      containers {
        image = local.downloader_image_url

        ports {
          container_port = var.port
        }

        resources {
          limits = {
            cpu    = "2000m"
            memory = "1Gi"
          }
        }

        env {
          name  = "VIDEO_FILES_BUCKET"
          value = google_storage_bucket.videos.name
        }

        env {
          name  = "VIDEO_SAVED_TOPIC"
          value = google_pubsub_topic.uploadable_video_saved.name
        }

        env {
          name  = "PROJECT_ID"
          value = var.project_id
        }

        env {
          name  = "DOWNLOADER_COOKIES"
          value = var.downloader_cookies
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
  service  = google_cloud_run_service.bot.name
  location = google_cloud_run_service.bot.location
  role     = "roles/run.invoker"
  member   = "allUsers"
} 