# Terraform - ACK 基础设施

## 前置条件
- Terraform >= 1.5
- 阿里云账号与访问密钥
- 已创建的 ECS 密钥对（用于 `key_name`）
- 已规划 OSS 桶名（全局唯一）

## 快速使用
```bash
terraform init
terraform apply \
  -var "worker_instance_types=[\"ecs.<your-gpu-instance-type>\"]" \
  -var "key_name=your-keypair-name" \
  -var "oss_bucket_name=maas-prod-oss"
```

## 说明
- `worker_instance_types` 请替换为实际 GPU 规格
- `k8s_version` 留空表示使用默认版本
- `oss_bucket_name` 需要全局唯一
