resource "google_cloudfunctions_function" "convertor" {
  name        = "convertor-function"
  description = "Function to convert videos"
  runtime     = "go120"
  entry_point = "Convert"

  environment_variables = {
    BOT_CALLBACK = "value1"
  }

  depends_on = [google_project_service.cloud_tasks_googleapis_com]

  trigger_http = true
}