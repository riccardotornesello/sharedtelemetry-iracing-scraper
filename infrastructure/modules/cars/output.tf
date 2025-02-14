output "db" {
  value = google_sql_database.database
}

output "db_user" {
  value = google_sql_user.cars_downloader
}
