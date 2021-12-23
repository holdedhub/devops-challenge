resource "google_artifact_registry_repository" "my-repo" {
  provider = google-beta

  project = var.project_id
  location = var.region
  repository_id = "${var.project_id}-repo"
  description = "example docker repository"
  format = "DOCKER"
}
