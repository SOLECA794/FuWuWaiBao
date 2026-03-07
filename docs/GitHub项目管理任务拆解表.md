# GitHub项目管理任务拆解表

> 用途：把当前项目拆成可以直接在 GitHub 上创建 issue 的任务。
> 使用方式：每一行都可以直接建成一张卡，标题可直接复用，标签和依赖也已建议好。
> 更新时间：2026-03-06

---

## 1. 推荐建卡方式

每个 issue 建议至少包含以下字段：

- 标题
- 背景
- 交付物
- 影响目录
- 依赖任务
- 验收标准
- 标签
- 负责人

建议标题格式：

- `[Backend] 修复 /api/v1 教师端路由注册`
- `[AI] 实现解析任务消费 Worker`
- `[Student FE] 接入真实 SSE 问答流`

---

## 2. Epic 划分

| Epic | 目标 | 建议里程碑 |
|---|---|---|
| E0 基线稳定 | 修正目录、路由、文档、启动方式不一致问题 | P0-稳定联调 |
| E1 课件闭环 | 打通上传、解析、写库、讲稿、预览 | P1-MVP闭环 |
| E2 学生学习闭环 | 打通播放、问答、断点、续接 | P1-MVP闭环 |
| E3 AI 能力增强 | 真实流式、知识点、弱点诊断、重讲 | P2-体验增强 |
| E4 开放平台与安全 | 开放接口、签名、安全、可观测性 | P2-体验增强 |

---

## 3. 当前建议优先级

### 3.1 必须先做的 P0

- 路由与接口基线修正
- 文档和真实目录同步
- 前端统一接口体系
- 启动和联调路径跑通

### 3.2 然后做的 P1

- 上传后真实解析
- 页级脚本写库
- 学生问答闭环
- 学习断点和续接

### 3.3 最后做的 P2

- 真流式 SSE
- 弱点诊断优化
- 开放适配层补强
- 安全与观测

---

## 4. 任务拆解表（可直接建 issue）

