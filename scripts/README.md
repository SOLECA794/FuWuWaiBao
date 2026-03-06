# 启动脚本说明

仓库根目录下提供了方便在 Windows (PowerShell) 上一键启动/停止本项目前端、后端、AI 引擎和依赖服务的脚本。

文件：
- `scripts/start-all.ps1`：一键启动脚本，会执行（默认）`docker-compose -f backend/docker-compose.yml up -d`，然后启动 AI 引擎（通过 `conda run -n fuww_ai python ai_engine/main.py`）、Go 后端和两个前端（使用 `npm run dev`）。
- `scripts/stop-all.ps1`：停止由 `start-all.ps1` 启动的进程并执行 `docker-compose -f backend/docker-compose.yml down`。

使用方法（在 PowerShell 中，以仓库根目录为当前目录运行）：

启动（含 docker-compose）：
```powershell
.\scripts\start-all.ps1
```

若想跳过 docker-compose（例如你已经用容器外的服务），可以：
```powershell
.\scripts\start-all.ps1 -SkipDocker
```

停止：
```powershell
.\scripts\stop-all.ps1
```

注意：
- 这些脚本假设系统上已安装并可用的命令：`docker-compose`、`conda`、`go`、`npm`。请确保将 `conda` 命令加入到 PATH，或在运行脚本前从 Anaconda/Miniconda 的 PowerShell 快捷方式打开。
- AI 引擎使用的 conda 环境默认名为 `fuww_ai`（与仓库内说明一致）。如果你的环境名不同，请修改 `scripts/start-all.ps1` 中对应的 `-n` 参数或在命令前激活正确的 conda 环境。
- 日志保存在 `logs/` 目录下（脚本会自动创建），PID 信息保存到 `scripts/pids.json`，供 `stop-all.ps1` 使用。

如需我在当前机器上运行并验证这些脚本，我可以继续执行。需要我现在运行吗？
