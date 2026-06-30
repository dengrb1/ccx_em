# 仓库协作指南

## 重要约定
- **始终使用简体中文回复**。
- 遵循 SOLID / KISS / DRY / YAGNI；优先修复根因，避免无关重构。
- 修改文档时，优先以以下文件为事实源：`VERSION`、`Makefile`、`backend-go/Makefile`、`frontend/package.json`、`backend-go/main.go`。
- 不要手动编辑生成产物：`dist/`、`frontend/dist/`、`backend-go/frontend/dist/`、`desktop/frontend/dist/`、`desktop/bin/`、`desktop/dist/`。

## 项目概览
- CCX 是一个多上游 AI API 代理与协议转换网关，当前正式支持六类渠道：`messages`、`chat`、`responses`、`gemini`、`images`、`vectors`。
- 对外代理入口覆盖 Claude Messages、OpenAI Chat Completions、OpenAI Responses、Gemini 原生协议、OpenAI Images、OpenAI Embeddings。
- 根目录 `VERSION` 是唯一发布版本源；后端构建时通过 `backend-go/Makefile` 的 `-ldflags` 注入运行时版本信息。

## 项目结构与模块
- `backend-go/`：主 Go 服务（Gin），负责路由、认证、调度、协议转换、指标、日志、会话与热重载配置；前端构建产物嵌入到 `backend-go/frontend/dist/`。
- `frontend/`：Vue 3 + Vite + Vuetify 管理界面。
- `desktop/`：Wails 3 桌面端壳层；`desktop/frontend/` 为桌面端前端。
- `shared/`：前后端/桌面端共享资源与约定。
- `dist/`：发布构建产物，禁止手动编辑。
- `.config/`：运行时配置与持久化目录，如 `config.json`、`metrics.db`、`conversation_state.json`、`backups/`。
- `refs/`：外部参考项目存档，仅供对照，默认只读。
- 文档入口：`README.md`、`README.zh-CN.md`、`backend-go/README.md`、`docs/guide/architecture.md`、`docs/guide/development.md`、`docs/guide/environment.md`、`docs/guide/release.md`、`docs/guide/desktop/`。

## 构建 / 测试 / 开发命令
- 安装依赖：根目录 `make install`（前端、后端、桌面端依赖与开发工具）。
- 全栈开发：根目录 `make dev`（前端 `bun run dev` + 后端 `air` 热重载）。
- 根目录常用命令：`make run`、`make build`、`make clean`、`make frontend-dev`、`make frontend-build`、`make desktop-dev`、`make desktop-build`。
- 仅后端：`cd backend-go && make dev|run|build|build-local|test|test-cover|fmt|lint|deps`。
- 前端：`cd frontend && bun install` 后执行 `bun run dev|build|preview|type-check|lint|test`。
- 桌面端：`cd desktop && wails3 task dev|package`。
- Docker：`docker-compose up -d` 用于接近生产环境的验证。

## 代码风格
- Go：保持包职责单一、接口清晰；修改后运行 `cd backend-go && go fmt ./...`。
- 前端：遵循现有 Vue 3 / TypeScript / Vuetify / Prettier / ESLint 风格，保持 strict 类型约束。
- 桌面端：遵循现有 Wails 3、Go 与前端风格，避免引入与 Web 端重复但未抽象的实现。
- 配置/密钥：`.env`、`.json` 只提交示例文件或文档化内容，禁止提交真实密钥。

## 路由与能力边界
- 实际代理与管理路由以 `backend-go/main.go` 为准。
- 常见代理入口包括：
  - `/v1/messages`
  - `/v1/messages/count_tokens`
  - `/v1/models`
  - `/v1/models/:model`
  - `/v1/chat/completions`
  - `/v1/responses`
  - `/v1/responses/compact`
  - `/v1/images/generations`
  - `/v1/images/edits`
  - `/v1/images/variations`
  - `/v1/embeddings`
  - `/v1beta/models/*`
- 管理入口按类型分组，统一位于：
  - `/api/messages/channels/*`
  - `/api/responses/channels/*`
  - `/api/chat/channels/*`
  - `/api/gemini/channels/*`
  - `/api/images/channels/*`
  - `/api/vectors/channels/*`
- capability-test / snapshot 当前只适用于 `messages`、`chat`、`responses`、`gemini`；不要假设 `images` 或 `vectors` 具备对应路由。
- `vectors` 通道对应 OpenAI Embeddings 代理入口 `/v1/embeddings`，管理维度使用 `/api/vectors/channels/*`。
- `responses` 支持 `GET /v1/responses` 的 WebSocket fallback 语义，处理兼容客户端回退。

## 测试规范
- 新增或修改后端逻辑尽量补 `_test.go`，优先表驱动 + `httptest`。
- 前端当前已有 `vitest` 能力；增加复杂逻辑时优先补轻量单测，并至少通过 `bun run test`、`bun run type-check`。
- 涉及桌面端联动时，至少验证核心构建或启动流程，不要只改文档或前端静态代码就结束。
- 文档或接口说明调整后，至少验证：`make build`、`cd backend-go && make test`、`cd frontend && bun run build`。

## 安全与配置提示
- 部署前必须设置强 `PROXY_ACCESS_KEY`；如需分离管理权限，再配置 `ADMIN_ACCESS_KEY`。
- 生产环境下不要使用默认访问密钥；后端启动时会校验访问密钥配置。
- `.config/config.json` 支持热重载；修改 `backend-go/.env` 后通常需要重启服务。
- 代理端点统一鉴权（`x-api-key` / `Authorization: Bearer`）；生产环境不建议依赖 query `key`。
- 记录或展示日志时注意脱敏，尤其是 API Key、Authorization 头和 multipart 请求内容。
