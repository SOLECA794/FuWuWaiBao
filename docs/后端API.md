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

## 六、开放平台接口

以下接口位于 `/api/v1` 下，并带签名校验中间件：

- `/lesson/parse`
- `/lesson/generateScript`
- `/lesson/generateAudio`
- `/qa/interact`
- `/qa/voiceToText`
- `/progress/track`
- `/progress/adjust`
- `/platform/users`
- `/platform/users/{userId}`
- `/platform/courses`
- `/platform/courses/{courseId}`
- `/platform/classes`
- `/platform/classes/{classId}`
- `/platform/enrollments`
- `/platform/enrollments/{enrollmentId}`
- `/platform/courses` `POST/PUT/DELETE`
- `/platform/classes` `POST/PUT/DELETE`
- `/platform/enrollments` `POST/PUT/DELETE`
- `/platform/courses/{courseId}`
- `/platform/syncCourse`
- `/platform/syncUser`
- `/platform/overview`

---

## 七、说明

1. 当前实现优先保证前后端统一接口可用。
2. 平台接口除了 `syncCourse`、`syncUser` 外，现已增加 `users/courses/classes/enrollments` 列表、详情和基础运营写接口，以及 `overview` 用于查看用户、课程、班级、选课关系的落库总览。
3. 旧版接口继续保留，便于历史页面与 mock 联调。
3. AI 相关能力在真实 AI 引擎不可用时，部分接口会回退到 mock 结果，避免前端阻塞。
