variable "region" {
  type    = string
  default = "cn-hangzhou"
}

variable "cluster_name" {
  type    = string
  default = "maas-ack"
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "vswitch_cidr" {
  type    = string
  default = "10.0.1.0/24"
}

variable "pod_cidr" {
  type    = string
  default = "172.20.0.0/16"
}

variable "service_cidr" {
  type    = string
  default = "172.21.0.0/20"
}

variable "k8s_version" {
  type     = string
  nullable = true
  default  = null
}

variable "worker_instance_types" {
  type = list(string)
}

variable "key_name" {
  type = string
}

variable "node_pool_desired_size" {
  type    = number
  default = 1
}

variable "node_pool_min_size" {
  type    = number
  default = 1
}

variable "node_pool_max_size" {
  type    = number
  default = 3
}

variable "system_disk_category" {
  type    = string
  default = "cloud_efficiency"
}

variable "system_disk_size" {
  type    = number
  default = 100
}

variable "data_disk_category" {
  type    = string
  default = "cloud_essd"
}

variable "data_disk_size" {
  type    = number
  default = 120
}

variable "oss_bucket_name" {
  type = string
}

variable "oss_bucket_acl" {
  type    = string
  default = "private"
}
