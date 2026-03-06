<<<<<<< HEAD
# 泛雅 AI 智课系统 (Smart Teaching)

本项目是一个包含 Go 后端、FastAPI AI 引擎以及三个前端应用的综合系统。

## 项目结构
- **backend/**: Go 企业级后端 (Gin + GORM + Redis + MinIO)
- **ai_engine/**: Python AI 核心服务 (PPT/PDF 解析 + 讲稿生成 + 问答)
- **frontend/**: 包含所有前端应用
    - **student/**: 学生端前端 (Vue 3 + Element Plus)
    - **teacher/**: 教师端前端 (Vue 3)
    - **ai-tool/**: AI 指令辅助/测试前端
- **docs/**: 项目接口与设计文档
- **uploads/**: 用于存储上传的临时文档 (PPT/PDF)

## 快速开始

### 1. 基础设施 (Postgres, Redis, MinIO)
```bash
# 启动数据库/存储依赖
docker-compose -f backend/docker-compose.yml up -d
```

### 2. AI 引擎 (Python)
```bash
# 进入目录并切换 conda 环境
conda activate fuww_ai
pip install -r requirements.txt
python ai_engine/main.py
```

### 3. 后端服务 (Go)
```bash
# 拷贝配置示例
cp backend/config/config.yaml.example backend/config/config.yaml
# 运行
go run backend/cmd/api/main.go
```

### 4. 前端 (学生端示例)
```bash
cd student-frontend
npm install
npm run dev
```

## 功能清单
- [x] 多模态（PDF/PPTX）分页解析
- [x] 智能讲稿生成
- [x] 实时对话与案例重讲
- [ ] 知识点地图展示

=======
# Vue 3 + Vite

This template should help get you started developing with Vue 3 in Vite. The template uses Vue 3 `<script setup>` SFCs, check out the [script setup docs](https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup) to learn more.

Learn more about IDE Support for Vue in the [Vue Docs Scaling up Guide](https://vuejs.org/guide/scaling-up/tooling.html#ide-support).
>>>>>>> b79727f64ad1860d8e9dc554eec4fdaef2859d48
