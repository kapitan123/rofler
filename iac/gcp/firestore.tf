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