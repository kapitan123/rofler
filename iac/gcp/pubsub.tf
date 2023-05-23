resource "google_pubsub_topic" "convertor_video_converted_topic" {
  name = "convertor_video_converted"
}

resource "google_pubsub_topic" "bot_video_link_published_topic" {
  name = "bot_video_link_published"
}

resource "google_pubsub_topic" "dead_letter_topic" {
  name = "global_dead_letter_received"
}

resource "google_pubsub_subscription" "convertor_to_published_videos" {
  name  = "convertor-to-published-videos"
  topic = google_pubsub_topic.bot_video_link_published_topic.name

  push_config {
    push_endpoint = "${local.bot_url}/pubsub/subscriptions/video-published"

    attributes = {
      x-goog-version = "v1"
    }
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead_letter_topic.id
    max_delivery_attempts = 5
  }

  ack_deadline_seconds = 20
}

resource "google_pubsub_subscription" "bot_to_converted_videos" {
  name  = "bot-to-converted-videos"
  topic = google_pubsub_topic.convertor_video_converted_topic.name

  push_config {
    push_endpoint = "${local.convertor_url}/pubsub/subscriptions/video-converted"

    attributes = {
      x-goog-version = "v1"
    }
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead_letter_topic.id
    max_delivery_attempts = 5
  }

  ack_deadline_seconds = 20

}
