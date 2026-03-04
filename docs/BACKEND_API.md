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