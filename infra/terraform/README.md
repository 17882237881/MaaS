# Terraform - ACK 基础设施

## 前置条件
- Terraform >= 1.5
- 阿里云账号与访问密钥
- 已创建的 ECS 密钥对（用于 `key_name`）

## 快速使用
```bash
terraform init
terraform apply \
  -var "worker_instance_types=[\"ecs.<your-gpu-instance-type>\"]" \
  -var "key_name=your-keypair-name"
```

## 说明
- `worker_instance_types` 请替换为实际 GPU 规格
- `k8s_version` 留空表示使用默认版本
