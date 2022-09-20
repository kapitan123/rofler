
// AK TODO If admin is never used I see no point in submitting it
resource "google_service_account" "telegrofler_adm" {
  account_id  = "telegrofler-adm"
  description = "Telegrofler admin"
  project     = var.project_id
}