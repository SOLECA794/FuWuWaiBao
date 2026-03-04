# 智能学习平台接口说明文档

## 一、AI 解析与问答模块

| 接口地址 | 请求方式 | 功能描述 | 请求参数（JSON 格式） |
|---------|---------|---------|----------------------|
| `/api/ai/parseKnowledge` | POST | 解析 PDF/PPTX 文件内容，拆解知识点层级结构（章节→知识点→子知识点） | `{ "fileContent": "解析后的文件文本内容", "fileType": "pdf/pptx", "studentId": "2025001" }` |
| `/api/ai/question` | POST | 多模态提问（文字/语音/图片），AI 智能答疑 | `{ "courseId": 1, "pageNum": 3, "question": "fillna怎么用？", "type": "text" }` |
| `/api/ai/traceQuestion` | POST | 溯源定位提问（精准定位课件区域提问） | `{ "courseId": 1, "pageNum": 3, "x": 200, "y": 150, "question": "这里不懂" }` |

---

## 二、薄弱点诊断与练习模块

| 接口地址 | 请求方式 | 功能描述 | 请求参数（JSON 格式） |
|---------|---------|---------|----------------------|
| `/api/weakPoint/getList` | GET | 获取指定学生的 AI 诊断薄弱点列表 | `{ "studentId": "2025001", "courseId": "PY202501" }` |
| `/api/weakPoint/getExplain` | POST | 获取指定薄弱点的详细讲解内容 | `{ "weakPointName": "缺失值填充", "studentId": "2025001" }` |
| `/api/weakPoint/getTest` | POST | 根据薄弱点生成对应的检测习题（单选/多选） | `{ "weakPointName": "缺失值填充", "studentId": "2025001", "questionType": "single" }` |
| `/api/weakPoint/checkAnswer` | POST | 校验习题答案，并返回解析和掌握状态 | `{ "studentId": "2025001", "questionId": "Q20250101", "userAnswer": "interpolate()" }` |

---

## 三、课件与学习数据模块

| 接口地址 | 请求方式 | 功能描述 | 请求参数 |
|---------|---------|---------|---------|
| `/api/courseware/page` | POST | 获取指定页码的课件信息（翻页核心接口） | `{ "courseId": 1, "currentPage": 3 }` |
| `/api/teacher/card-data/{courseId}` | GET | 获取指定课件的学习卡点 Log 数据 | 路径参数：`courseId`（课程 ID，字符串/数字） |
| `/api/student/studyData` | GET | 获取学生课程学习数据统计（专注度/薄弱点等） | `{ "studentId": 2025001, "courseId": 1 }` |
| `/api/student/breakpoint` | GET | 获取学生上次学习的断点页码（用于续播） | `{ "studentId": 2025001, "courseId": 1 }` |
| `/api/student/saveNote` | POST | 保存学生课件笔记（扩展功能） | `{ "studentId": 2025001, "pageNum": 3, "note": "fillna()填充缺失值" }` |

---

## 核心说明
- **请求方式**：GET 请求参数通常拼接在 URL 或 Query 中；POST 请求参数统一放在 Request Body 中，格式为 JSON。
- **必填字段**：所有接口中标记为“必填”的字段不可缺失，否则接口会返回参数校验错误。
- **数据类型**：注意区分字符串（如 `"2025001"`）和数字（如 `1`）类型，避免类型不匹配导致接口调用失败。

---

要不要我帮你把这份接口说明整理成一份**可直接导入 Postman 的集合文件**，方便你和后端联调？