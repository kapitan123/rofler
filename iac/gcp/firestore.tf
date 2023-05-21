resource "google_firestore_database" "database" {
  project                     = var.project_id
  name                        = "(default)"
  type                        = "FIRESTORE_NATIVE"
  location_id                 = var.region
  concurrency_mode            = "OPTIMISTIC"
  app_engine_integration_mode = "DISABLED"
}

resource "google_firestore_index" "pidors_by_date_and_chat" {
  project = var.project_id

  collection = "pidors"

  fields {
    field_path = "chosen_on"
    order      = "ASCENDING"
  }

  fields {
    field_path = "chat_id"
    order      = "ASCENDING"
  }
}

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