resource "google_firestore_index" "posts_by_date_and_chat" {
  project = var.project_id

  collection = "posts"

  fields {
    field_path = "chosen_on"
    order      = "ASCENDING"
  }

  fields {
    field_path = "chat_id"
    order      = "ASCENDING"
  }
}

resource "google_firestore_field" "expire_at_field" {
  project    = var.project_id
  collection = "media"
  database   = "(default)"
  field      = "expire_at"

  index_config {
    indexes {
      order = "ASCENDING"
    }
  }

  ttl_config {}
}