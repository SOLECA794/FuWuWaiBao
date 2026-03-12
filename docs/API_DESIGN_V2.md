# 智能学习平台 API 接口文档（v2.1 已落地版）

> 本文档用于描述当前仓库中 `/api/v1` 统一接口体系与开放适配层的**已落地实现**。  
> 与旧版兼容接口的差异，请参考 [后端API.md](后端API.md) 与 [当前运行接口清单.md](当前运行接口清单.md)。

---

## 一、全局规范

1. **基础路径**：`/api/v1`
2. **数据格式**：`application/json; charset=utf-8`
3. **资源 ID**：统一使用 `String`
4. **内部接口响应结构**：
   ```json
   {
     "code": 200,
     "message": "请求成功",
     "data": {}
   }
   ```
5. **开放适配接口响应结构**：
   ```json
   {
     "code": 200,
     "msg": "操作成功",
     "data": {},
     "requestId": "req_xxx"
   }
   ```
6. **签名校验**：开放适配接口在配置 `OPEN_API_STATIC_KEY` 时启用 `enc + time` 校验。

---

## 二、教师端接口

> 基础路径：`/api/v1/teacher/coursewares`

### 2.1 获取课件列表
- **接口**：`GET /api/v1/teacher/coursewares`
- **状态**：已实现

### 2.2 上传课件
- **接口**：`POST /api/v1/teacher/coursewares/upload`
- **Content-Type**：`multipart/form-data`
- **参数**：`file`、`title`
- **状态**：已实现

### 2.3 删除课件
- **接口**：`DELETE /api/v1/teacher/coursewares/{courseId}`
- **状态**：已实现

### 2.4 获取页面讲稿
- **接口**：`GET /api/v1/teacher/coursewares/{courseId}/scripts/{pageNum}`
- **状态**：已实现

### 2.5 更新页面讲稿
- **接口**：`PUT /api/v1/teacher/coursewares/{courseId}/scripts/{pageNum}`
- **请求体**：
```json
{
  "content": "修改后的讲稿内容"
}
```
- **状态**：已实现

### 2.6 AI 生成讲稿
- **接口**：`POST /api/v1/teacher/coursewares/{courseId}/scripts/ai-generate`
- **请求体**：
```json
{
  "pageNum": 1,
  "mode": "llm"
}
```
- **状态**：已实现

### 2.7 发布课件
- **接口**：`POST /api/v1/teacher/coursewares/{courseId}/publish`
- **请求体**：
```json
{
  "scope": "all"
}
```
- **状态**：已实现

### 2.8 教师学情统计
- **接口**：`GET /api/v1/teacher/coursewares/{courseId}/stats`
- **状态**：已实现
- **备注**：当前返回真实统计数据，不再填充默认关键词

### 2.9 历史提问记录
- **接口**：`GET /api/v1/teacher/coursewares/{courseId}/questions?page=1&pageSize=20`
- **状态**：已实现

### 2.10 卡点数据
- **接口**：`GET /api/v1/teacher/coursewares/{courseId}/card-data`
- **状态**：已实现
- **备注**：无行为日志时返回零值，不再构造估算型伪数据

---

## 三、AI 学伴与互动答疑

> 基础路径：`/api/v1/ai/coursewares`

### 3.1 获取知识图谱
- **接口**：`GET /api/v1/ai/coursewares/{courseId}/knowledge-graph`
- **状态**：已实现
- **说明**：当前从 `knowledge_points` 表构建树结构；若无数据则返回空数组

### 3.2 智能答疑
- **接口**：`POST /api/v1/ai/coursewares/{courseId}/ask`
- **状态**：已实现
- **请求体**：
```json
{
  "pageNum": 3,
  "type": "text",
  "studentId": "student001",
  "question": "这里提到的 fillna 怎么用？",
  "tracePoint": {
    "x": 200,
    "y": 150
  }
}
```
- **说明**：支持普通提问和带 `tracePoint` 的定位提问