| ID | Epic | 建议 Issue 标题 | 负责人 | 优先级 | 依赖 | 影响目录 | 标签 | DoD |
|---|---|---|---|---|---|---|---|---|
| T001 | E0 | `[Backend] 审计并修复 /api/v1 路由注册问题` | Go | P0 | 无 | `backend/api/main.go` | `area/backend` `type/bug` `priority/P0` | 所有文档声明的 `/api/v1` 路由可命中；列出修复项；前端主要请求不再 404 |
| T002 | E0 | `[Docs] 同步当前真实目录与启动方式` | 文档/负责人 | P0 | T001 | `README.md` `docs/*.md` | `area/docs` `type/task` `priority/P0` | README 和 docs 中不再出现已废弃目录；新人按文档可启动项目 |
| T003 | E0 | `[Student FE] 统一学生端接口调用到明确版本` | 学生前端 | P0 | T001 | `frontend/student/src/services/**` | `area/frontend-student` `type/refactor` `priority/P0` | 学生端明确只走一套主链路接口；保留兼容调用需有注释说明 |
| T004 | E0 | `[Teacher FE] 验证教师端 /api/v1 接口全链路可用` | 教师前端 | P0 | T001 | `frontend/teacher/src/**` | `area/frontend-teacher` `type/task` `priority/P0` | 列表、上传、保存讲稿、发布、统计至少手测通过一遍 |
| T005 | E0 | `[DevOps] 补齐本地联调环境说明和默认端口约定` | 负责人/后端 | P0 | T002 | `README.md` `docs/联调清单.md` | `area/devops` `area/docs` `type/task` `priority/P0` | 明确 Go、AI、前端端口；docker 启动方式可复制运行 |
| T006 | E0 | `[Security] 梳理开放接口签名与鉴权缺口` | Go | P0 | T002 | `backend/internal/handler/compat_*.go` | `area/backend` `type/task` `priority/P0` | 输出安全缺口清单，至少覆盖签名、重放、鉴权、CORS |
| T007 | E1 | `[Backend] 为 Course 增加解析状态与任务标识` | Go | P1 | T001 | `backend/internal/model/models.go` `backend/api/main.go` | `area/backend` `type/feature` `priority/P1` | 课件具备 `status/taskId/errorMessage` 类状态字段或等价设计；上传后状态可查询 |
| T008 | E1 | `[Backend] 上传课件后入解析队列` | Go | P1 | T007 | `backend/internal/service/course.go` `backend/internal/repository/**` | `area/backend` `type/feature` `priority/P1` | 上传成功后产生解析任务；状态变为 parsing；Redis 可见任务 |
| T009 | E1 | `[AI] 实现解析任务消费 Worker` | Python AI | P1 | T008 | `ai_engine/**` | `area/ai` `type/feature` `priority/P1` | Worker 能消费任务并拉取课件内容；失败时有错误日志 |
| T010 | E1 | `[AI] 输出统一解析结果 schema` | Python AI | P1 | T009 | `ai_engine/parser.py` `ai_engine/schema.py` | `area/ai` `type/feature` `priority/P1` | 输出包含页级内容、总页数、可持久化字段；样例固化到文档 |
| T011 | E1 | `[AI] 基于解析结果生成页级脚本节点` | Python AI | P1 | T010 | `ai_engine/generator.py` `ai_engine/generate.py` | `area/ai` `type/feature` `priority/P1` | 每页生成 nodes 结构，至少包含 opening/explain/transition |
| T012 | E1 | `[Backend] 解析结果写入 CoursePage 和 KnowledgePoint` | Go | P1 | T010 | `backend/internal/model/**` `backend/internal/handler/**` | `area/backend` `type/feature` `priority/P1` | 解析完成后可通过接口查到 page/script/knowledge graph |
| T013 | E1 | `[Teacher FE] 补充解析状态展示和刷新机制` | 教师前端 | P1 | T007 T012 | `frontend/teacher/src/**` | `area/frontend-teacher` `type/feature` `priority/P1` | 教师端可看到 parsing/ready/failed；失败有明确提示 |
| T014 | E1 | `[Backend] 提供解析状态查询接口` | Go | P1 | T007 | `backend/internal/handler/**` | `area/backend` `type/feature` `priority/P1` | 前端无需猜测状态；接口返回稳定状态结构 |
| T015 | E2 | `[Backend] 统一学生学习会话状态结构` | Go | P1 | T003 | `backend/internal/handler/compat_common.go` `backend/internal/model/**` | `area/backend` `type/refactor` `priority/P1` | 明确 sessionId/currentPage/currentNodeId/interrupted/resume* 字段 |
| T016 | E2 | `[Backend] Redis 持久化学习进度与续接状态` | Go | P1 | T015 | `backend/internal/repository/redis.go` `backend/internal/handler/**` | `area/backend` `type/feature` `priority/P1` | 问答前后状态能落 Redis；服务重启后不会全部丢失 |
| T017 | E2 | `[Student FE] 播放器接入 currentPage/currentNodeId 上报` | 学生前端 | P1 | T015 | `frontend/student/src/**` | `area/frontend-student` `type/feature` `priority/P1` | 进入学习、切页、切节点均能更新后端进度 |
| T018 | E2 | `[Backend] 完善学生脚本接口的 nodes 输出` | Go | P1 | T012 | `backend/internal/handler/compat_student.go` | `area/backend` `type/feature` `priority/P1` | 学生端能拿到稳定 nodes，字段名固定 |
| T019 | E2 | `[Student FE] 完成问答结束后的自动续接` | 学生前端 | P1 | T016 T018 | `frontend/student/src/**` | `area/frontend-student` `type/feature` `priority/P1` | `final` 事件后页面自动恢复到正确 page/node |
| T020 | E2 | `[Backend] 记录 QuestionLog 并串联学习统计` | Go | P1 | T016 | `backend/internal/handler/student.go` `backend/internal/handler/teacher.go` | `area/backend` `type/feature` `priority/P1` | 学情、提问记录、卡点统计使用真实日志数据 |
| T021 | E3 | `[AI] 提供真实流式问答接口` | Python AI | P1 | T010 | `ai_engine/main.py` `ai_engine/qa.py` | `area/ai` `type/feature` `priority/P1` | AI 服务按 token 或 sentence 流式返回，而不是一次性整包 |
| T022 | E3 | `[Backend] 将 Python 流式问答透传为 SSE` | Go | P1 | T021 | `backend/internal/service/ai_engine_client.go` `backend/internal/handler/**` | `area/backend` `type/feature` `priority/P1` | Go 不再先拿完整答案再拆分；首 token 明显提前 |
| T023 | E3 | `[Student FE] 对接 token/sentence/final/error 四类 SSE 事件` | 学生前端 | P1 | T022 | `frontend/student/src/services/**` `frontend/student/src/**` | `area/frontend-student` `type/feature` `priority/P1` | UI 正确处理四类事件，异常时可恢复 |
| T024 | E3 | `[AI] 完善知识点树与弱点诊断数据来源` | Python AI | P2 | T010 T020 | `ai_engine/**` | `area/ai` `type/feature` `priority/P2` | 弱点不再仅靠提问次数估算；有明确计算逻辑 |
| T025 | E3 | `[Backend] 重构弱点诊断持久化与题目生成流程` | Go | P2 | T024 | `backend/internal/handler/weakpoint.go` `backend/internal/model/**` | `area/backend` `type/refactor` `priority/P2` | 题目、答题记录、掌握度变化可追踪 |
| T026 | E3 | `[Teacher FE] 优化学情分析页的数据来源与展示` | 教师前端 | P2 | T020 T025 | `frontend/teacher/src/**` | `area/frontend-teacher` `type/feature` `priority/P2` | 页面只展示真实统计，不展示伪数据 |
| T027 | E4 | `[Backend] 为开放适配层增加时效校验与重放保护` | Go | P2 | T006 | `backend/internal/handler/compat_common.go` | `area/backend` `type/feature` `priority/P2` | `time` 有有效期校验；重复请求可识别 |
| T028 | E4 | `[Backend] 为开放适配层补充集成测试样例` | Go | P2 | T027 | `backend/**` `docs/**` | `area/backend` `area/docs` `type/task` `priority/P2` | 至少覆盖 parse、interact、track 三类接口 |
| T029 | E4 | `[Docs] 生成联调用 Postman 或 HTTP 样例集合` | 文档/后端 | P2 | T014 T022 | `docs/**` | `area/docs` `type/task` `priority/P2` | 队友无需猜请求格式；关键接口有现成样例 |
| T030 | E4 | `[QA] 建立端到端冒烟验收清单` | 负责人/测试 | P1 | T013 T019 T022 | `docs/联调清单.md` | `area/docs` `type/task` `priority/P1` | 教师上传到学生问答续接可按步骤验收 |

