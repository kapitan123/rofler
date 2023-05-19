resource "google_pubsub_topic" "convertor_video_converted_topic" {
  name = "convertor_video_converted"
}

resource "google_pubsub_topic" "bot_video_link_published_topic" {
  name = "bot_video_link_published"
}