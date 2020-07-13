provider "google" {
  project = "kaito2"
  region  = "asia-northeast1-b"
}

// cloud kms を有効化
resource "google_project_service" "kms" {
  service = "cloudkms.googleapis.com"
}