---

## 5. 建议先创建的 10 张卡

如果你不想一次建 30 张卡，先建下面这 10 张就够了：

1. T002 同步当前真实目录与启动方式
2. T003 统一学生端接口调用到明确版本
3. T004 验证教师端 `/api/v1` 接口全链路可用
4. T007 为 Course 增加解析状态与任务标识
5. T008 上传课件后入解析队列
6. T009 实现解析任务消费 Worker
7. T012 解析结果写入 CoursePage 和 KnowledgePoint
8. T016 Redis 持久化学习进度与续接状态
9. T021 提供真实流式问答接口
10. T022 将 Python 流式问答透传为 SSE

这 10 张卡基本决定项目能不能从“演示代码”进入“可联调交付”。

---

## 6. 指派建议

### Go 后端同学

- T001 T006 T007 T008 T012 T014 T015 T016 T018 T020 T022 T025 T027 T028

### Python AI 同学

- T009 T010 T011 T021 T024

### 教师端前端同学

- T004 T013 T026

### 学生端前端同学

- T003 T017 T019 T023

### 项目负责人或文档负责人

- T002 T005 T029 T030

---

## 7. Done 标准

每张卡关闭前，至少满足以下 4 条：

1. 有代码提交或文档提交。
2. 有最少一条可复现的验收方式。
3. 如果改了接口，已同步接口文档。
4. 没有把依赖问题留给下一个人“猜”。
