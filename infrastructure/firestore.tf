resource "google_firestore_database" "database" {
  name        = "(default)"
  location_id = "eur3"
  type        = "FIRESTORE_NATIVE"
}
