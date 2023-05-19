// AK TODO remame resource to bot
// Callback processing service
resource "google_cloud_run_service" "bot" {
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

        # AK TODO REMOVE IF ID passing works as expected
        env {
          name  = "MESSAGE_DELETION_QUEUE_NAME"
          value = var.message_deletion_queue_name
        }

        env {
          name  = "MESSAGE_DELETION_QUEUE_ID"
          value = google_cloud_tasks_queue.tg_message_deletion.id
        }

        env {
          name  = "VIDEO_PUBLISHED_TOPIC"
          value = google_pubsub_topic.bot_video_link_published_topic.id
        }

        env {
          name  = "PROJECT_ID"
          value = var.project_id
        }

        env {
          name  = "VIDEO_FILES_BUCKET_URL"
          value = google_storage_bucket.converted_videos.url
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
          name  = "VIDEO_FILES_BUCKET_URL"
          value = google_storage_bucket.converted_videos.url
        }

        env {
          name  = "VIDEO_CONVERTED_TOPIC"
          value = google_pubsub_topic.convertor_video_converted_topic.id
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
  service  = google_cloud_run_service.bot.name
  location = google_cloud_run_service.bot.location
  role     = "roles/run.invoker"
  member   = "allUsers"
} 