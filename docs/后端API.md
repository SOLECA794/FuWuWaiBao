# 智能互动讲课系统 - 接口文档 v2.0（教师端 + 学生端）

## 基础信息
- **基础URL**: `http://localhost:8080/api/v1`
- **响应格式**: 所有接口返回 JSON 格式
- **字符编码**: UTF-8
- **身份认证**: 需要鉴权的接口在 HTTP Header 中携带 `Authorization: Bearer <Token>`

---

## 目录
1. [课件管理（教师端）](#1-课件管理教师端)
2. [讲稿编辑（教师端）](#2-讲稿编辑教师端)
3. [学情分析（教师端）](#3-学情分析教师端)
4. [提问记录（教师端）](#4-提问记录教师端)
5. [学习卡点分析（教师端）](#5-学习卡点分析教师端)
6. [AI学伴与互动答疑（学生端）](#6-ai学伴与互动答疑学生端)
7. [薄弱点诊断与练习（学生端）](#7-薄弱点诊断与练习学生端)
8. [学习过程数据（学生端）](#8-学习过程数据学生端)
9. [课件预览（公开）](#9-课件预览公开)

---

## 1. 课件管理（教师端）
> 基础路径: `/teacher/coursewares`

### 1.1 获取课件列表
- **接口**: `GET /`
- **功能**: 获取所有已上传的课件
- **请求参数**: 无

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": [
    {
      "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
      "title": "测试课件",
      "fileType": "pdf",
      "status": "published",
      "createdAt": "2026-03-02T19:29:06+08:00"
    }
  ]
}
```

### 1.2 上传课件
- **接口**: `POST /upload`
- **功能**: 上传并解析课件（支持 PDF/PPT）
- **Content-Type**: `multipart/form-data`

**请求参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | File | 是 | 课件文件（PDF/PPT/PPTX） |
| title | String | 否 | 课件标题（默认"未命名课件"） |

**响应示例**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "id": "课件ID",
    "title": "课件标题",
    "file_url": "文件访问URL",
    "file_type": "pdf",
    "created_at": "2026-03-03T10:00:00+08:00"
  }
}
```

### 1.3 删除课件
- **接口**: `DELETE /{courseId}`
- **功能**: 删除指定课件

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功"
}
```

### 1.4 发布课件
- **接口**: `POST /{courseId}/publish`
- **功能**: 发布课件给学生端

**请求体**: 无

**响应示例**:
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
    "publishedAt": "2026-03-05T10:30:00+08:00"
  }
}
```

---

## 2. 讲稿编辑（教师端）
> 基础路径: `/teacher/coursewares/{courseId}/scripts`

### 2.1 获取讲稿
- **接口**: `GET /{pageNum}`
- **功能**: 获取指定课件指定页码的讲稿

**路径参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| courseId | String | 课件ID |
| pageNum | Integer | 页码 |

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
    "pageNum": 1,
    "content": "这是第一页的讲稿内容"
  }
}
```

### 2.2 保存讲稿
- **接口**: `PUT /{pageNum}`
- **功能**: 保存或更新讲稿

**请求体**:
```json
{
  "content": "这是更新后的讲稿内容"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "保存成功"
}
```

### 2.3 AI生成讲稿
- **接口**: `POST /{courseId}/scripts/ai-generate`
- **功能**: AI自动生成讲稿内容

**请求体**:
```json
{
  "pageNum": 1,
  "mode": "llm"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
    "pageNum": 1,
    "content": "## Go语言基础教程 第1页：课程导入\n\n### 教学目标\n- 了解本章节的核心概念\n..."
  }
}
```

---

## 3. 学情分析（教师端）
> 基础路径: `/teacher/coursewares/{courseId}`

### 3.1 获取班级宏观学情
- **接口**: `GET /stats`
- **功能**: 获取课件的学情分析数据（停留时长、提问频次、词云）

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "pageStayTime": [
      {"page": 1, "stayAvg": 120.5},
      {"page": 2, "stayAvg": 85.3}
    ],
    "questionFreq": [
      {"page": 1, "count": 3},
      {"page": 2, "count": 2}
    ],
    "wordCloud": [
      {"word": "依赖注入", "count": 3},
      {"word": "微服务", "count": 2}
    ]
  }
}
```

---

## 4. 提问记录（教师端）
> 基础路径: `/teacher/coursewares/{courseId}`

### 4.1 获取提问记录
- **接口**: `GET /questions`
- **功能**: 分页获取学生的提问记录

**查询参数**:
| 参数名 | 类型 | 必填 | 默认 | 说明 |
|--------|------|------|------|------|
| page | Integer | 否 | 1 | 页码 |
| pageSize | Integer | 否 | 20 | 每页条数（最大100） |

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "list": [
      {
        "id": "记录ID",
        "user_id": "user001",
        "course_id": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
        "page_index": 1,
        "question": "什么是依赖注入？",
        "answer": "依赖注入是一种设计模式...",
        "created_at": "2026-03-03T10:00:00+08:00"
      }
    ],
    "total": 7,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## 5. 学习卡点分析（教师端）

### 5.1 获取学习卡点数据
- **接口**: `GET /teacher/card-data/{courseId}`
- **功能**: 获取指定课件的学习卡点分析数据（提问量、停留时长、卡点指数）

**路径参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| courseId | String | 课件ID |

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "pageStats": [
      {"page": 1, "questionCount": 2, "stayTime": 20, "cardIndex": 2.0},
      {"page": 3, "questionCount": 4, "stayTime": 32, "cardIndex": 3.6}
    ],
    "topPages": [
      {"page": 5, "value": 8, "ratio": 53.33}
    ],
    "totalQuestions": 8
  }
}
```

---

## 6. AI学伴与互动答疑（学生端）
> 基础路径: `/api/v1/ai/coursewares`

### 6.1 获取课件知识图谱
- **接口**: `GET /{courseId}/knowledge-graph`
- **功能**: 获取 AI 解析拆解出的知识点层级结构

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "courseId": "PY202501",
    "structure": [
      {
        "chapter": "第一章：数据清洗基础",
        "knowledgePoints": [
          {"name": "缺失值处理", "subPoints": ["fillna()", "interpolate()", "dropna()"]}
        ]
      }
    ]
  }
}
```

### 6.2 智能多模态答疑
- **接口**: `POST /{courseId}/ask`
- **功能**: 在学习某页时向 AI 提问（包含常规提问和圈选溯源提问）

**请求体**:
```json
{
  "pageNum": 3,
  "type": "text",
  "question": "这里提到的 fillna 怎么用？",
  "tracePoint": {
    "x": 200,
    "y": 150
  }
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "answer": "fillna() 是用于填充缺失值的函数..."
  }
}
```

---

## 7. 薄弱点诊断与练习（学生端）
> 基础路径: `/api/v1/student`

### 7.1 获取个人薄弱点列表
- **接口**: `GET /coursewares/{courseId}/weak-points`
- **功能**: 获取学生在指定课程中的薄弱点列表

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| studentId | String | 是 | 学生ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": [
    {
      "id": "wp_001",
      "name": "缺失值填充",
      "description": "数据中空值的处理方法",
      "count": 5,
      "mastery": 60
    }
  ]
}
```

### 7.2 薄弱点 AI 详细讲解
- **接口**: `GET /weak-points/{weakPointId}/explain`
- **功能**: AI 对该薄弱知识点进行重构和二次讲解

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "id": "wp_001",
    "name": "缺失值填充",
    "title": "缺失值填充 · 知识点讲解",
    "content": "缺失值是数据中为空的部分，常用处理方法：...",
    "examples": ["df.fillna(0)", "df.interpolate()"],
    "relatedConcepts": ["数据清洗", "异常值处理"]
  }
}
```