---

## 四、薄弱点诊断与练习

> 基础路径：`/api/v1/student`

### 4.1 获取薄弱点列表
- **接口**：`GET /api/v1/student/coursewares/{courseId}/weak-points?studentId=student001`
- **状态**：已实现

### 4.2 获取薄弱点讲解
- **接口**：`GET /api/v1/student/weak-points/{weakPointId}/explain?name=第3页知识点`
- **状态**：已实现

### 4.3 生成练习题
- **接口**：`POST /api/v1/student/weak-points/{weakPointId}/generate-test`
- **状态**：已实现
- **请求体**：
```json
{
  "studentId": "student001",
  "questionType": "single",
  "weakPointName": "第3页知识点"
}
```

### 4.4 校验答案
- **接口**：`POST /api/v1/student/tests/{questionId}/check`
- **状态**：已实现
- **请求体**：
```json
{
  "studentId": "student001",
  "userAnswer": "A"
}
```

---

## 五、学习过程数据

> 基础路径：`/api/v1/student/coursewares`

### 5.1 获取学习断点
- **接口**：`GET /api/v1/student/coursewares/{courseId}/breakpoint?studentId=student001`
- **状态**：已实现

### 5.2 更新学习断点
- **接口**：`PUT /api/v1/student/coursewares/{courseId}/breakpoint`
- **状态**：已实现
- **请求体**：
```json
{
  "studentId": "student001",
  "pageNum": 5
}
```

### 5.3 保存随堂笔记
- **接口**：`POST /api/v1/student/coursewares/{courseId}/notes`
- **状态**：已实现
- **请求体**：
```json
{
  "studentId": "student001",
  "pageNum": 3,
  "content": "我的笔记内容...",
  "x": 0,
  "y": 0
}
```

### 5.4 学生学习统计
- **接口**：`GET /api/v1/student/coursewares/{courseId}/stats?studentId=student001`
- **状态**：已实现

---

## 六、学生端播放与会话接口（当前主运行接口）

> 以下接口不是 `/api/v1` 规范的一部分，但已经真实实现，并被学生端页面直接使用。

> 说明：为配合页面层完全围绕 `v1` 组织，当前这些能力也已补充 `v1` 路由别名。

### 6.0 学生课件列表（v1）
- **接口**：`GET /api/v1/student/coursewares`
- **状态**：已实现

### 6.0.1 开始学习会话（v1）
- **接口**：`POST /api/v1/student/sessions`
- **状态**：已实现

### 6.0.2 上报学习进度（v1）
- **接口**：`POST /api/v1/student/sessions/progress`
- **状态**：已实现
- **说明**：当前支持 `currentTimeSec`，用于持久化播放秒级游标，而不再只记录页码/节点

### 6.0.3 获取学生播放脚本（v1）
- **接口**：`GET /api/v1/student/coursewares/{courseId}/scripts/{pageNum}`
- **状态**：已实现
- **说明**：返回节点级时间轴、`audio_meta` 与 `playback_mode`，可同时支持纯时长驱动与未来音频驱动播放

### 6.0.4 获取页面播放音频元数据（v1）
- **接口**：`GET /api/v1/student/coursewares/{courseId}/scripts/{pageNum}/audio`
- **状态**：已实现
- **说明**：返回当前页的占位音频包元数据，包括 `sections`、`total_duration_sec`、`playback_mode`

### 6.0.5 SSE 问答流（v1）
- **接口**：`POST /api/v1/student/qa/stream`
- **状态**：已实现

### 6.1 获取学生课件列表
- **接口**：`GET /api/student/courseware-list`
- **状态**：已实现


### 6.2 开始学习会话
- **接口**：`POST /api/student/session/start`
- **状态**：已实现

### 6.3 上报播放进度
- **接口**：`POST /api/student/progress/update`
- **状态**：已实现

