# Node 02 / Step 02 - Helm 中间件基线（PostgreSQL / Redis / Kafka）

## 背景/目标
- 为平台提供数据库、缓存、消息队列的统一部署方式
- 通过 Helm 形成可复用的部署基线

## 范围与非目标
- 范围：Helm umbrella chart + 依赖声明 + 基础 values
- 非目标：生产级高可用参数与备份策略

## 设计与接口
- Chart：`deploy/helm/maas-middleware`
- 依赖：Bitnami PostgreSQL / Redis / Kafka
- values：统一配置入口（密码、持久化、资源）

## 实现步骤
1. 新建 `deploy/helm/maas-middleware/Chart.yaml` 声明依赖
2. 新建 `values.yaml` 提供基础配置
3. 新建 `templates/NOTES.txt` 提示使用方式

## 测试与验收
- `helm dependency update` 可拉取依赖
- `helm template` 可渲染无语法错误

## 回滚策略
- `git revert <commit>` 回滚本步骤提交

## 变更记录
- commit: `feat: add helm charts for middleware`
