# 智能互动讲课系统后端接口说明

本文档以当前仓库实现为准，目标是保留现有统一接口设计，同时兼容旧版联调接口。

## 基础信息

- 服务地址：`http://localhost:18080`
- 健康检查：`GET /health`
- 旧版兼容接口前缀：`/api`
- 统一接口前缀：`/api/v1`
- 返回编码：UTF-8

---

## 一、教师端统一接口

### 1. 课件管理

- `GET /api/v1/teacher/coursewares`
- `POST /api/v1/teacher/coursewares/upload`
- `DELETE /api/v1/teacher/coursewares/:courseId`
- `POST /api/v1/teacher/coursewares/:courseId/publish`

### 2. 讲稿编辑

- `GET /api/v1/teacher/coursewares/:courseId/scripts/:pageNum`
- `PUT /api/v1/teacher/coursewares/:courseId/scripts/:pageNum`
- `POST /api/v1/teacher/coursewares/:courseId/scripts/ai-generate`

### 3. 教学分析

- `GET /api/v1/teacher/coursewares/:courseId/stats`
- `GET /api/v1/teacher/coursewares/:courseId/questions`
- `GET /api/v1/teacher/coursewares/:courseId/card-data`
- `GET /api/v1/teacher/coursewares/:courseId/node-insights`（节点洞察：提问数/笔记数/练习正确率/重讲率/7日趋势）

---

## 二、学生端统一接口

### 1. 学习过程

- `GET /api/v1/student/coursewares`
- `POST /api/v1/student/sessions`
- `POST /api/v1/student/sessions/progress`
- `GET /api/v1/student/coursewares/:courseId/scripts/:pageNum`
- `POST /api/v1/student/qa/stream`

### 2. 学习记录

- `GET /api/v1/student/coursewares/:courseId/breakpoint`
- `PUT /api/v1/student/coursewares/:courseId/breakpoint`
- `POST /api/v1/student/coursewares/:courseId/notes`
- `GET /api/v1/student/notes`（支持 `studentId/courseId/pageNum/page/pageSize`）
- `POST /api/v1/student/favorites`
- `GET /api/v1/student/favorites`（支持分页）
- `DELETE /api/v1/student/favorites/:favoriteId`（需 `studentId` 归属校验）
- `GET /api/v1/student/coursewares/:courseId/stats`

### 3. 薄弱点与练习

- `GET /api/v1/student/coursewares/:courseId/weak-points`
- `GET /api/v1/student/weak-points/:weakPointId/explain`
- `POST /api/v1/student/weak-points/:weakPointId/generate-test`
- `POST /api/v1/student/tests/:questionId/check`
- `POST /api/v1/student/practice/generate`（AI 生成 + 模板兜底）
- `POST /api/v1/student/practice/submit`（幂等提交：`taskId + studentId`）
- `POST /api/v1/student/nodes/:nodeId/explain`（专项讲解）

---

## 三、AI 统一接口

- `POST /api/v1/ai/parse-knowledge`
- `GET /api/v1/ai/coursewares/:courseId/knowledge-graph`
- `POST /api/v1/ai/coursewares/:courseId/ask`

---

## 四、旧版兼容接口

以下接口仍保留，供旧前端或联调脚本继续使用：

### 教师端

- `GET /api/teacher/courseware-list`
- `POST /api/teacher/upload-courseware`
- `DELETE /api/teacher/courseware/:courseId`
- `POST /api/teacher/publish-courseware`
- `GET /api/teacher/script/:courseId/:page`
- `POST /api/teacher/script/save`
- `POST /api/teacher/ai-generate-script`
- `GET /api/teacher/student-stats/:courseId`
- `GET /api/teacher/question-records/:courseId`
- `GET /api/teacher/card-data/:courseId`

### 学生端

- `GET /api/student/courseware-list`
- `POST /api/student/session/start`
- `POST /api/student/progress/update`
- `GET /api/student/script/:courseId/:page`
- `POST /api/student/qa/stream`
- `POST /api/student/courseware/page`
- `POST /api/student/ai/question`
- `POST /api/student/ai/traceQuestion`
- `GET /api/student/studyData`
- `GET /api/student/breakpoint`
- `PUT /api/student/breakpoint`
- `POST /api/student/saveNote`

