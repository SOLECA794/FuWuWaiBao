# 前端对接说明（当前实现版）

> 本文档用于指导前端按当前仓库真实接口接入教师端与学生端页面。

---

## 1. 以哪份文档为准？

按优先级从高到低：
1. [当前运行接口清单.md](当前运行接口清单.md)
2. [后端API.md](后端API.md)
3. [API_DESIGN_V2.md](API_DESIGN_V2.md)
4. [泛雅API规范.md](泛雅API规范.md)

说明：
- 当前页面联调，以“当前运行接口清单”优先。
- `/api/v1` 是推荐主接口。
- `/api/teacher/*`、`/api/student/*`、`/api/weakPoint/*` 仍保留兼容。

---

## 2. 当前前端服务文件

### 2.1 学生端
- 入口聚合：`frontend/student/src/services/v1/index.js`
- 资源模块：
  - `frontend/student/src/services/v1/coursewareApi.js`
  - `frontend/student/src/services/v1/sessionApi.js`
  - `frontend/student/src/services/v1/qaApi.js`
  - `frontend/student/src/services/v1/weakPointApi.js`
  - `frontend/student/src/services/v1/knowledgeApi.js`
- API 基址配置：`frontend/student/src/config/api.js`
- 默认后端地址：`http://localhost:18080`

### 2.2 教师端
- 入口聚合：`frontend/teacher/src/services/v1/index.js`
- 资源模块：
  - `frontend/teacher/src/services/v1/coursewareApi.js`
  - `frontend/teacher/src/services/v1/analyticsApi.js`
- API 基址配置：`frontend/teacher/src/config/api.js`
- 默认后端地址：`http://localhost:18080`

---

## 3. 教师端当前使用接口

| 功能 | 方法 | 路径 |
|---|---|---|
| 健康检查 | GET | `/health` |
| 课件列表 | GET | `/api/v1/teacher/coursewares` |
| 删除课件 | DELETE | `/api/v1/teacher/coursewares/{courseId}` |
| 获取讲稿 | GET | `/api/v1/teacher/coursewares/{courseId}/scripts/{page}` |
| 保存讲稿 | PUT | `/api/v1/teacher/coursewares/{courseId}/scripts/{page}` |
| AI 生成讲稿 | POST | `/api/v1/teacher/coursewares/{courseId}/scripts/ai-generate` |
| 上传课件 | POST | `/api/v1/teacher/coursewares/upload` |
| 发布课件 | POST | `/api/v1/teacher/coursewares/{courseId}/publish` |
| 学情统计 | GET | `/api/v1/teacher/coursewares/{courseId}/stats` |
| 卡点数据 | GET | `/api/v1/teacher/coursewares/{courseId}/card-data` |
| 提问记录 | GET | `/api/v1/teacher/coursewares/{courseId}/questions?page=1&pageSize=100` |

---

## 4. 学生端当前使用接口

| 功能 | 方法 | 路径 |
|---|---|---|
| 健康检查 | GET | `/health` |
| 课件列表 | GET | `/api/student/courseware-list` |
| 开始会话 | POST | `/api/student/session/start` |
| 获取页面脚本 | GET | `/api/student/script/{courseId}/{page}` |
| 上报进度 | POST | `/api/student/progress/update` |
| SSE 问答流 | POST | `/api/student/qa/stream` |
| 获取断点 | GET | `/api/v1/student/coursewares/{courseId}/breakpoint?studentId=...` |
| 保存断点 | PUT | `/api/v1/student/coursewares/{courseId}/breakpoint` |
| 获取学习统计 | GET | `/api/v1/student/coursewares/{courseId}/stats?studentId=...` |
| 获取薄弱点 | GET | `/api/v1/student/coursewares/{courseId}/weak-points?studentId=...` |
| 获取讲解 | GET | `/api/v1/student/weak-points/{weakPointId}/explain?name=...` |
| 生成习题 | POST | `/api/v1/student/weak-points/{weakPointId}/generate-test` |
| 校验答案 | POST | `/api/v1/student/tests/{questionId}/check` |
| 知识点解析 | POST | `/api/ai/parseKnowledge` |

---

## 5. 学生端最小闭环（当前已落地）

学生端当前最小闭环已是：
1. 获取课件列表
2. 创建学习会话
3. 获取某页结构化脚本（`nodes[]`）
4. 更新当前播放进度（`currentPage/currentNodeId`）
5. 发起 SSE 问答流
6. 根据 `final.resume_node_id` 恢复播放节点
7. 读取断点 / 保存断点
8. 查看薄弱点 / 做题 / 判题

---

## 6. SSE 约定（当前已实现）

### 6.1 事件类型
- `token`：增量文本，用于打字机效果
- `sentence`：完整句子
- `final`：结束帧

### 6.2 `final` 典型字段
```json
{
  "need_reteach": false,
  "source_page": 5,
  "resume_page": 5,
  "resume_node_id": "p5_n13"
}
```

### 6.3 前端建议维护的状态
- `sessionId`
- `currentPage`
- `currentNodeId`
- `aiReply`
- `qaHistory`
- `resumeNodeId`

---

## 7. 当前前端实现注意事项

### 7.1 学生端
- 学生端已不再复用教师端课件列表接口。
- 学生端问答优先走 `/api/student/qa/stream`。
- 学生端脚本播放定位依赖 `node_id` 格式：`p{page}_n{index}`。

### 7.2 教师端
- 教师端已切换到 `/api/v1/teacher/coursewares/*`。
- 讲稿保存已从旧版 `POST /api/teacher/script/save` 切换为 `PUT /api/v1/.../scripts/{pageNum}`。

---

## 8. 若要继续扩展前端

优先建议：
1. 把学生端“溯源定位”真正接到 `/api/v1/ai/coursewares/{courseId}/ask` 的 `tracePoint` 字段。
2. 把学生端页面播放逻辑从“页级”继续细化到“node 级”。
3. 为教师端增加 `/knowledge-graph` 的可视化展示。
4. 后续若开放给平台嵌入，再接开放适配层而不是直接暴露内部接口。

---

## 9. 更新记录

| 日期 | 版本 | 内容 |
|---|---|---|
| 2026-03-06 | v2.1 | 前端服务层改为 `services/v1` 资源模块结构，页面层围绕 v1 API 组织 |
