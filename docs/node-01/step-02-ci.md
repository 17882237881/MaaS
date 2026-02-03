# Node 01 / Step 02 - CI 基线（GitHub Actions）

## 背景/目标
- 为每次提交与 PR 提供自动化校验
- 保证 Go 代码可编译、可测试、可静态检查

## 范围与非目标
- 范围：基础 CI 流水线（checkout/setup-go/tidy/test/vet）
- 非目标：发布、部署、性能测试

## 设计与接口
- 触发条件：push 到 main / PR 到 main
- 任务：`go mod tidy`、`go test ./...`、`go vet ./...`

## 实现步骤
1. 新增 `.github/workflows/ci.yml`
2. 配置 Go 版本与缓存
3. 执行 `go mod tidy` 并校验仓库无差异
4. 执行单元测试与 vet

## 测试与验收
- 推送后 GitHub Actions 通过
- CI 在无代码变更情况下无 `git diff`

## 回滚策略
- `git revert <commit>` 回滚本步骤提交

## 变更记录
- commit: `chore: add github actions pipeline`
