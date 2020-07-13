resource "google_kms_key_ring" "sops_key_ring" {
  // project  = "kaito2"
  name     = "sops-key-ring"
  location = "asia-northeast1"
}

resource "google_kms_crypto_key" "sops_crypto_key" {
  name     = "sops-crypto-key"
  key_ring = google_kms_key_ring.sops_key_ring.self_link
}