### 7.3 生成随堂检测题
- **接口**: `POST /weak-points/{weakPointId}/generate-test`
- **功能**: 根据薄弱点生成对应的检测习题

**请求体**:
```json
{
  "questionType": "single"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "questionId": "q_001",
    "type": "single",
    "content": "处理缺失值时，以下哪种方法最适合时间序列数据？",
    "options": [
      {"key": "A", "value": "fillna(0)"},
      {"key": "B", "value": "interpolate()"}
    ]
  }
}
```

### 7.4 提交并校验答案
- **接口**: `POST /tests/{questionId}/check`
- **功能**: 校验习题答案，返回解析和掌握状态

**请求体**:
```json
{
  "userAnswer": "B"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "isCorrect": true,
    "correctAnswer": "B",
    "explanation": "interpolate() 线性插值适合时间序列数据",
    "masteryDelta": 10,
    "newMastery": 75
  }
}
```

---

## 8. 学习过程数据（学生端）
> 基础路径: `/api/v1/student`

### 8.1 开始学习会话
- **接口**: `POST /session/start`
- **功能**: 初始化学生在某课件上的学习会话

**请求体**:
```json
{
  "userId": "user001",
  "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "sessionId": "sess_xxx",
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d"
  }
}
```

### 8.2 上报播放进度
- **接口**: `POST /progress/update`
- **功能**: 上报当前播放游标（页码 + node）

