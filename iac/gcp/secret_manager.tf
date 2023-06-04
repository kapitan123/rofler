resource "google_secret_manager_secret" "downloader_cookies" {
  secret_id = "downloader_cookies"

  labels = {
    owner = "downloader"
  }

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret" "bot_token" {
  secret_id = "telegram_bot_token"

  labels = {
    owner = "bot"
  }

  replication {
    automatic = true
  }
}