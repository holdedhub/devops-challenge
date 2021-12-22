terraform {
  backend "gcs" {
    bucket      = "devops-challenge-tfstate"
    prefix      = "terraform/state"
  }
}