### 薄弱点与知识点

- `GET /api/weakPoint/getList`
- `POST /api/weakPoint/getExplain`
- `POST /api/weakPoint/getTest`
- `POST /api/weakPoint/checkAnswer`
- `POST /api/ai/parseKnowledge`

---

## 五、公开接口

- `GET /api/courseware/:courseId/page/:pageNum`
- `GET /api/v1/courseware/:courseId/page/:pageNum`

---

## 六、平台 V1 接口（标准 REST，无需签名）

以下接口位于 `/api/v1/platform`，返回统一格式 `{ "code": 200, "message": "请求成功", "data": ... }`，供内部管理端调用。

### 1. 平台总览
- `GET /api/v1/platform/overview`

### 2. 用户管理
- `GET /api/v1/platform/users`
- `GET /api/v1/platform/users/:userId`
- `POST /api/v1/platform/syncUser`

### 3. 课程管理
- `GET /api/v1/platform/courses`
- `POST /api/v1/platform/courses`
- `GET /api/v1/platform/courses/:courseId`
- `PUT /api/v1/platform/courses/:courseId`
- `DELETE /api/v1/platform/courses/:courseId`
- `POST /api/v1/platform/syncCourse`

### 4. 班级管理
- `GET /api/v1/platform/classes`
- `POST /api/v1/platform/classes`
- `GET /api/v1/platform/classes/:classId`
- `PUT /api/v1/platform/classes/:classId`
- `DELETE /api/v1/platform/classes/:classId`

### 5. 选课管理
- `GET /api/v1/platform/enrollments`
- `POST /api/v1/platform/enrollments`
- `GET /api/v1/platform/enrollments/:enrollmentId`
- `PUT /api/v1/platform/enrollments/:enrollmentId`
- `DELETE /api/v1/platform/enrollments/:enrollmentId`

---

## 七、开放平台接口（需签名，泛雅等外部对接）

以下接口位于 `/api/v1/open/platform`，需携带 OpenAPI 签名（如 `enc` + `time`）校验：

- `GET /api/v1/open/platform/overview`
- `GET /api/v1/open/platform/users`
- `GET /api/v1/open/platform/users/:userId`
- `POST /api/v1/open/platform/syncUser`
- `GET /api/v1/open/platform/courses`
- `POST /api/v1/open/platform/courses`
- `GET /api/v1/open/platform/courses/:courseId`
- `PUT /api/v1/open/platform/courses/:courseId`
- `DELETE /api/v1/open/platform/courses/:courseId`
- `POST /api/v1/open/platform/syncCourse`
- `GET /api/v1/open/platform/classes`
- `POST /api/v1/open/platform/classes`
- `GET /api/v1/open/platform/classes/:classId`
- `PUT /api/v1/open/platform/classes/:classId`
- `DELETE /api/v1/open/platform/classes/:classId`
- `GET /api/v1/open/platform/enrollments`
- `POST /api/v1/open/platform/enrollments`
- `GET /api/v1/open/platform/enrollments/:enrollmentId`
- `PUT /api/v1/open/platform/enrollments/:enrollmentId`
- `DELETE /api/v1/open/platform/enrollments/:enrollmentId`

其他开放接口（需签名）：

- `POST /api/v1/lesson/parse`
- `POST /api/v1/lesson/generateScript`
- `POST /api/v1/lesson/generateAudio`
- `POST /api/v1/qa/interact`
- `POST /api/v1/qa/voiceToText`
- `POST /api/v1/progress/track`
- `POST /api/v1/progress/adjust`

---

## 八、说明

1. 当前实现优先保证前后端统一接口可用。
2. 平台标准接口（六）无需签名，直接请求 `/api/v1/platform/*`；对外/泛雅对接使用（七）`/api/v1/open/platform/*` 并携带签名。
3. 旧版接口 continue 保留，便于历史页面与 mock 联调。
4. AI 相关能力在真实 AI 引擎不可用时，部分接口会回退到 mock 结果，避免前端阻塞。
5. `student/practice/generate` 当前为“AI 出题优先，失败自动回退模板题”，以保证接口稳定可用。
6. `student/favorites/:favoriteId` 删除接口已增加用户归属校验，防止越权删除。

