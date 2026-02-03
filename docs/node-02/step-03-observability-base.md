# Node 02 / Step 03 - 对象存储与可观测基线（OSS + Prometheus/Grafana）

## 背景/目标
- 为模型结果与中间产物提供对象存储
- 为平台提供基础监控与可观测能力

## 范围与非目标
- 范围：OSS 桶资源、监控 Helm 基线
- 非目标：日志索引/可视化细化、告警策略与告警路由

## 设计与接口
- OSS：`infra/terraform/oss.tf`
- 监控：`deploy/helm/maas-observability`
- 依赖：kube-prometheus-stack（Prometheus/Grafana）

## 实现步骤
1. 新增 `oss.tf` 与 OSS 变量/输出
2. 新增 `maas-observability` Helm chart
3. 预留 Loki/Promtail 作为可选日志组件

## 测试与验收
- `terraform plan` 包含 OSS 桶资源
- `helm dependency update` 可拉取监控依赖

## 回滚策略
- `git revert <commit>` 回滚本步骤提交

## 变更记录
- commit: `feat: add observability base`
