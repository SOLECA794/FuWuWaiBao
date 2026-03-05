# 智能学习平台 API 接口文档 (v2.0 统一重构版)

## 一、全局规范

1. **基础 URL**: `/api/v1`
2. **数据交互格式**: `application/json; charset=utf-8` （除文件上传外）
3. **主键类型**: 所有资源 ID（`courseId`, `studentId`, `questionId` 等）统一使用 **字符串 (String)** 类型，建议使用 UUID。
4. **统一响应结构**:
   ```json
   {
     "code": 200,             // 200: 成功, 400: 参数错误, 401: 未授权, 404: 不存在, 500: 服务器错误
     "message": "请求成功",   // 提示信息
     "data": {}               // 具体响应数据，可能为 Object 或 Array
   }
   ```
5. **身份认证**: 所有需要鉴权的接口，均在 HTTP Header 中携带 `Authorization: Bearer <Token>`。

---

## 二、课件与讲稿管理 (教师端)
> 基础路径: `/api/v1/teacher/coursewares`

### 2.1 获取课件列表
- **接口**: `GET /`
- **功能**: 获取当前教师的所有课件。
- **响应**: `data` 返回课件数组（包含 `courseId`, `title`, `fileType`, `status`, `createdAt`）。

### 2.2 上传并解析课件
- **接口**: `POST /upload`
- **Content-Type**: `multipart/form-data`
- **参数**: 
  - `file` (File): PDF 或 PPTX 文件
  - `title` (String, 可选): 课件标题
- **说明**: 异步/耗时接口，上传后底层自动触发 AI 解析与切片。返回 `courseId`。

### 2.3 获取页面讲稿
- **接口**: `GET /{courseId}/scripts/{pageNum}`
- **功能**: 获取某页的具体讲稿内容。
- **响应**: `data: { "courseId": "...", "pageNum": 1, "content": "..." }`

### 2.4 更新页面讲稿
- **接口**: `PUT /{courseId}/scripts/{pageNum}`
- **体参数**: `{ "content": "修改后的讲稿内容" }`

### 2.5 AI 生成讲稿
- **接口**: `POST /{courseId}/scripts/ai-generate`
- **体参数**: `{ "pageNum": 1, "mode": "llm" }`
- **功能**: 调用大模型根据当前页内容重新生成讲稿文字。

### 2.6 发布课件
- **接口**: `POST /{courseId}/publish`
- **功能**: 将课件状态改为已发布，向学生开放。

---

## 三、AI 学伴与互动答疑 (学生端)
> 基础路径: `/api/v1/ai/coursewares`

### 3.1 获取课件知识图谱 (解析结果)
- **接口**: `GET /{courseId}/knowledge-graph`
- **功能**: 获取 AI 解析拆解出的知识点层级结构（章节 $\rightarrow$ 知识点 $\rightarrow$ 子知识点）。

### 3.2 智能多模态答疑
- **接口**: `POST /{courseId}/ask`
- **功能**: 在学习某页时向 AI 提问（包含常规提问和圈选溯源提问）。
- **体参数**: 
  ```json
  {
    "pageNum": 3,
    "type": "text",              // text / audio / image
    "question": "这里提到的 fillna 怎么用？",
    "tracePoint": {              // 可选：如果不为空，则是圈选定位溯源提问
      "x": 200, 
      "y": 150 
    }
  }
  ```
- **响应**: 返回 AI 智能解答内容。

---

## 四、薄弱点诊断与练习 (学生端)
> 基础路径: `/api/v1/student`

### 4.1 获取个人薄弱点列表
- **接口**: `GET /coursewares/{courseId}/weak-points`
- **响应**: 返回知识树中被诊断为薄弱的节点列表及掌握度评估。

### 4.2 薄弱点 AI 详细讲解
- **接口**: `GET /weak-points/{weakPointId}/explain`
- **功能**: AI 对该薄弱知识点进行重构和二次讲解（Reteach）。

### 4.3 生成随堂检测题
- **接口**: `POST /weak-points/{weakPointId}/generate-test`
- **体参数**: `{ "questionType": "single" }` // single(单选), multiple(多选)
- **响应**: 返回 AI 动态生成的题目、选项及 `questionId`。

