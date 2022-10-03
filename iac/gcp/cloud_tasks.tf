resource "google_cloud_tasks_queue" "tg_message_deletion" {
  name     = local.deletionQueueName
  location = var.region
  project  = var.project_id

  rate_limits {
    max_concurrent_dispatches = 100
  }

  retry_config {
    max_attempts       = 5
    max_retry_duration = "4s"
    max_backoff        = "3s"
    min_backoff        = "2s"
    max_doublings      = 1
  }

  depends_on = [google_project_service.cloud_tasks_googleapis_com]
}

locals {
  deletionQueueName = "tg-message-deletion"
}