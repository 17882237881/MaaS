# Node 01 / Step 01 - 项目初始化与工程基线

## 背景/目标
- 初始化仓库的基础结构与 Go Module
- 为后续 CloudWeGo（Hertz/Kitex/Protobuf）接入预留目录与规范
- 提供最小化的工程操作入口（Makefile）

## 范围与非目标
- 范围：目录结构、Go Module、基础脚手架文件
- 非目标：任何具体业务服务、IDL 或框架代码

## 设计与接口
- 统一仓库布局（cmd/internal/api/deploy/infra/docs/scripts）
- 采用 Go Module `github.com/17882237881/MaaS`
- 说明：CloudWeGo 具体服务骨架将在后续节点生成

## 实现步骤
1. 创建目录结构：`cmd/`、`internal/`、`api/proto/`、`deploy/`、`infra/`、`docs/`、`scripts/`
2. 初始化 Go Module：`go mod init github.com/17882237881/MaaS`
3. 增加 `.gitignore`（Go/IDE/OS/Env）
4. 增加 `Makefile`（tidy/test/vet）
5. 增加 `README.md`（项目简介与结构说明）
6. 使用 `.gitkeep` 保留空目录

## 测试与验收
- `make test` 应能通过（当前无业务代码）
- `make vet` 无报错
- 仓库结构与文档路径符合约定

## 回滚策略
- `git revert <commit>` 回滚本步骤提交

## 变更记录
- commit: `chore: bootstrap repo structure with cloudwego`
