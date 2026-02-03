# Node 02 / Step 01 - Terraform 基础设施基线（ACK + GPU 节点池）

## 背景/目标
- 使用 IaC 方式定义基础设施
- 建立 VPC/VSwitch/ACK 集群与 GPU 节点池的最小可用模板

## 范围与非目标
- 范围：Terraform 目录结构、基础资源定义
- 非目标：生产级多地域/多可用区拓扑与网络安全细节

## 设计与接口
- 目录：`infra/terraform`
- 资源：VPC、VSwitch、ACK Managed Cluster、GPU Node Pool
- 变量：region、CIDR、instance_types、key_name 等

## 实现步骤
1. 新增 `infra/terraform/versions.tf` 定义 provider 与版本
2. 新增 `infra/terraform/variables.tf` 定义可配置项
3. 新增 `infra/terraform/main.tf` 定义 VPC/VSwitch/ACK/NodePool
4. 新增 `infra/terraform/outputs.tf` 输出关键资源 ID
5. 新增 `infra/terraform/README.md` 使用说明

## 测试与验收
- `terraform fmt` 可通过（如需要）
- `terraform init` 可完成 provider 下载
- `terraform plan` 在提供必填变量后可生成计划

## 回滚策略
- `git revert <commit>` 回滚本步骤提交

## 变更记录
- commit: `feat: add terraform base for k8s`
