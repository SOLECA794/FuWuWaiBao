# 智能互动讲课系统 - 后端接口文档

> 本文档描述当前仓库中 **已经实现并可运行** 的 Go 后端接口。  
> 若关注统一风格接口，请同时参考 [API_DESIGN_V2.md](API_DESIGN_V2.md)。  
> 若关注当前联调可直接调用的接口，请优先参考 [当前运行接口清单.md](当前运行接口清单.md)。

---

## 1. 基础信息

- 基础地址：`http://localhost:18080`
- 旧版内部接口前缀：`/api`
- 统一内部接口前缀：`/api/v1`
- 开放适配接口前缀：`/api/v1`
- 返回格式：JSON（SSE 接口除外）
- 字符编码：UTF-8

### 1.1 通用响应格式（内部接口）

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 1.2 通用响应格式（开放适配接口）

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": {},
  "requestId": "req_xxx"
}
```

### 1.3 健康检查

- **接口**：`GET /health`
- **功能**：检查 Go 后端服务是否正常启动

---

## 2. 教师端旧版接口（兼容保留）

### 2.1 获取课件列表
- **接口**：`GET /api/teacher/courseware-list`
- **功能**：获取课件列表

### 2.2 上传课件
- **接口**：`POST /api/teacher/upload-courseware`
- **Content-Type**：`multipart/form-data`
- **参数**：
  - `file`：PDF / PPT / PPTX 文件
  - `title`：课件标题（可选）

### 2.3 删除课件
- **接口**：`DELETE /api/teacher/courseware/{courseId}`

### 2.4 发布课件
- **接口**：`POST /api/teacher/publish-courseware`
- **请求体**：
```json
{
  "courseId": "xxx",
  "scope": "all"
}
```

### 2.5 获取讲稿
- **接口**：`GET /api/teacher/script/{courseId}/{page}`

### 2.6 保存讲稿
- **接口**：`POST /api/teacher/script/save`
- **请求体**：
```json
{
  "courseId": "xxx",
  "page": 1,
  "content": "讲稿内容"
}
```

### 2.7 AI 生成讲稿
- **接口**：`POST /api/teacher/ai-generate-script`
- **请求体**：
```json
{
  "courseId": "xxx",
  "page": 1,
  "courseName": "Go语言基础教程"
}
```

### 2.8 教师学情统计
- **接口**：`GET /api/teacher/student-stats/{courseId}`
- **说明**：返回真实统计结果；当前版本已移除默认补词逻辑

### 2.9 提问记录
- **接口**：`GET /api/teacher/question-records/{courseId}?page=1&pageSize=20`

### 2.10 卡点数据
- **接口**：`GET /api/teacher/card-data/{courseId}`
- **说明**：当前版本不再按页码推导伪造停留时长，若无日志则返回零值数据

---

## 3. 学生端旧版接口（兼容保留 + 补齐）

### 3.1 获取学生可学习课件列表
- **接口**：`GET /api/student/courseware-list`
- **说明**：优先返回已发布课件；若无发布课件，则回退返回全部课件

### 3.2 开始学习会话
- **接口**：`POST /api/student/session/start`
- **功能**：初始化学习会话
- **请求体**：
```json
{
  "userId": "student001",
  "courseId": "course001"
}
```

### 3.3 上报播放进度
- **接口**：`POST /api/student/progress/update`
- **请求体**：
```json
{
  "sessionId": "sess_xxx",
  "userId": "student001",
  "courseId": "course001",
  "currentPage": 3,
  "currentNodeId": "p3_n2"
}
```

### 3.4 获取学生播放脚本
- **接口**：`GET /api/student/script/{courseId}/{page}`
- **说明**：返回 `nodes[]` 结构化脚本，供前端播放和续接定位

### 3.5 SSE 问答流
- **接口**：`POST /api/student/qa/stream`
- **响应类型**：`text/event-stream; charset=utf-8`
- **请求体**：
```json
{
  "sessionId": "sess_xxx",
  "courseId": "course001",
  "page": 3,
  "nodeId": "p3_n2",
  "question": "这里的公式是什么意思？"
}
```

#### SSE 事件
- `token`：增量文本
- `sentence`：完整句子
- `final`：结束结构化信息

### 3.6 获取课件页内容
- **接口**：`POST /api/student/courseware/page`

### 3.7 普通问答
- **接口**：`POST /api/student/ai/question`

### 3.8 溯源问答
- **接口**：`POST /api/student/ai/traceQuestion`

### 3.9 学习数据统计
- **接口**：`GET /api/student/studyData?studentId=xxx&courseId=xxx`

### 3.10 获取断点
- **接口**：`GET /api/student/breakpoint?studentId=xxx&courseId=xxx`

### 3.11 更新断点
- **接口**：`PUT /api/student/breakpoint`

### 3.12 保存笔记
- **接口**：`POST /api/student/saveNote`

---

## 4. 薄弱点与知识点接口（旧版兼容）

### 4.1 获取薄弱点列表
- **接口**：`GET /api/weakPoint/getList?studentId=xxx&courseId=xxx`
- **说明**：当前版本只返回真实统计结果，不再强行补“暂无薄弱点”假数据

### 4.2 获取薄弱点讲解
- **接口**：`POST /api/weakPoint/getExplain`

### 4.3 生成薄弱点习题
- **接口**：`POST /api/weakPoint/getTest`

### 4.4 校验答案
- **接口**：`POST /api/weakPoint/checkAnswer`

### 4.5 解析知识点树
- **接口**：`POST /api/ai/parseKnowledge`

---

## 5. 课件预览接口

### 5.1 获取预览页图片
- **接口**：`GET /api/courseware/{courseId}/page/{pageNum}`
- **成功行为**：302 重定向到预览图 URL
- **失败行为**：返回 404 JSON

---

## 6. `/api/v1` 统一内部接口

> 以下接口已经在后端路由中真实注册，前端新代码优先使用这一套。

### 6.1 教师端

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/teacher/coursewares` | 获取课件列表 |
| POST | `/api/v1/teacher/coursewares/upload` | 上传课件 |
| DELETE | `/api/v1/teacher/coursewares/{courseId}` | 删除课件 |
| GET | `/api/v1/teacher/coursewares/{courseId}/scripts/{pageNum}` | 获取讲稿 |
| PUT | `/api/v1/teacher/coursewares/{courseId}/scripts/{pageNum}` | 更新讲稿 |
| POST | `/api/v1/teacher/coursewares/{courseId}/scripts/ai-generate` | AI 生成讲稿 |
| POST | `/api/v1/teacher/coursewares/{courseId}/publish` | 发布课件 |
| GET | `/api/v1/teacher/coursewares/{courseId}/stats` | 学情统计 |
| GET | `/api/v1/teacher/coursewares/{courseId}/questions` | 提问记录 |
| GET | `/api/v1/teacher/coursewares/{courseId}/card-data` | 卡点数据 |