### 6.4 获取学生播放脚本
- **接口**：`GET /api/student/script/{courseId}/{page}`
- **状态**：已实现

### 6.4.1 获取学生播放音频元数据
- **接口**：`GET /api/student/script/{courseId}/{page}/audio`
- **状态**：已实现

### 6.5 SSE 问答流
- **接口**：`POST /api/student/qa/stream`
- **状态**：已实现
- **Content-Type**：`text/event-stream; charset=utf-8`

#### 事件约定
- `token`：逐片文本
- `sentence`：整句文本
- `final`：结束帧，包含 `need_reteach`、`resume_page`、`resume_node_id`、`resume_sec`

---

## 七、开放适配层（OpenAPI Adapter）

### 7.1 设计目标
对齐泛雅平台接口风格，提供“课件解析 → 脚本生成 → 问答交互 → 进度适配”的开放调用入口。

### 7.2 签名规则
- 参数：`enc`、`time`
- 算法：`enc = MD5(参数有序拼接 + staticKey + time)`
- `time` 格式：`yyyy-MM-ddHH:mm:ss`
- 仅当环境变量 `OPEN_API_STATIC_KEY` 存在时启用强校验

### 7.3 已实现开放接口

| 模块 | 接口 | 方法 | 状态 | 说明 |
|---|---|---|---|---|
| 课件解析 | `/api/v1/lesson/parse` | POST | 已实现 | 保存课程基础信息并返回 `parseId` |
| 脚本生成 | `/api/v1/lesson/generateScript` | POST | 已实现 | 调用 AI 生成脚本结构 |
| 语音合成 | `/api/v1/lesson/generateAudio` | POST | 骨架已实现 | 当前返回占位音频资产元数据，可为后续真实 TTS 对接保留统一结构 |
| 问答交互 | `/api/v1/qa/interact` | POST | 已实现 | 基于课件页内容答疑 |
| 语音识别 | `/api/v1/qa/voiceToText` | POST | 骨架已实现 | 当前支持占位文本返回 |
| 进度追踪 | `/api/v1/progress/track` | POST | 已实现 | 记录进度到内存/Redis/DB |
| 节奏调整 | `/api/v1/progress/adjust` | POST | 已实现 | 返回 `adjustPlan` |
| 平台用户列表 | `/api/v1/platform/users` | GET | 已实现 | 支持按角色、组织、关键字过滤平台用户 |
| 平台用户详情 | `/api/v1/platform/users/{userId}` | GET | 已实现 | 返回用户基础信息及关联课程、班级、选课统计 |
| 平台课程列表 | `/api/v1/platform/courses` | GET | 已实现 | 支持按教师、状态、组织过滤课程 |
| 平台课程创建 | `/api/v1/platform/courses` | POST | 已实现 | 支持通过教师、学期、组织信息创建平台课程 |
| 平台课程详情 | `/api/v1/platform/courses/{courseId}` | GET | 已实现 | 返回教师、班级、成员和选课统计 |
| 平台课程更新 | `/api/v1/platform/courses/{courseId}` | PUT | 已实现 | 支持按字段更新课程元数据 |
| 平台课程删除 | `/api/v1/platform/courses/{courseId}` | DELETE | 已实现 | 删除课程并级联清理关联班级、选课 |
| 平台班级列表 | `/api/v1/platform/classes` | GET | 已实现 | 支持按课程、教师、状态过滤班级 |
| 平台班级创建 | `/api/v1/platform/classes` | POST | 已实现 | 支持指定课程创建教学班 |
| 平台班级详情 | `/api/v1/platform/classes/{classId}` | GET | 已实现 | 返回课程、教师和班级成员数据 |
| 平台班级更新 | `/api/v1/platform/classes/{classId}` | PUT | 已实现 | 支持调整归属课程、教师、容量和状态 |
| 平台班级删除 | `/api/v1/platform/classes/{classId}` | DELETE | 已实现 | 删除班级并清理班级下选课关系 |
| 平台选课列表 | `/api/v1/platform/enrollments` | GET | 已实现 | 支持按课程、班级、用户、角色、状态过滤 |
| 平台选课创建 | `/api/v1/platform/enrollments` | POST | 已实现 | 支持手动创建学生/教师选课关系 |
| 平台选课详情 | `/api/v1/platform/enrollments/{enrollmentId}` | GET | 已实现 | 返回选课关系及用户/课程/班级摘要 |
| 平台选课更新 | `/api/v1/platform/enrollments/{enrollmentId}` | PUT | 已实现 | 支持调整选课角色、状态与关联实体 |
| 平台选课删除 | `/api/v1/platform/enrollments/{enrollmentId}` | DELETE | 已实现 | 支持手动解除选课关系 |
| 平台同步 | `/api/v1/platform/syncCourse` | POST | 已实现 | 已落到平台课程/班级/教师关系模型，并返回 classes、teachers、courseMeta |
| 平台同步 | `/api/v1/platform/syncUser` | POST | 已实现 | 已落到平台用户/选课关系模型，并返回 linkedCourses、userProfile |
| 平台总览 | `/api/v1/platform/overview` | GET | 已实现 | 返回平台用户、课程、班级、选课统计与最近同步记录 |

