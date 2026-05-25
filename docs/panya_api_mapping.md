# 接口对照与数据映射模板（基于超星AI互动智课服务系统开放API PDF）

说明：此模板用于将现有系统接口与泛雅（学习通）接口逐项对照，记录请求/响应映射、鉴权与测试信息，便于开发与联调。

本文档已根据你提供的 PDF 抽取出一版“目标接口清单”。当前阶段不再依赖原有接口文档，而是以该 PDF 作为唯一接口来源。

## 目标接口清单（已抽取）

| 序号 | 功能模块 | 接口地址 | 方法 | 核心用途 | 关键入参 |
|---:|---|---|---|---|---|
| 1 | 课件上传与解析 | /api/v1/lesson/parse | POST | 上传 PPT/PDF 并返回解析任务与课件结构预览 | schoolId, userId, courseId, fileType, fileUrl, isExtractKeyPoint, enc |
| 2 | 智课脚本生成 | /api/v1/lesson/generateScript | POST | 基于解析结果生成结构化讲授脚本 | parseId, teachingStyle, speechSpeed, customOpening, enc |
| 3 | 语音合成 | /api/v1/lesson/generateAudio | POST | 将讲授脚本转换为语音音频 | scriptId, voiceType, audioFormat, sectionIds, enc |
| 4 | 问答交互 | /api/v1/qa/interact | POST | 学生文字/语音提问并返回多轮问答结果 | schoolId, userId, courseId, lessonId, sessionId, questionType, questionContent, currentSectionId, historyQa, enc |
| 5 | 语音提问识别 | /api/v1/qa/voiceToText | POST | 将学生语音问题转写为文本 | voiceUrl, voiceDuration, language, enc |
| 6 | 学习进度追踪 | /api/v1/progress/track | POST | 记录进度、最近问答和学习轨迹 | schoolId, userId, courseId, lessonId, currentSectionId, progressPercent, lastOperateTime, qaRecordId, enc |
| 7 | 学习节奏调整 | /api/v1/progress/adjust | POST | 根据理解程度与进度返回后续讲授建议 | userId, lessonId, currentSectionId, understandingLevel, qaRecordId, enc |

## 通用约定（来自 PDF）

- 协议：HTTP/HTTPS
- 数据格式：JSON，UTF-8
- 签名：enc = MD5(参数有序拼接 + staticKey + time)
- 时间格式：yyyy-MM-ddHH:mm:ss
- 返回格式：{ code, msg, data, requestId }
- 成功码：200
- 客户端错误：4xx
- 服务端错误：5xx

使用说明：复制下方表格到 Excel/Google Sheets 以便多人填写，或直接在本文件中新增条目。

## 字段说明
- 现有系统端点：本系统的 API 地址或功能名。
- 现有请求示例：JSON 或表单示例。
- 现有响应示例：JSON 示例（成功/失败）。
- 泛雅端点：泛雅对应的 API 地址（或预估接口名称）。
- 泛雅请求映射：字段映射、类型转换说明。
- 泛雅响应映射：如何将泛雅响应转换为本系统模型。
- 鉴权方式：如 `API Key`、`签名`、`OAuth2`、Cookie 等。
- 幂等/唯一键：请求幂等策略与唯一标识字段。
- 错误码映射：泛雅错误码 → 本系统错误码/处理方式。
- 限流/速率：调用频率限制与建议。
- 沙盒/测试URL：测试环境地址与凭证（如有）。
- 备注：其他注意事项（合规、时间窗口、批量限制等）。

---

## 表格模板（Markdown）

| 序号 | 现有系统端点（功能） | 方法 | 现有请求示例 | 现有响应示例 | 泛雅端点（预估） | 泛雅方法 | 请求映射（字段对应） | 响应映射 | 鉴权 | 幂等 | 错误码映射 | 限流 | 沙盒/测试URL | 备注 |
|---:|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| 1 | 获取课程列表 | GET | `GET /api/v1/courses?userId=...` | `{ "courses": [...] }` | `/course/list` | GET | `userId -> userId` | `courses -> courses` | API Key | 是（query:userId） | 200/0 -> OK; 401 -> AUTH_FAIL | 1000/min | https://sandbox.fanya.example/api | 示例 |


## CSV 模板（表头，便于导入）

现有系统端点,方法,现有请求示例,现有响应示例,泛雅端点,泛雅方法,请求映射,响应映射,鉴权,幂等,错误码映射,限流,沙盒测试URL,备注


## 示例条目（填写示例）
- 序号：1
- 现有系统端点（功能）：创建作业
- 方法：POST
- 现有请求示例：
```
POST /api/v1/homework
{
  "courseId": "c123",
  "title": "第1次作业",
  "dueDate": "2026-06-01T23:59:00Z",
  "questions": [ ... ]
}
```
- 现有响应示例：
```
{ "code": 0, "data": { "homeworkId": "h456" } }
```
- 泛雅端点（预估）：`/homework/create`
- 泛雅方法：POST
- 请求映射：`courseId -> courseId`, `title -> title`, `dueDate -> deadline (格式: yyyy-MM-dd)`，`questions -> items`（题型需转换）
- 响应映射：`泛雅返回 {"id":"..."} 映射为 data.homeworkId`
- 鉴权：签名令牌（client_id/client_secret + 签名）
- 幂等：客户端传 `clientRequestId` 保证幂等
- 错误码映射：泛雅 400 -> 参数错误 -> 本系统 1002；泛雅 429 -> 速率限制 -> 重试/退避
- 限流：每用户每分钟 30 次
- 沙盒/测试URL：https://sandbox.fanya.example/homework/create
- 备注：题库题型需做映射表

---

## 联调/验收检查清单（建议）
- [ ] 获取/确认沙盒 URL 与测试账号
- [ ] 确认鉴权方式并获得测试凭证
- [ ] 填写接口映射表（优先基础用户/课程/作业/成绩/资源类接口）
- [ ] 由开发实现适配层并在沙盒完成端到端测试
- [ ] 记录异常场景与错误码映射，并设计重试策略
- [ ] 性能/并发测试（按 SLA 预估）


如需，我可以把该模板转换为 Excel (`docs/panya_api_mapping.xlsx`) 并填入示例行，或者把首批 10 个关键接口的映射根据我检索到的公开线索先行预填一版。