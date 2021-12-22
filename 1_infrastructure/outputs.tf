output "project_id" {
  value = var.project_id
}

output "region" {
  value = module.gke.region
}

output "cluster_name" {
  description = "Cluster name"
  value       = module.gke.name
}

output "kubernetes_endpoint" {
  sensitive = true
  value = module.gke.endpoint
}

output "master_kubernetes_version" {
  description = "The master Kubernetes version"
  value       = module.gke.master_version
}
