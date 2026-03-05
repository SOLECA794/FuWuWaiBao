# 智能互动讲课系统 - 教师端接口文档

## 基础信息
- **基础URL**: `http://localhost:8080/api`
- **响应格式**: 所有接口返回 JSON 格式
- **字符编码**: UTF-8

---

## 目录
1. [课件管理](#1-课件管理)
2. [讲稿编辑](#2-讲稿编辑)
3. [学情分析](#3-学情分析)
4. [提问记录](#4-提问记录)
5. [课件预览](#5-课件预览)

---

## 1. 课件管理

### 1.1 获取课件列表
- **接口**: `GET /teacher/courseware-list`
- **功能**: 获取所有已上传的课件
- **请求参数**: 无

**响应示例**:
```json
{
  "code": 200,
  "data": [
    {
      "id": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
      "title": "测试课件",
      "file_url": "http://localhost:9000/courses/xxx/test.pdf",
      "file_type": "pdf",
      "total_page": 0,
      "created_at": "2026-03-02T19:29:06+08:00"
    }
  ]
}
```

### 1.2 上传课件
- **接口**: `POST /teacher/upload-courseware`
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
- **接口**: `DELETE /teacher/courseware/{courseId}`
- **功能**: 删除指定课件

**响应示例**:
```json
{
  "code": 200,
  "message": "删除成功"
}
```

### 1.4 发布课件
- **接口**: `POST /teacher/publish-courseware`
- **功能**: 发布课件给学生端

**请求体**:
```json
{
  "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
  "scope": "all"
}
```

**响应示例**:
```json
{
  "code": 200,
  "message": "发布成功",
  "data": {
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
    "scope": "all",
    "publishedAt": "2026-03-03 10:30:00"
  }
}
```

---

## 2. 讲稿编辑

### 2.1 获取讲稿
- **接口**: `GET /teacher/script/{courseId}/{page}`
- **功能**: 获取指定课件指定页码的讲稿

**路径参数**:
| 参数名 | 类型 | 说明 |
|--------|------|------|
| courseId | String | 课件ID |
| page | Integer | 页码 |

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
    "page": 1,
    "content": "这是第一页的讲稿内容"
  }
}
```

### 2.2 保存讲稿
- **接口**: `POST /teacher/script/save`
- **功能**: 保存或更新讲稿

**请求体**:
```json
{
  "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
  "page": 1,
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
- **接口**: `POST /teacher/ai-generate-script`
- **功能**: AI自动生成讲稿内容

**请求体**:
```json
{
  "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
  "page": 1,
  "courseName": "Go语言基础教程"
}
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "courseId": "8abc34a7-4d05-41c5-b3b9-7b629463444d",
    "page": 1,
    "content": "## Go语言基础教程 第1页：课程导入\n\n### 教学目标\n- 了解本章节的核心概念\n..."
  }
}
```

---

## 3. 学情分析

### 3.1 获取学情数据
- **接口**: `GET /teacher/student-stats/{courseId}`
- **功能**: 获取课件的学情分析数据

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "pageStats": [
      {"page": 1, "count": 3},
      {"page": 2, "count": 2},
      {"page": 3, "count": 2}
    ],
    "keywords": [
      {"word": "依赖注入", "count": 3},
      {"word": "微服务", "count": 2},
      {"word": "AOP", "count": 1}
    ],
    "totalQuestions": 7,
    "activeUsers": 5
  }
}
```

---

## 4. 提问记录

### 4.1 获取提问记录
- **接口**: `GET /teacher/question-records/{courseId}`
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
    "pageSize": 20,
    "totalPage": 1
  }
}
```

---

## 5. 课件预览

### 5.1 获取预览图片
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

## 6. 学习卡点分析

### 6.1 获取学习卡点数据
- **接口**: `GET /api/teacher/card-data/{courseId}`
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
      {"page": 3, "questionCount": 4, "stayTime": 32, "cardIndex": 3.6},
      {"page": 5, "questionCount": 8, "stayTime": 45, "cardIndex": 6.25}
    ],
    "topPages": [
      {"page": 5, "value": 8, "ratio": 53.33},
      {"page": 3, "value": 4, "ratio": 26.67},
      {"page": 1, "value": 2, "ratio": 13.33}
    ],
    "totalQuestions": 8
  }
}
```

| 字段 | 说明 |
|------|------|
| pageStats.page | 页码 |
| pageStats.questionCount | 提问量 |
| pageStats.stayTime | 平均停留时长(秒) |
| pageStats.cardIndex | 卡点指数(0-10) |
| topPages | TOP5卡点页面 |
| totalQuestions | 总提问数 |

---

## 7. AI薄弱点功能

### 7.1 知识点解析
- **接口**: `POST /api/ai/parseKnowledge`
- **功能**: 解析课件内容，拆解知识点层级结构

**请求体**:
```json
{
  "fileContent": "数据清洗是数据分析的重要步骤...",
  "fileType": "pptx",
  "studentId": "2025001"
}
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "structure": [
      {
        "chapter": "第一章：数据清洗基础",
        "knowledgePoints": [
          {"name": "缺失值处理", "subPoints": ["fillna()", "interpolate()", "dropna()"]},
          {"name": "异常值识别", "subPoints": ["Z-Score", "IQR"]}
        ]
      }
    ]
  }
}
```

### 7.2 获取薄弱点列表
- **接口**: `GET /api/weakPoint/getList`
- **功能**: 获取指定学生的AI诊断薄弱点列表

**查询参数**:
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| studentId | String | 是 | 学生ID |
| courseId | String | 是 | 课程ID |

**响应示例**:
```json
{
  "code": 200,
  "data": [
    {"name": "缺失值填充", "count": 5, "mastery": 60},
    {"name": "异常值识别", "count": 3, "mastery": 45},
    {"name": "重复值处理", "count": 2, "mastery": 80}
  ]
}
```

| 字段 | 说明 |
|------|------|
| name | 薄弱点名称 |
| count | 出现次数 |
| mastery | 掌握程度(0-100) |

### 7.3 获取薄弱点讲解
- **接口**: `POST /api/weakPoint/getExplain`
- **功能**: 获取指定薄弱点的详细讲解内容

**请求体**:
```json
{
  "weakPointName": "缺失值填充",
  "studentId": "2025001"
}
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "title": "缺失值填充 · 知识点讲解",
    "content": "缺失值是数据中为空的部分，常用方法：\n1. fillna() 填充常数、均值、中位数\n2. interpolate() 线性插值（适合时序）\n3. dropna() 直接删除行/列",
    "examples": [
      "df.fillna(0)  # 用0填充",
      "df.interpolate()  # 线性插值",
      "df.dropna()  # 删除缺失值"
    ]
  }
}
```

### 7.4 生成检测习题
- **接口**: `POST /api/weakPoint/getTest`
- **功能**: 根据薄弱点生成对应的检测习题

**请求体**:
```json
{
  "weakPointName": "缺失值填充",
  "studentId": "2025001",
  "questionType": "single"
}
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "questionId": "Q001",
    "content": "处理缺失值时，以下哪种方法最适合时间序列数据？",
    "type": "single",
    "options": [
      "A. fillna(0) 用0填充",
      "B. interpolate() 线性插值",
      "C. dropna() 删除缺失值",
      "D. fillna(method='ffill') 向前填充"
    ]
  }
}
```

### 7.5 校验答案
- **接口**: `POST /api/weakPoint/checkAnswer`
- **功能**: 校验习题答案，返回解析和掌握状态

**请求体**:
```json
{
  "studentId": "2025001",
  "questionId": "Q001",
  "userAnswer": "B"
}
```

**响应示例**:
```json
{
  "code": 200,
  "data": {
    "isCorrect": true,
    "correctAnswer": "B",
    "explanation": "interpolate() 线性插值适合时间序列数据，可以基于前后值推算缺失值。",
    "masteryDelta": 10
  }
}
```

| 字段 | 说明 |
|------|------|
| isCorrect | 答案是否正确 |
| correctAnswer | 正确答案 |
| explanation | 答案解析 |
| masteryDelta | 掌握度变化值 |

---

## 8. 数据说明

### 卡点指数计算公式
```
卡点指数 = 提问量 × 0.5 + 停留时长 ÷ 20
```

### 掌握度计算
- 初始掌握度：根据历史答题正确率计算
- 每次正确答题：掌握度 +10
- 每次错误答题：掌握度 -5
- 掌握度范围：0-100

---

## 更新记录

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2026-03-04 | v1.1 | 添加学习卡点分析和AI薄弱点功能接口 |

---



## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 注意事项

1. **文件上传**仅支持 PDF、PPT、PPTX 格式
2. **分页接口**默认每页20条，最大支持100条
3. **所有响应**均包含 `code` 字段，`code=200` 表示成功
4. **中文编码**已统一为 UTF-8
5. **时间格式**为 ISO8601: `2006-01-02T15:04:05+08:00`

---

## 更新记录

| 日期 | 版本 | 更新内容 |
|------|------|----------|
| 2026-03-03 | v1.0 | 完成所有教师端接口文档 |