resource "google_pubsub_topic" "uploadable_video_saved" {
  name = "uploadable_video_saved"
}

resource "google_pubsub_topic" "video_url_published" {
  name = "video_url_published"
}

resource "google_pubsub_topic" "dead_letter_received" {
  name = "dead_letter_received"
}

resource "google_pubsub_subscription" "downloader_to_video_url_published" {
  name  = "downloader_to_video_url_published"
  topic = google_pubsub_topic.video_url_published.name

  push_config {
    push_endpoint = "${local.bot_url}/api/pubsub/subscriptions/video-url-published"

    attributes = {
      x-goog-version = "v1"
    }
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead_letter_received.id
    max_delivery_attempts = 5
  }

  ack_deadline_seconds = 20
}

resource "google_pubsub_subscription" "bot_to_saved_videos" {
  name  = "bot_to_saved_videos"
  topic = google_pubsub_topic.uploadable_video_saved.name

  push_config {
    push_endpoint = "${local.downloader_url}/api/pubsub/subscriptions/video-saved"

    attributes = {
      x-goog-version = "v1"
    }
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead_letter_received.id
    max_delivery_attempts = 5
  }

  ack_deadline_seconds = 20

}
