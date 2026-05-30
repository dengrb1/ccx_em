# Homebrew 安装与更新

通过 [Homebrew](https://brew.sh) 安装和更新 CCX Desktop macOS 版本。

## 安装

```bash
brew tap BenedictKing/ccx
brew install --cask ccx-desktop
```

## 更新

```bash
brew update
brew upgrade --cask ccx-desktop
```

## 仓库结构

自建 Tap 仓库：`BenedictKing/homebrew-ccx`

```
homebrew-ccx/
├── Casks/
│   └── ccx-desktop.rb
└── .github/
    └── workflows/
        └── update-cask.yml
```

- `Casks/ccx-desktop.rb` — Cask 定义，支持 arm64/amd64 双架构自动选择；由 Tap workflow 自动生成。
- `update-cask.yml` — 响应主仓库 `repository_dispatch` 事件，自动更新 Cask。

Cask 模板文件维护在主仓库 `packaging/homebrew/` 目录下，供参考。首次初始化 Tap 时，只需要先复制 `packaging/homebrew/update-cask.yml` 到 Tap 仓库的 `.github/workflows/update-cask.yml`，然后手动触发主仓库 `Update Homebrew Tap` workflow 生成真实的 `Casks/ccx-desktop.rb`。

## 自动更新流程

```
推送 vX.Y.Z tag
       │
       ▼
  release.yml
  build-macos → 构建 Draft Release 和 DMG (arm64 + amd64)
       │
       ▼
  维护者确认并发布 GitHub Release
       │
       ▼
  update-homebrew.yml (release: published)
  下载 DMG SHA256 → repository_dispatch
       │
       ▼
  homebrew-ccx/update-cask.yml
  更新 Casks/ccx-desktop.rb → commit + push
```

主仓库 `release.yml` 只负责构建 Draft Release。维护者确认并发布 GitHub Release 后，`.github/workflows/update-homebrew.yml` 会从已发布 Release 下载 DMG SHA256 校验文件，然后向 Tap 仓库发送 `repository_dispatch` 事件，Tap 仓库自动更新 Cask 定义并推送。

也可以手动触发 `Update Homebrew Tap` workflow，输入已发布的 tag（例如 `v2.8.17`）重新同步 Tap。

## 维护者配置

在主仓库 `BenedictKing/ccx` 的 GitHub Secrets 中添加：

| Secret | 说明 |
|---|---|
| `HOMEBREW_TAP_TOKEN` | 用于触发 Tap 仓库 `repository_dispatch` 的 GitHub PAT |

Token 权限要求：

- **Fine-grained PAT（推荐）**：仅需对 `BenedictKing/homebrew-ccx` 的 **Contents: Read and write** 权限。
- **Classic PAT**：需要 `repo` scope。
- 不使用 `GITHUB_TOKEN`，因为它无法触发其他仓库的 dispatch。

## 故障排查

```bash
# 检查 Tap 是否已添加
brew tap

# 查看 Cask 信息
brew info --cask BenedictKing/ccx/ccx-desktop

# 检查 livecheck 是否能检测到新版本
brew livecheck --cask BenedictKing/ccx/ccx-desktop

# 强制重新安装
brew reinstall --cask ccx-desktop
```