### 4.4 提交并校验答案
- **接口**: `POST /tests/{questionId}/check`
- **体参数**: `{ "userAnswer": "A" }`
- **响应**: 返回对错判断、详细解析以及更新后的知识点掌握状态。

---

## 五、学习过程数据 (学生端)
> 基础路径: `/api/v1/student/coursewares`

### 5.1 获取/更新学习断点
- **接口**: `GET /{courseId}/breakpoint` (获取) / `PUT /{courseId}/breakpoint` (更新)
- **功能**: 记录学生上次学到了哪一页，用于续播。
- **体参数** (PUT时): `{ "pageNum": 5 }`

### 5.2 保存随堂笔记
- **接口**: `POST /{courseId}/notes`
- **体参数**: `{ "pageNum": 3, "content": "我的笔记内容...", "x": 0, "y": 0 }`

---

## 六、学情分析数据展示

### 6.1 教师端 - 班级宏观学情
- **接口**: `GET /api/v1/teacher/coursewares/{courseId}/stats`
- **响应**: 包含各页面的停留时长统计、提问频次图表、高频提问词云（如：“依赖注入”、“AOP”等）。

### 6.2 教师端 - 历史提问记录
- **接口**: `GET /api/v1/teacher/coursewares/{courseId}/questions`
- **查询参数**: `?page=1&pageSize=20`
- **功能**: 教师可翻阅学生对该课件的所有 AI 提问历史与解答记录。

### 6.3 学生端 - 个人微观学情
- **接口**: `GET /api/v1/student/coursewares/{courseId}/stats`
- **功能**: 返回学生本人的学习总时长、提问次数、笔记数及雷达图预估掌握情况。

---

## 七、后续 API 设计与扩展规范

为保证项目的长期可维护性与多端一致性，后续新增接口需严格遵守以下设计要求：

### 1. 命名与路由规范
- **动词与路径**: 资源使用复数名词（如 `/coursewares`），其后的操作使用标准 HTTP 动词（GET 获取, POST 创建, PUT 全量更新, PATCH 局部更新, DELETE 删除）。针对非 CRUD 的行为，在路径末尾使用动词描述（如 `/publish`, `/generate-test`）。
- **连字符命名**: URL 路径中的多单词组合统一使用小写连字符 `kebab-case`（如 `/weak-points`），严禁使用下划线或驼峰。
- **字段命名**: JSON 请求与响应体中的字段统一使用 `lowerCamelCase`（小写驼峰，如 `courseId`, `pageNum`）。

### 2. 参数处理要求
- **路径参数 (Path)**: 用于唯一定位资源（如 `/{courseId}`）。
- **查询参数 (Query)**: 用于过滤、排序、分页（如 `?page=1&pageSize=20&status=published`）。
- **请求体 (Body)**: 复杂的业务数据必须放在 Body 中，且必须校验必填项，空值应返回 400 状态码及具体错误字段描述。

### 3. AI 接口专用要求
- **超时处理**: 涉及 LLM 或文件解析的接口（耗时超过 5s），必须在响应头或 Body 中包含 `processingTime`，并建议在前端实现 Loading 状态管理。
- **流式响应 (Optional)**: 针对长文本生成的问答接口，后续应优先考虑支持 `SSE (Server-Sent Events)` 流式返回，提升用户体验。
- **上下文关联**: 提问类接口必须携带 `context` 或 `sourceId`，严禁脱离课件背景的孤立问答。

### 4. 健壮性与安全性
- **状态码准则**: 
  - `200/201`: 操作成功。
  - `400`: 参数校验失败（前端传错）。
  - `401`: 未授权（Token 失效）。
  - `403`: 无权访问某课件（权限控制）。
  - `500`: 后端逻辑或 AI 引擎宕机。
- **幂等性**: GET、PUT、DELETE 接口必须保证幂等，多次调用不应产生副作用。
- **版本控制**: 如果接口逻辑发生重大破坏性变更，必须升级路径版本号（如从 `/v1/` 到 `/v2/`）。

### 5. 文档同步要求
- **实时更新**: 任何后端（Go/Python）接口代码的变更，必须同步更新至此 `API_DESIGN_V2.md`，确保“代码即文档”。
- **示例覆盖**: 新增复杂接口必须附带 Request Body 示例和 Response Data 示例。
