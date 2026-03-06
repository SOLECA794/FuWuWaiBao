# 泛雅 AI 智课系统

一个集成 Go 后端、Python AI 引擎与多套前端的教学平台样例工程。

## 概览
- `backend/`：Go 后端服务（Gin + GORM + Redis + MinIO）
- `ai_engine/`：Python AI 服务（PDF/PPTX 解析、讲稿生成、问答）
- `frontend/`：包含学生、教师等前端应用（Vue）
- `docs/`：项目设计与 API 文档（已生成）
- `uploads/`：上传的临时文件存放目录

> 本仓库为教学/演示用途，目录与运行方式可根据部署环境调整。

## 快速开始（本地开发）

先决条件：已安装 Docker、Go、Node.js、npm/yarn、以及 Python（推荐通过 conda 管理虚拟环境）。

1) 启动基础服务（Postgres/Redis/MinIO）

```bash
docker-compose -f backend/docker-compose.yml up -d
```

2) 启动 AI 引擎（Python）

```bash
cd ai_engine
conda activate fuww_ai   # 或使用你自己的虚拟环境
pip install -r requirements.txt
python main.py
```

3) 运行后端（Go）

```bash
cp backend/config/config.yaml.example backend/config/config.yaml
cd backend
go run ./cmd/api
```

4) 启动前端示例（学生端）

```bash
cd frontend/student
npm install
npm run serve
```

## 文档

详细接口与设计文档保存在 `docs/` 目录，请参考该目录下的 Markdown 文件以获取 API 说明、联调清单与系统架构等内容。

## 开发建议与调试
- 开发 AI 模块时，优先在 `ai_engine/` 下使用提供的测试数据（`ai_engine/tests/mocks/`）进行快速迭代。
- 后端配置文件位于 `backend/config/`，上线前请检查数据库与 MinIO 的连接配置。

## 常见命令汇总

- 启动全部依赖：

```bash
docker-compose -f backend/docker-compose.yml up -d
```

- 启动 AI 服务：

```bash
cd ai_engine
python main.py
```

- 运行后端：

```bash
cd backend
go run ./cmd/api
```

## 贡献与联系
如果你希望贡献代码或有问题，请在仓库中创建 issue 或联系维护者。文档已放在 `docs/` 中，优先阅读后提交改动。

---
更新于：项目结构与文档已同步到 `docs/`，README 已适配最新目录。

