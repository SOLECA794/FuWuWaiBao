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
- `GET /api/v1/student/coursewares/:courseId/stats`

### 3. 薄弱点与练习

- `GET /api/v1/student/coursewares/:courseId/weak-points`
- `GET /api/v1/student/weak-points/:weakPointId/explain`
- `POST /api/v1/student/weak-points/:weakPointId/generate-test`
- `POST /api/v1/student/tests/:questionId/check`

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

以下接口位于 `/api/v1/platform`，返回统一格式 `{ "code": 200, "message": "请求成功", "data": ... }`，供内部或管理端调用。

- `/lesson/parse`
- `/lesson/generateScript`
- `/lesson/generateAudio`
- `/qa/interact`
- `/qa/voiceToText`
- `/progress/track`
- `/progress/adjust`
- `/platform/syncCourse`
- `/platform/syncUser`
=======
### 1. 平台总览

- `GET /api/v1/platform/overview` — 返回用户、课程、班级、选课统计与最近同步记录

### 2. 用户管理

- `GET /api/v1/platform/users` — 用户列表，支持 `role` / `orgCode` / `keyword` / `page` / `pageSize`
- `GET /api/v1/platform/users/:userId` — 用户详情（含关联课程/班级/选课）
- `POST /api/v1/platform/syncUser` — 用户同步（Body: JSON）

### 3. 课程管理

- `GET /api/v1/platform/courses` — 课程列表，支持 `teacherId` / `status` / `orgCode` / `keyword` / `page` / `pageSize`
- `POST /api/v1/platform/courses` — 创建课程（Body: JSON）
- `GET /api/v1/platform/courses/:courseId` — 课程详情
- `PUT /api/v1/platform/courses/:courseId` — 更新课程
- `DELETE /api/v1/platform/courses/:courseId` — 删除课程（级联班级与选课）
- `POST /api/v1/platform/syncCourse` — 课程同步（Body: JSON）

### 4. 班级管理

- `GET /api/v1/platform/classes` — 班级列表，支持 `courseId` / `teacherId` / `status` / `keyword` / `page` / `pageSize`
- `POST /api/v1/platform/classes` — 创建班级（Body: JSON）
- `GET /api/v1/platform/classes/:classId` — 班级详情
- `PUT /api/v1/platform/classes/:classId` — 更新班级
- `DELETE /api/v1/platform/classes/:classId` — 删除班级（级联选课）

### 5. 选课管理

- `GET /api/v1/platform/enrollments` — 选课列表，支持 `courseId` / `classId` / `userId` / `role` / `status` / `page` / `pageSize`
- `POST /api/v1/platform/enrollments` — 创建选课（Body: JSON）
- `GET /api/v1/platform/enrollments/:enrollmentId` — 选课详情
- `PUT /api/v1/platform/enrollments/:enrollmentId` — 更新选课
- `DELETE /api/v1/platform/enrollments/:enrollmentId` — 删除选课
>>>>>>> Stashed changes

---

## 七、开放平台接口（需签名，泛雅等外部对接）

以下接口位于 `/api/v1/open/platform`，需携带 OpenAPI 签名（如 `enc` + `time`）校验：

- `GET /api/v1/open/platform/overview`
- `GET /api/v1/open/platform/users`、`GET /api/v1/open/platform/users/:userId`
- `GET/POST/PUT/DELETE /api/v1/open/platform/courses`、`/api/v1/open/platform/courses/:courseId`
- `GET/POST/PUT/DELETE /api/v1/open/platform/classes`、`/api/v1/open/platform/classes/:classId`
- `GET/POST/PUT/DELETE /api/v1/open/platform/enrollments`、`/api/v1/open/platform/enrollments/:enrollmentId`
- `POST /api/v1/open/platform/syncCourse`、`POST /api/v1/open/platform/syncUser`

其他开放接口（需签名）：

- `POST /api/v1/lesson/parse`、`/api/v1/lesson/generateScript`、`/api/v1/lesson/generateAudio`
- `POST /api/v1/qa/interact`、`/api/v1/qa/voiceToText`
- `POST /api/v1/progress/track`、`/api/v1/progress/adjust`

---

## 八、说明

1. 当前实现优先保证前后端统一接口可用。
2. 旧版接口继续保留，便于历史页面与 mock 联调。
3. AI 相关能力在真实 AI 引擎不可用时，部分接口会回退到 mock 结果，避免前端阻塞。
=======
2. 平台标准接口（六）无需签名，直接请求 `/api/v1/platform/*`；对外/泛雅对接使用（七）`/api/v1/open/platform/*` 并携带签名。
3. 旧版接口继续保留，便于历史页面与 mock 联调。
4. AI 相关能力在真实 AI 引擎不可用时，部分接口会回退到 mock 结果，避免前端阻塞。
>>>>>>> Stashed changes
