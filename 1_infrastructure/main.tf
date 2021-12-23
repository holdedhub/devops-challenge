provider "google" {
  project = var.project_id
  region  = var.region
}

#provider "kubernetes" {
#  host                   = "https://${module.gke.endpoint}"
#  token                  = data.google_client_config.default.access_token
#  cluster_ca_certificate = base64decode(module.gke.ca_certificate)
#}

data "google_client_config" "default" {}

module "network" {
  source  = "terraform-google-modules/network/google"
  version = "4.0.1"
  project_id   = var.project_id
  network_name = "${var.project_id}-${var.network}"
  subnets = [
    {
      subnet_name   = "${var.project_id}-${var.subnetwork}"
      subnet_ip     = "10.10.0.0/16"
      subnet_region = var.region
    },
  ]
  secondary_ranges = {
    "${var.project_id}-${var.subnetwork}" = [
      {
        range_name    = "${var.project_id}-${var.ip_range_pods}"
        ip_cidr_range = "10.104.0.0/14"
      },
      {
        range_name    = "${var.project_id}-${var.ip_range_services}"
        ip_cidr_range = "10.108.0.0/20"
      },
    ]
  }
}

module "gke" {
  source  = "terraform-google-modules/kubernetes-engine/google"
  version = "17.3.0"
  project_id                  = var.project_id
  name                        = "${var.project_id}-gke-cluster"
  regional                    = true
  region                      = var.region
  network                     = module.network.network_name
  subnetwork                  = module.network.subnets_names[0]
  ip_range_pods               = module.network.subnets_secondary_ranges[0].*.range_name[0]
  ip_range_services           = module.network.subnets_secondary_ranges[0].*.range_name[1]
  create_service_account      = false

  node_pools = [
    {
      name                      = "node-pool"
      machine_type              = "e2-medium"
      min_count                 = 1
      max_count                 = 2
      disk_size_gb              = 30
    },
  ]
}

