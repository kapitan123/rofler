data "google_secret_manager_secret" "downloader_cookies" {
  secret   = "downloader_cookies"
}

data "google_secret_manager_secret" "bot_token" {
  secret   = "telegram_bot_token"
}