**请求体**:
```json
{
  "sessionId": "sess_xxx",
  "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
  "page": 5,
  "nodeId": "p5_n12"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "ok"
}
```

### 8.3 获取学习断点
- **接口**: `GET /coursewares/{courseId}/breakpoint`
- **功能**: 获取学生上次学习的页码（用于续播）

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | String | 是 | 学生ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "lastPageNum": 5
  }
}
```

### 8.4 更新学习断点
- **接口**: `PUT /coursewares/{courseId}/breakpoint`
- **功能**: 更新学生学习断点

**请求体**:
```json
{
  "userId": "user001",
  "pageNum": 5
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "ok"
}
```

### 8.5 保存随堂笔记
- **接口**: `POST /coursewares/{courseId}/notes`
- **功能**: 保存学生笔记

**请求体**:
```json
{
  "userId": "user001",
  "pageNum": 3,
  "content": "fillna() 用于填充缺失值",
  "x": 0,
  "y": 0
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "status": "saved"
  }
}
```

### 8.6 获取某页讲稿（学生播放用）
- **接口**: `GET /coursewares/{courseId}/pages/{pageNum}`
- **功能**: 获取指定课件指定页的结构化脚本

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
    "page": 5,
    "nodes": [
      {"node_id": "p5_n1", "type": "opening", "text": "..."}
    ],
    "page_summary": "..."
  }
}
```

### 8.7 获取个人微观学情
- **接口**: `GET /coursewares/{courseId}/stats`
- **功能**: 返回学生本人的学习数据

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| userId | String | 是 | 学生ID |

**响应示例**:
```json
{
  "code": 200,
  "message": "请求成功",
  "data": {
    "totalQuestions": 12,
    "studyHours": 2.5,
    "mastery": {
      "章节一": 85,
      "章节二": 70
    }
  }
}
```

### 8.8 问答流式接口
- **接口**: `POST /qa/stream`
- **功能**: 打断提问，SSE 流式返回回答
- **Content-Type**: `text/event-stream; charset=utf-8`

**请求体**:
```json
{
  "sessionId": "sess_xxx",
  "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
  "page": 5,
  "nodeId": "p5_n12",
  "question": "这个公式里的 x 表示什么？"
}
```

**SSE 输出示例**:
```text
event: token
data: {"text":"x"}

event: token
data: {"text":"表示"}

event: sentence
data: {"text":"这里的 x 通常表示..."}

event: final
data: {"need_reteach":false,"source_page":5,"resume_page":5,"resume_node_id":"p5_n13"}
```

---

## 9. 课件预览（公开）

### 9.1 获取预览图片
- **接口**: `GET /courseware/{courseId}/page/{pageNum}`
- **功能**: 获取课件指定页码的预览图片

**路径参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| courseId | String | 课件ID |
| pageNum | Integer | 页码 |

**响应**:
- 成功: 302 重定向到图片URL
- 失败:
```json
{
  "code": 404,
  "message": "预览图不存在"
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 无权访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 注意事项

1. **文件上传**仅支持 PDF、PPT、PPTX 格式
2. **分页接口**默认每页20条，最大支持100条
3. **所有响应**均包含 `code` 字段，`code=200` 表示成功
4. **中文编码**已统一为 UTF-8
5. **时间格式**为 ISO8601: `2006-01-02T15:04:05+08:00`
6. **SSE 流式接口**使用 `text/event-stream` 格式，每帧以空行结尾

---

## 更新记录

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2026-03-05 | v2.0 | 统一基础URL为 `/api/v1`，按规范重构所有接口路径，添加薄弱点诊断、学习会话、SSE问答等新接口 |