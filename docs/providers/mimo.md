# 小米 MiMo 配置指南

## 获取 API Key

1. 访问 [小米 MiMo 平台](https://platform.xiaomimimo.com/)
2. 注册并登录账号
3. 在 [套餐管理](https://platform.xiaomimimo.com/console/plan-manage) 页面选择适合的套餐
4. 在 [API Keys](https://platform.xiaomimimo.com/console/api-keys) 页面创建新的 API Key 并复制

::: tip
账号余额可在 [余额页面](https://platform.xiaomimimo.com/console/balance) 查看。
:::

## 在 CCX 中添加渠道

| 字段 | 值 |
|------|-----|
| 名称 | `MiMo`（自定义） |
| 服务类型 | `openai` |
| Base URL | `https://api.mimo.xiaomi.com/v1` |
| API Keys | 你的 MiMo API Key |

### 配置步骤

1. 进入 CCX 管理界面，选择 **Chat** 入口
2. 点击「添加渠道」
3. 填写以下信息：
   - **名称**：`MiMo`
   - **服务类型**：选择 `OpenAI Chat`
   - **Base URL**：`https://api.mimo.xiaomi.com/v1`
   - **API Keys**：粘贴你的 API Key
4. 点击保存

### 模型白名单（可选）

```
mimo-v2.5-pro
mimo-v2.5
mimo-v2-flash
```

### 模型映射（可选）

```json
{
  "mimo-pro": "mimo-v2.5-pro",
  "mimo": "mimo-v2.5",
  "mimo-flash": "mimo-v2-flash"
}
```

## 可用模型

| 模型 | 说明 |
|------|------|
| `mimo-v2.5-pro` | 最新旗舰，1.02T 总参 / 42B 激活 |
| `mimo-v2.5` | 310B 总参 / 15B 激活，原生多模态 |
| `mimo-v2-flash` | 309B 总参 / 15B 激活，高速推理 |

## 注意事项

- MiMo 通过兼容 OpenAI 协议的平台访问
- MiMo 是推理模型，支持 `reasoning_content` 字段返回思考过程

### 视觉支持

MiMo 各模型的视觉支持情况：

| 模型 | 视觉支持 |
|------|----------|
| `MiMo-V2.5-Pro` | 不支持 |
| `MiMo-V2.5` | 支持（原生多模态） |
| `MiMo-V2-Flash` | 不支持 |

::: warning
`MiMo-V2.5-Pro` 不支持图片输入。如果需要处理包含图片的请求，必须配置**视觉回退模型**。
:::

**配置方式：** 编辑渠道，在「视觉回退模型」字段填入 `MiMo-V2.5`。当请求包含图片且目标模型（如 `MiMo-V2.5-Pro`）不支持视觉时，CCX 会自动使用 `MiMo-V2.5` 替代模型处理该请求。

![MiMo 渠道配置：模型重定向将 haiku/opus/sonnet 统一映射到 mimo-v2.5-pro；下方「视觉回退模型」字段填入 mimo-v2.5，当请求包含图片且目标模型不支持视觉时自动切换到该模型](/screenshots/channel-mimo-config.png)

<figcaption>MiMo 渠道配置 — 模型重定向（haiku/opus/sonnet → mimo-v2.5-pro）与视觉回退模型（mimo-v2.5）</figcaption>

如果留空视觉回退模型，包含图片的请求将跳过该渠道，failover 到下一个支持视觉的渠道。

### 回传思考内容

MiMo 作为推理模型，思考过程中产生的 `reasoning_content` 需要回传给 API，否则会返回 HTTP 400 错误。

**必须启用：** 编辑渠道时打开「回传思考内容」开关（`PassbackReasoningContent`）。

启用后 CCX 会自动处理：
- **请求方向：** 为缺少 `thinking` 块的 assistant 消息注入占位符，满足 MiMo 的回传要求
- **响应方向：** 将上游返回的 `reasoning_content` 转换为 Claude 原生的 `thinking` 内容块，下游客户端可正常解析

::: warning
不开启此开关会导致包含历史对话的推理请求失败（HTTP 400）。
:::

### 推荐的模型映射

将 Claude 模型名映射到 MiMo 模型：

| 请求模型 | 重定向到 | 说明 |
|----------|----------|------|
| `haiku` | `mimo-v2.5-pro` | 旗舰推理 |
| `opus` | `mimo-v2.5-pro` | 旗舰推理 |
| `sonnet` | `mimo-v2.5-pro` | 旗舰推理 |

::: tip
所有 Claude 模型统一映射到 `MiMo-V2.5-Pro` 以获得最佳推理能力。配合视觉回退模型 `MiMo-V2.5` 使用，图片请求会自动切换到支持视觉的模型。
:::