### 6.2 AI 学伴

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/ai/coursewares/{courseId}/knowledge-graph` | 获取知识图谱 |
| POST | `/api/v1/ai/coursewares/{courseId}/ask` | 学生问答 |

### 6.3 学生端

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/student/coursewares/{courseId}/weak-points` | 获取薄弱点 |
| GET | `/api/v1/student/weak-points/{weakPointId}/explain` | 获取讲解 |
| POST | `/api/v1/student/weak-points/{weakPointId}/generate-test` | 生成习题 |
| POST | `/api/v1/student/tests/{questionId}/check` | 校验答案 |
| GET | `/api/v1/student/coursewares/{courseId}/breakpoint` | 获取断点 |
| PUT | `/api/v1/student/coursewares/{courseId}/breakpoint` | 更新断点 |
| POST | `/api/v1/student/coursewares/{courseId}/notes` | 保存笔记 |
| GET | `/api/v1/student/coursewares/{courseId}/stats` | 学习统计 |

---

## 7. 开放适配层接口（泛雅对接）

> 当配置环境变量 `OPEN_API_STATIC_KEY` 时，以下接口会启用 `enc + time` 校验。  
> 若未配置，则默认放行，便于本地联调。

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/v1/lesson/parse` | 课件解析 |
| POST | `/api/v1/lesson/generateScript` | 脚本生成 |
| POST | `/api/v1/lesson/generateAudio` | 音频生成（当前为骨架实现） |
| POST | `/api/v1/qa/interact` | 问答交互 |
| POST | `/api/v1/qa/voiceToText` | 语音转文本（当前为骨架实现） |
| POST | `/api/v1/progress/track` | 进度跟踪 |
| POST | `/api/v1/progress/adjust` | 节奏调整 |
| POST | `/api/v1/platform/syncCourse` | 同步课程 |
| POST | `/api/v1/platform/syncUser` | 同步用户 |

---

## 8. AI 引擎内部依赖接口

> 以下接口由 Python `ai_engine` 服务提供，Go 后端通过 HTTP 调用。

| 方法 | 路径 | 用途 |
|---|---|---|
| POST | `/generate-script` | 生成讲稿 |
| POST | `/ask-with-context` | 上下文问答 |
| POST | `/parse-knowledge` | 知识点树解析 |
| POST | `/upload` | 文档上传解析（当前 Go 未直接调用） |
| GET | `/lessons/{doc_id}` | 获取生成课程（当前 Go 未直接调用） |
| POST | `/chat` | 原生问答（当前 Go 未直接调用） |
| GET | `/health` | AI 服务健康检查 |

---

## 9. 更新记录

| 日期 | 版本 | 内容 |
|---|---|---|
| 2026-03-06 | v2.0 | 文档同步到当前实际实现，补充 `/api/v1` 与开放适配层 |