### 7.4 当前 TTS 占位层约定

- `GenerateAudio` 已在 Go 后端与 AI 引擎接口层预留统一请求/响应结构。
- 数据层新增了 `AudioAsset`、`CoursePage.audio_*`、`TeachingNode.audio_*`、`TeachingNode.tts_status` 等字段，用于沉淀音频状态与时间轴切片。
- 当前阶段由后端占位逻辑生成结构化元数据，不强依赖真实 TTS 服务，因此可以先完成前端播放协议与断点续接。
- 当任一节点具备真实 `audio_url` 且 `tts_status=ready` 时，播放模式自动切换为 `audio_timeline`；否则保持 `duration_timeline`。

### 7.5 平台级模型落点

- `PlatformUser`：平台统一用户实体，承载外部用户 ID、角色、姓名、组织信息。
- `TeachingCourse`：面向教学运营的课程实体，不再与单个课件文件强绑定。
- `CourseClass`：教学班/行政班映射，用于班级发布、统计与权限控制。
- `CourseEnrollment`：选课/授课关系，承载学生、教师、课程、班级之间的关联。
- 现有 `Course` 继续承担课件资源与解析结果实体，后续通过外键或映射表挂接到 `TeachingCourse`。

---

## 八、与 Python AI 引擎的内部映射

| Go 后端能力 | AI 引擎接口 | 说明 |
|---|---|---|
| 讲稿生成 | `POST /generate-script` | 已接入 |
| 上下文问答 | `POST /ask-with-context` | 已接入 |
| 知识点解析 | `POST /parse-knowledge` | 已接入 |
| 音频生成 | `POST /generate-audio` | 接口已预留，默认由后端占位层兜底 |

其余 AI 引擎原生接口：
- `POST /upload`
- `GET /lessons/{doc_id}`
- `POST /chat`
- `GET /health`

---

## 九、状态说明

本次版本变更后：
- `/api/v1` 路由已真实注册
- 学生端会话 / SSE 能力已落地
- 旧版接口保留用于兼容
- 开放适配层已具备可调用骨架
- 播放链路已支持 `resume_sec`、音频元数据和秒级续接
- 平台级核心模型已入库，开放适配层已支持总览、列表、详情和同步接口
- 教师学情 / 卡点 / 薄弱点已移除伪造默认数据

---

## 十、更新记录

| 日期 | 版本 | 更新内容 |
|---|---|---|
| 2026-03-11 | v2.2 | 补充学生端播放音频元数据、`resume_sec` 续接协议、TTS 占位层与平台级模型说明 |
| 2026-03-06 | v2.1 | 文档与当前实现同步，标记已实现接口与开放适配层状态 |
