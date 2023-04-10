
resource "google_cloud_scheduler_job" "choose_random_p" {
  name        = "choose_random_p"
  description = "test http job"
  schedule    = "0 8 * * *" // 9:00 Berlin time
  time_zone   = "Europe/London"
  region      = var.region

  http_target {
    http_method = "POST"
    uri         = "${local.telegrofler_url}/chat/-1001201899231/pidoroftheday"
  }
}