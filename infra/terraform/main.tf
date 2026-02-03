provider "alicloud" {
  region = var.region
}

data "alicloud_zones" "available" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "main" {
  vpc_name   = var.cluster_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "main" {
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = var.vswitch_cidr
  zone_id      = data.alicloud_zones.available.zones[0].id
  vswitch_name = "${var.cluster_name}-vsw"
}

resource "alicloud_cs_managed_kubernetes" "main" {
  name_prefix          = var.cluster_name
  version              = var.k8s_version
  worker_vswitch_ids   = [alicloud_vswitch.main.id]
  pod_cidr             = var.pod_cidr
  service_cidr         = var.service_cidr
  slb_internet_enabled = true
  new_nat_gateway      = true
}

resource "alicloud_cs_kubernetes_node_pool" "gpu" {
  cluster_id     = alicloud_cs_managed_kubernetes.main.id
  node_pool_name = "${var.cluster_name}-gpu"
  vswitch_ids    = [alicloud_vswitch.main.id]

  instance_types = var.worker_instance_types
  key_name       = var.key_name
  desired_size   = var.node_pool_desired_size

  scaling_config {
    min_size = var.node_pool_min_size
    max_size = var.node_pool_max_size
  }

  system_disk_category  = var.system_disk_category
  system_disk_size      = var.system_disk_size
  install_cloud_monitor = true

  data_disks {
    size     = var.data_disk_size
    category = var.data_disk_category
  }
}
