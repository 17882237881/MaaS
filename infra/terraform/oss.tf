resource "alicloud_oss_bucket" "maas" {
  bucket = var.oss_bucket_name
  acl    = var.oss_bucket_acl
}
