resource "google_compute_network" "main" {
  name                    = "sharedtelemetryapp-network"
  auto_create_subnetworks = false
  routing_mode            = "REGIONAL"
}

resource "google_compute_subnetwork" "private_01" {
  name                     = "private-01"
  ip_cidr_range            = "10.2.0.0/24"
  region                   = "europe-west3"
  network                  = google_compute_network.main.id
  private_ip_google_access = true
}

resource "google_compute_global_address" "private_ip_address" {
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  address       = "10.3.0.0"
  network       = google_compute_network.main.id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}
