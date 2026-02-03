output "vpc_id" {
  value = alicloud_vpc.main.id
}

output "vswitch_id" {
  value = alicloud_vswitch.main.id
}

output "ack_cluster_id" {
  value = alicloud_cs_managed_kubernetes.main.id
}

output "gpu_node_pool_id" {
  value = alicloud_cs_kubernetes_node_pool.gpu.id
}
