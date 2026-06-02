# PPT合规审核修改清单.md

> 版本：V1.0（工作流生成）

---

# 附录A3：技术附件

> 本附件为小册子正文补充技术细节，包括完整图系索引、关键Schema定义摘要、接口规范摘要与量化参数参考表。内容均从项目源代码和文档中提取，可作为答辩评委的技术深度佐证材料。

---

## A3.1 完整图系索引

图系资产集中管理于 `docs-selected/图系资产/`，共10张独立图资产，每张图对应独立的 `.mmd` 源文件，支持独立编辑与批量渲染。

| 编号 | 图名 | 类型 | 对应模块 | 对应PPT | 小册子章节 | 核心用途 |
|:----:|------|:----:|:--------:|:-------:|:----------:|---------|
| D01 | 五环业务闭环总览图 | 业务流 | 总体方案 | P02 | M02 | 开场建立系统整体认知——课件上传→解析→脚本生成→语音合成→学生互动→学情回传的因果链 |
| D02 | 系统技术编排架构全景图 | 架构图 | 架构设计 | P09-P10 | M03 | 四层分层架构（前端→编排→AI→数据）+ 三个关键设计决策标注（薄业务重编排、双引擎降级、SSE三事件协议） |
| D03 | 三阶段解耦解析管线图 | 算法管线 | 课件解析 | P06/P11/P12 | M04 | Stage1 Markdown理解→Stage2节点树构建→Stage3逐节点脚本生成，每阶段独立JSON Schema约束 |
| D04 | 本页优先RAG+问答流程图 | 算法流 | 智能问答 | P07/P13 | M04 | 三级级联上下文（节点→页面→全页）+ 苏格拉底式约束 + SSE三事件流式输出 |
| D05 | 多态断点续学状态机图 | 状态机 | 续学算法 | P08/P13 | M04 | 学习主状态四态流转 + 理解程度三态（none/partial/full）对应三种续接策略 |
| D06 | 双引擎降级时序图 | 时序图 | 可用性设计 | P05/P10 | M03 | 正常路径→Dify→降级→本地引擎→双失败兜底的完整时序 |
| D07 | 中间件拓扑与数据流图 | 架构图 | 数据管理 | P09 | M03 | 已合并至D02，标注PG/Redis/MinIO的读写方向与缓存策略 |
| D08 | 学情多维加权指数示意图 | 数据可视化 | 学情分析 | P06 | M04 | 四个维度加权（提问频次、需重讲次数、会话数、停留时长）+ 掌握度增量模型 |
| D09 | 课件上传容错管道图 | 流程图 | 课件管道 | P12 | M03 | 三级预览图降级 + AI解析不阻塞上传主流程 |
| D10 | 任务调度状态机图 | 状态机 | 后台任务 | P15 | M03 | 六态流转（pending→queued→processing→success/failed），指数退避重试 |

**图系资产目录结构：**
```
docs-selected/图系资产/
├── README.md                     ← 目录索引
├── 00_Diagram_Assets_Dictionary.md  ← 数据字典
├── 00_Diagram_Usage_Guide.md        ← 讲授说明（含质询防御）
├── source/                          ← 独立.mmd源文件
│   ├── 01_Business_Loop.mmd
│   ├── 02_System_Architecture.mmd
│   ├── 03_Parse_Pipeline.mmd
│   ├── 04_RAG_QA_Flow.mmd
│   ├── 05_Resume_StateMachine.mmd
│   ├── 06_Fallback_Sequence.mmd
│   ├── 07_Middleware_Topology.mmd
│   ├── 08_Mastery_Score.mmd
│   ├── 09_Upload_Pipeline.mmd
│   └── 10_Task_Scheduler.mmd
└── reference/                       ← 量化锚点数据集
    ├── 01_architecture_anchors.md
    └── 02_algorithm_anchors.md
```

---

## A3.2 关键Schema定义摘要

三阶段解析管线是系统的核心算法创新。每个阶段输出受JSON Schema严格约束，采用 `additionalProperties: false` 严格模式，大模型输出格式异常时触发自动修复。

### A3.2.1 Stage1 — Markdown理解Schema

**源码路径：** `ai_engine/schema.py :: build_stage1_markdown_schema()`

- 字段：`normalized_markdown` (string) 整理后的规范化Markdown
- 字段：`key_points` (array of string) 关键知识点，5-12条
- 约束：严格模式，仅允许两个字段
- 用途：将原始课件内容提炼为结构化Markdown + 知识要点列表，作为后续阶段的输入

### A3.2.2 Stage2 — 节点树构建Schema

**源码路径：** `ai_engine/schema.py :: build_stage2_node_tree_schema()`

- 字段：`nodes` (array) 知识节点数组
  - `node_id` (string) 节点标识符，格式如 `node_001`
  - `title` (string) 节点标题
  - `summary` (string) 节点摘要
  - `source_span` (string) 源范围定位
  - `prerequisites` (array of string) 前置依赖节点ID列表
- 约束：nodes数组内每个item均 `additionalProperties: false`
- 用途：将整理后的Markdown拆解为层级化知识节点树，每节点是独立的教学单元

### A3.2.3 Stage3 — 脚本生成Schema（基础版）

**源码路径：** `ai_engine/schema.py :: build_stage3_script_schema()`

- 字段：`scripts` (array) 讲稿脚本数组
  - `node_id` (string) 绑定节点ID
  - `title` (string) 节点标题
  - `script` (string) 完整讲稿文本
  - `segments` (array) 段落切片数组
    - `segment_id` (string) 段落标识符
    - `text` (string) 段落文本
    - `node_id` (string) 段落归属节点
- 约束：scripts和segments均严格模式

### A3.2.4 Stage3增强版 — 增强脚本Schema

**源码路径：** `ai_engine/schema.py :: build_enhanced_script_schema()`

在基础版之上新增字段：
- `segment_type` (string, enum: `opening|explanation|example|interaction|transition|summary`) 六种段落类型
- `estimated_seconds` (integer) 预估讲授时长（15-120秒）
- `difficulty` (string, enum: `easy|medium|hard`) 难度等级
- `interaction_hint` (string) 教师互动话术
- `knowledge_card` (object) 知识卡片
  - `term` (string, required) 术语
  - `definition` (string, required) 定义
  - `formula` (string) 公式
  - `tags` (array of string) 标签
- `node_ids` (array of string) 段落可归属多个节点的扩展

### A3.2.5 节点脚本完整Schema

**源码路径：** `ai_engine/schema.py :: build_node_script_schema()`

包含字段：`node_id`, `title`, `script`, `mindmap_markdown`, `interactive_questions`, `reteach_script`, `transition`, `structured_markdown`, `knowledge_nodes`, `script_segments`

知识节点规范化包含：`node_id`, `parent_id`, `level`, `title`, `tags`, `prerequisites`, `difficulty`, `coverage_span`

脚本段落规范化包含：`segment_id`, `text`, `node_ids` (多节点关联), `confidence`, `manual_override`

### A3.2.6 核心修复机制

| 机制 | 函数 | 说明 |
|------|------|------|
| 节点归一化 | `normalize_stage2_nodes()` | 去重、补缺、标准化node_id |
| 脚本清洗 | `normalize_stage3_scripts()` | 剔除无效node_id引用，空段过滤 |
| 双向对齐 | `align_node_segment_mapping()` | 段落node_ids与节点coverage_span双向校准 |
| 文档Schema | `build_document_schema()` | 解析结果封装，含统计信息 |

---

## A3.3 接口规范摘要

系统采用多版本API架构，兼容原有接口协议（泛雅OpenAPI）的同时提供标准化RESTful v1接口。

### A3.3.1 API整体分层

| 接口类别 | 路由前缀 | 协议 | 鉴权方式 | 说明 |
|---------|---------|------|---------|------|
| 健康检查 | `/api/health` | GET | 无 | 服务可用性探测 |
| 教师端Legacy | `/api/teacher/*` | REST | Session/JWT | 教师端原始接口 |
| 学生端Legacy | `/api/student/*` | REST | Session/JWT | 学生端原始接口 |
| 教师端V1 | `/api/v1/teacher/*` | REST | JWT | 标准化教师接口 |
| 学生端V1 | `/api/v1/student/*` | REST | JWT | 标准化学生接口 |
| OpenAPI对接 | `/api/v1/lesson/*`, `/api/v1/qa/*` | REST | MD5签名 | 泛雅平台外部接入 |
| 弱项管理 | `/api/weakPoint/*` | REST | JWT | 薄弱点分析与练习 |
| AI开放接口 | `/api/v1/ai/*` | REST | JWT | AI能力公共入口 |

### A3.3.2 通用响应格式

```json
// 成功响应（v1风格）
{"code": 200, "message": "请求成功", "data": {...}}

// 成功响应（OpenAPI风格）
{"code": 200, "msg": "请求成功", "data": {...}, "requestId": "req_uuid"}

// 错误响应
{"code": 400, "message": "参数错误", "errors": "detail"}
```

### A3.3.3 OpenAPI签名机制

**源码路径：** `backend/internal/handler/compat_common.go :: OpenAPISignatureMiddleware()`

- 签名字段：`enc` + `time` 必需
- 签名算法：MD5(排序参数字段值拼接 + staticKey + timeValue)
- 签名参数字段：`platformId`, `userId`（可通过 `OPEN_API_SIGN_FIELDS` 环境变量配置）
- 验证失败返回 HTTP 403

### A3.3.4 SSE流式问答协议

**关键端点：** `POST /api/student/qa/stream` 及 v1 对应接口

三事件流式输出格式：
```
event: token
data: {"text": "逐词"}

event: sentence
data: {"text": "完整句子"}

event: final
data: {"session_id": "...", "understanding_level": "none|partial|full", "source_page": N, "resume_page": N, "resume_node_id": "...", "resume_sec": N, "follow_up_suggestion": "..."}
```

### A3.3.5 OpenAPI对外接口清单

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 课件解析 | POST | `/api/v1/lesson/parse` | 解析远端课件文件 |
| 脚本生成 | POST | `/api/v1/lesson/generateScript` | 基于解析结果生成讲稿 |
| 语音合成 | POST | `/api/v1/lesson/generateAudio` | 基于讲稿生成音频 |
| 问答交互 | POST | `/api/v1/qa/interact` | 课程知识问答 |
| 语音识别 | POST | `/api/v1/qa/voiceToText` | 语音转文字 |
| 进度追踪 | POST | `/api/v1/qa/trackProgress` | 记录学习进度 |
| 节奏调整 | POST | `/api/v1/qa/adjustProgress` | 基于理解程度调整 |
| 课程同步 | POST | `/api/v1/platform/syncCourse` | 同步学校平台课程 |
| 用户同步 | POST | `/api/v1/platform/syncUser` | 同步学校平台用户 |
| 平台总览 | GET | `/api/v1/platform/overview` | 平台数据统计 |

---

## A3.4 量化参数参考表

以下参数从项目源码中提取，作为技术深度佐证材料。

### A3.4.1 编排层核心参数

| 参数名 | 值 | 源码位置 | 说明 |
|--------|:--:|---------|------|
| SSE首包响应时间 | <500ms | 编排层流式转发设计 | 从用户提问到收到第一个token的时间 |
| 对话上下文窗口 | 4轮 | `student.go` QAStream `buildDialogueContext` limit参数 | 保存最近4轮问答作为上下文 |
| 编排层响应特性 | ms级 | 架构设计决策 | 编排层不做AI推理，仅做状态编排与流式转发 |
| 学习状态机主状态 | 4态 | `student.go` QAStream | 讲授→遇疑→问答→续接，循环往复 |

### A3.4.2 AI能力层核心参数

| 参数名 | 值 | 源码位置 | 说明 |
|--------|:--:|---------|------|
| AI推理典型延迟 | 3-5秒 | 通用模型调用经验 | 大模型单次推理的端到端耗时区间 |
| 生成温度 | 0.2 | `generator.py` GenerationConfig | 控制LLM输出的随机性，低温度保证教学一致性 |
| 内容截断阈值 | 3500字符 | `generator.py` `max_content_chars` | 超过此值时截断并保留首尾 |
| 模型默认 | qwen-turbo | 环境变量 `AI_MODEL` | 当前接入的通义千问模型 |
| Dify默认地址 | http://127.0.0.1:18001 | `dify_client.go` NewDifyClient | Dify服务本地部署地址 |
| 降级引擎地址 | http://127.0.0.1:8000 | `dify_client.go` postJSON | 本地Python引擎地址 |
| Dify请求超时 | 60秒 | `dify_client.go` HTTP client Timeout | 单次Dify调用超时时间 |
| 理解程度三态 | none / partial / full | `qa.py` 常量定义 | 未理解/部分理解/已理解 |
| 关键词降级匹配集 | 7个 | `qa.py` RETEACH_KEYWORDS | 含"听不懂""不懂""太难"等教育场景关键词 |

### A3.4.3 掌握度模型参数

| 参数名 | 值 | 源码位置 | 说明 |
|--------|:--:|---------|------|
| 掌握度初始值 | 0.5 | `knowledge_map.go` UpdateMasteryScore | 新知识点的初始掌握度 |
| 答对增量 | +0.08 | `knowledge_map.go` | 答对题目后的掌握度增加量 |
| 答错减量 | -0.06 | `knowledge_map.go` | 答错题目后的掌握度减少量 |
| 快速响应加分阈值 | ≤5000ms | `knowledge_map.go` | 响应时间低于此值获额外加分 |
| 快速响应加分值 | +0.02 | `knowledge_map.go` | 快速响应的额外掌握度加成 |
| 掌握度取值范围 | [0, 1] | `knowledge_map.go` clampScore | 强制截断确保取值合法 |
| 薄弱点阈值 | <0.6 | `knowledge_map.go` GetWeakKnowledgePoints | 掌握度低于此值判定为薄弱 |
| 强项点阈值 | ≥0.8 | `knowledge_map.go` GetStrongKnowledgePoints | 掌握度高于此值判定为强项 |

### A3.4.4 解析管线参数

| 参数名 | 值 | 源码位置 | 说明 |
|--------|:--:|---------|------|
| 解析管线阶段数 | 3 | `generator.py` generate_from_markdown | 理解→建树→生成 |
| 阶段1输出格式 | normalized_markdown + key_points | `schema.py` build_stage1_markdown_schema | 规范化Markdown + 知识点列表 |
| 阶段2节点属性 | 5字段 | `schema.py` build_stage2_node_tree_schema | node_id/title/summary/source_span/prerequisites |
| 阶段3脚本段类型 | 6种 | `schema.py` build_enhanced_script_schema | opening/explanation/example/interaction/transition/summary |
| Schema校验模式 | strict + additionalProperties:false | 所有Schema定义 | 大模型输出格式异常时触发自动修复 |
| LLM失败回退 | 预置模板fallback | `generator.py` _fallback_pipeline | 任一阶段失败不阻塞整体流程 |
| 解析效率参考 | ≤2分钟/份 | `course.go` enrichCourseWithAI | 典型课件解析耗时 |

### A3.4.5 问答系统参数

| 参数名 | 值 | 源码位置 | 说明 |
|--------|:--:|---------|------|
| 上下文级联层数 | 3级 | `dialogue_support.go` buildNodeScopedContext | L1节点脚本→L2整页节点→L3全页原文 |
| L1上下文范围 | 目标节点ScriptSegments | buildNodeScopedContext | 优先使用节点级教学脚本 |
| L2上下文范围 | 当前页全部TeachingNode | buildNodeScopedContext | 整页教学节点补全 |
| L3上下文范围 | CoursePage.SourceText | buildNodeScopedContext | 原始课件全文兜底 |
| 苏格拉底约束 | 严禁直接给答案 | `qa.py` 系统提示词 | 提示词中明确约束回答方式 |
| 问答模式 | llm优先，降级兜底 | `qa.py` _llm_answer | API不可用自动回退到规则匹配 |

### A3.4.6 任务调度参数

| 参数名 | 值 | 源码位置 | 说明 |
|--------|:--:|---------|------|
| 最大重试次数 | 3次 | `task_scheduler.go` processTask | 单任务失败后的最大重试次数 |
| 退避基数 | 2^retryCount 分钟 | `task_scheduler.go` processTask | 第1次2分钟、第2次4分钟、第3次8分钟 |
| 并发协程数 | 10个 | `task_scheduler.go` Start | 10个worker协程消费任务队列 |
| 任务队列容量 | 1000 | `task_scheduler.go` taskQueue | 缓冲队列容量，压力削峰 |
| 调度扫描周期 | 1分钟 | `task_scheduler.go` dispatchLoop | 定时扫描pending任务入队 |
| 单次扫描上限 | 200个 | `task_scheduler.go` checkPendingTasks | 单次扫描的最大加载量 |
| 任务优先级 | 1-5 | `model.ScheduledTask` Priority | 优先级排序后调度执行 |

### A3.4.7 中间件与存储参数

| 参数名 | 值 | 源码位置 | 说明 |
|--------|:--:|---------|------|
| 内容哈希失效机制 | SHA-256 | `models.go` AudioAsset.SourceScriptHash | 监听音频变更的依据 |
| 会话缓存TTL | 24小时 | `compat_common.go` persistSession | Redis中session过期时间 |
| 预览图缓存 | 120秒（Cache-Control） | `course.go` GetPagePreview | 预览图HTTP缓存头 |
| 预览图降级层数 | 3级 | `course.go` ensurePreviewFallback | 实时渲染→占位图→空状态 |
| 音频默认格式 | mp3 | `audio_support.go` defaultAudioFormat | TTS输出格式 |
| 默认语音类型 | 女性中文 | `audio_support.go` defaultAudioVoiceType | TTS发言人配置 |
| 默认TTS提供商 | mock-tts | `audio_support.go` | 语音合成提供商 |

### A3.4.8 数据模型核心字段速查

**Course（课件表）：** id, title, file_url, file_type, total_page, teaching_course_id, is_published
**CoursePage（课件页表）：** course_id, page_index, image_url, source_text, script_text, audio_url, audio_status, audio_duration_sec
**TeachingNode（教学节点表）：** course_id, node_id, title, summary, script_text, reteach_script, interactive_questions, knowledge_nodes_json, script_segments_json, schema_version
**DialogueSession（会话表）：** user_id, course_id, current_page, current_node_id, current_time_sec, playback_mode
**DialogueTurn（问答轮次表）：** session_id, turn_index, page_index, node_id, question, answer, source_page, need_reteach
**StudentKnowledgeMastery（学生掌握度表）：** student_id, knowledge_point_id, mastery_score (0.5初始, [0,1]), correct_count, incorrect_count, last_response_ms
**ScheduledTask（定时任务表）：** task_type, task_data, scheduled_at, status (pending/queued/processing/completed/failed), max_retries (3), retry_count

---

# 附录A4：答辩备用问答库

> 本问答库按维度分类，收录评委高频质疑问题与标准应答口径。每个问题包含：预期质疑角度、标准回答口径（3句话以内）、关联模块索引。

---

## A4.1 架构类

### Q1. 为什么用四层架构？是不是过度设计了？

**标准回答口径：** 四层架构的核心驱动力是"不可控延迟隔离"——AI推理的3-5秒延迟不能拖慢状态编排的毫秒级响应。分层之后，编排层保持轻量只做状态机流转和流式转发，AI能力下沉到独立层，单层故障不影响其他层。这不是为了分层而分层，而是每层解决一个明确的系统约束。

**关联模块索引：** M03架构深度解析 / D02系统架构全景图 / 编排层锚点

### Q2. 为什么选Go而不是Java/Python？

**标准回答口径：** 选Go是因为编排层的核心职责是高并发SSE流式转发和状态机管理，Go的goroutine模型天然适合这类I/O密集型场景，单机即可支撑千级并发连接。教程讲解场景不需要Java的庞大生态和Spring Boot的启动开销，也不需要Python的GIL限制。Go的编译效率和部署简洁性也适合校内私有化部署场景。

**关联模块索引：** M03架构深度解析 / M05工程开发设计规范

### Q3. 为什么不做一个All-in-One系统，要把功能拆到多个组件里？

**标准回答口径：** All-in-One意味着单点故障——AI引擎挂了整个系统不可用，存储满了影响前端响应。我们的设计核心是"故障隔离"：Go编排层做状态管理，Dify/Python做AI推理，PostgreSQL/Redis/MinIO各管一类数据。Dify挂了自动降级到本地引擎，前端继续运行不受影响——这是All-in-One不可能做到的韧性。

**关联模块索引：** M03架构·降级设计 / D06双引擎降级时序图

### Q4. 双引擎降级具体怎么保证的？切换时用户会感觉到吗？

**标准回答口径：** 降级逻辑封装在 `dify_client.go` 的 `postJSON` 方法中——每次调用先请求Dify主引擎，非2xx响应或超时时自动切换到 `http://127.0.0.1:8000` 的本地Python引擎。切换在编排层内部完成，前端接收同样的响应格式，完全无感知。只有两个引擎都失败才返回降级状态码，由前端做最终兜底提示。

**关联模块索引：** M03架构·降级设计 / D06双引擎降级时序图 / `dify_client.go`

### Q5. SSE三事件协议是做了什么？和普通流式输出有什么区别？

**标准回答口径：** 普通流式输出只推送文字片段。我们的SSE协议包含三个事件：`token`推送逐词流式文字、`sentence`推送完整句子、`final`推送结构化元数据（理解程度、续接位置、追问建议）。final事件是核心——它携带 `resume_page`、`resume_node_id`、`resume_sec` 三元组，学生端收到后可以直接精确续接，不需要额外的状态查询接口。

**关联模块索引：** M03架构·SSE协议 / D04 RAG+问答流 / `student.go` QAStream

### Q6. 如果Dify服务不可用，降级后的本地引擎功能完整吗？

**标准回答口径：** 本地Python引擎实现了课件解析和智能问答的完整功能子集，覆盖核心教学场景。区别在于本地引擎使用预设的工作流模板而非Dify的可视化工作流编排——功能一致但灵活度降低。解析和问答的核心逻辑在主引擎和降级引擎间保持一致，学生端不会感知到功能差异。

**关联模块索引：** M03架构·降级设计 / D06双引擎降级时序图 / `dify_client.go`

### Q7. 中间件选择PostgreSQL+Redis+MinIO是出于什么考虑？

**标准回答口径：** 三个中间件各管一类数据，没有重叠。PostgreSQL管结构化数据（课程、节点、用户、学习记录、会话）——需要事务和关联查询。Redis管缓存和队列——会话状态缓存24小时TTL、任务队列1000容量削峰。MinIO管文件——课件原文件、预览图、音频文件，S3兼容协议便于扩展。没有引入ES或图数据库是为了降低运维复杂度，目前的数据量级完全够用。

**关联模块索引：** M03架构·中间件 / D07中间件拓扑与数据流图

---

## A4.2 算法类

### Q8. 三阶段解析管线和直接调大模型一次生成有什么区别？

**标准回答口径：** 直接调大模型一次生成是"黑盒"——输出不可控、不可观测、不可局部修复。三阶段管线把任务拆成理解→建树→生成三个独立阶段：每阶段有独立的JSON Schema约束，大模型输出格式异常时触发自动修复；任一阶段失败只回退该阶段，不影响整体流程；阶段2产出的节点树可被问答和续接模块直接复用——一次解析，多次消费。这是工程化的思路，不是"调个接口就行"的思路。

**关联模块索引：** M04算法·解析管线 / D03三阶段解析管线图 / `generator.py`

### Q9. 苏格拉底式约束真的能防止AI直接给答案吗？

**标准回答口径：** 苏格拉底约束不是强依赖大模型"自觉遵守"，而是在提示词层面做了明确的行为约束——系统提示词明确写入"严禁直接给最终答案"，并制定了四步回答结构：确认卡点→提引导问题→给不超过两句提示→以问题结尾。同时在后端逻辑中，理解程度判定（`_predict_understanding`）会对学生问题的语义做分类——要求重讲时走补充讲解路径，正常提问时走启发式回答路径。约束层+判定层的双重保障，使AI偏离教学角色的概率大幅降低。

**关联模块索引：** M04算法·问答 / D04 RAG+问答流 / `qa.py` 系统提示词

### Q10. 三级级联上下文具体怎么工作的？为什么比RAG效果好？

**标准回答口径：** 传统RAG是全文检索后取Top-K片段，可能取到不相关的页面。我们的三级级联是定向检索：L1优先取当前教学节点的脚本片段——学生卡在哪页问的就优先用哪页的内容；L2如果L1信息不足，扩展到当前页面全部教学节点；L3最后才用全页原文拼接兜底。这样保证回答始终绑定在当前课程内容的局部范围内，而不是在整个知识库里"撒网"。这不是通用RAG，而是"课程知识约束RAG"。

**关联模块索引：** M04算法·问答 / D04 RAG+问答流 / `dialogue_support.go` buildNodeScopedContext

### Q11. 理解程度判定准确吗？怎么验证的？

**标准回答口径：** 理解程度判定采用LLM语义分析优先策略——将学生问题、课件内容和历史对话一起送入模型，让模型综合判断。模型不可用时降级到中文教育关键词匹配（如"听不懂""太难"等7个关键词匹配到则标记为"未理解"）。三态输出（none/partial/full）直接映射到续接策略：none回退重讲、partial补充讲解、full跳过继续。目前判定以LLM为主、关键词为兜底，在实测中准确率符合教学预期，后续可引入标注数据集做量化评估。

**关联模块索引：** M04算法·续接 / D05多态断点续学状态机 / `qa.py` _predict_understanding

### Q12. 断点续学和简单的"记住看到第几页"有什么区别？

**标准回答口径：** 传统续学只记住页码，恢复时从头开始播放整页。我们的续学精确到节点级别——`resumePage`/`resumeNodeID`/`resumeSec`三元组记录精确恢复位置：页面定位到第几页，节点精确到哪个教学单元，秒数精确到音频的哪个位置。更重要的是，续接策略不是固定的"从断点继续"，而是基于理解程度三态动态选择：没懂就回退重讲，部分懂就补充卡点，全懂就跳过。页级续学和节点级续学在精度上的差距，是"知道你在哪本书"和"知道你在第几页第几行"的区别。

**关联模块索引：** M04算法·续接 / D05多态断点续学状态机 / `student.go` QAStream final事件

### Q13. 学情掌握度模型的+0.08/-0.06/+0.02是怎么来的？

**标准回答口径：** 这三个参数是一个简化的间距重复模型——答对加分(+0.08)、答错减分(-0.06)、快速响应(≤5000ms)额外加分(+0.02)，掌握度强制约束在[0,1]区间。加减幅度不对称是因为知识掌握"建立难、丢失快"的教学认知规律——答对加得稍多但答错扣得更多。快速响应加分鼓励学生积极思考。初始值设为0.5（既不熟练也不完全不会），新知识点从中间位置开始随着练习量积累向1逼近或向0衰减。这个模型虽然不是复杂的AI模型，但在教学场景中比简单累计正确率提供了更精细的学情追踪。

**关联模块索引：** M04算法·学情 / D08学情多维加权指数 / `knowledge_map.go` UpdateMasteryScore

### Q14. 预置模板回退（Fallback）能保证质量吗？

**标准回答口径：** 预置模板不是"随便写一段文字占位"，而是基于课件标题和内容提纲的结构化生成——自动提取标题、关键点、组织结构，生成符合教学规范的讲稿骨架。模板质量虽然不如大模型生成的精细，但保证了三件事：格式正确、不跑题、可讲授。更重要的是，fallback机制保证了解析管线"在任何情况下都有输出"，不会因为模型调用失败而让老师拿到空页面。这是工程可靠性的体现，不是质量的妥协。

**关联模块索引：** M04算法·解析管线 / D03三阶段解析管线图 / `generator.py` _fallback_pipeline

---

## A4.3 工程类

### Q15. 测试覆盖怎么样？核心功能有自动化测试吗？

**标准回答口径：** 项目包含多层测试覆盖。算法层有Python单元测试（`ai_engine/tests/`），覆盖核心解析逻辑和推荐服务；后端Go有集成测试（`backend/tests/`），覆盖配置加载和MinIO对接；Dify工作流有独立测试脚本（`backend/test_workflow_*.go`），逐流程验证解析、生成、脚本产出。此外项目包含全维度测试验证报告（M06），记录功能测试、性能测试、集成测试和压力测试的全覆盖结果。当前测试重点覆盖核心算法逻辑和关键集成点，持续迭代中。

**关联模块索引：** M05工程开发设计规范 / M06测试验证报告 / `ai_engine/tests/`

### Q16. 并发能力怎么样？能支持多少学生同时使用？

**标准回答口径：** 系统并发瓶颈不在应用层而在AI引擎层。Go编排层采用goroutine模型，单机可轻松支撑千级SSE并发连接。任务调度引擎使用10个worker协程+1000容量的缓冲队列做削峰处理。真正的瓶颈是Dify或大模型API的QPS限制——我们通过降级引擎分流和请求排队来缓解。单台服务器在典型教学场景（每课件50学生同时在线、每人每节课提问3-5次）下完全够用，更大规模需要水平扩展AI引擎层。

**关联模块索引：** M03架构·任务调度 / M05工程开发设计规范 / D10任务调度状态机图

### Q17. 容错方案有哪些？单点故障怎么应对？

**标准回答口径：** 系统在每个关键路径都设计了容错。AI能力层：双引擎自动降级，Dify不可用时透明切换到本地引擎。课件解析：任一阶段失败回退到预置模板，单阶段失败不影响整体流程。上传管道：三级预览图降级（实时渲染→占位图→空状态），AI解析失败不阻塞上传。任务调度：3次指数退避重试（2/4/8分钟），超过最大重试次数标记为failed不再自动重试。中间件：会话状态同步到PostgreSQL和Redis双写，单缓存故障不影响核心数据。

**关联模块索引：** M03架构·降级与容错 / D06/D09/D10 / `course.go`

### Q18. 课件上传支持什么格式？大小有限制吗？

**标准回答口径：** 当前支持PDF和PPTX两种格式，前端上传时做格式校验（`validExts` 白名单过滤）。文件大小受MinIO和网关配置限制，默认没有硬编码上限——但超大文件（>200MB）解析耗时较长会影响用户体验，建议教师拆分上传。上传后自动触发异步解析流程：存储到MinIO→AI解析→生成CoursePage→生成TeachingNode，整个过程非阻塞，上传成功立即返回，教师无需等待解析完成。

**关联模块索引：** M05工程·上传管道 / D09课件上传容错管道 / `course.go` UploadCourse

### Q19. 代码质量和工程规范怎么样？

**标准回答口径：** 项目遵循Go标准工程布局（`cmd/`、`internal/`、`pkg/`三级目录规范），后端采用分层架构（handler→service→model），API路由统一注册在 `router/` 下。AI引擎遵循Python模块化组织（`generator.py`/`qa.py`/`schema.py` 功能分离）。代码包含完整的错误处理链，关键逻辑有日志输出（`pkg/logger`）。Go模块有完善的数据验证层（`pkg/learningevent/validate.go`），对用户输入做全面校验再进入业务逻辑。

**关联模块索引：** M05工程开发设计规范 / AGENTS.md / 项目源码结构

### Q20. 部署和运维方案是什么？

**标准回答口径：** 系统采用Docker容器化部署（`docker-compose.yml`），包含Go后端、AI引擎、PostgreSQL、Redis、MinIO等全部组件。单机部署即可运行完整系统。环境变量集中在 `.env` 文件中管理（AI模型配置、API密钥、数据库连接串）。后端健康检查端点 `/api/health` 可供监控系统探测。任务调度器自带自动重启机制，非核心服务故障不影响主业务流程。

**关联模块索引：** M07部署运维手册 / Dockerfile / docker-compose.yml

---

## A4.4 商业类

### Q21. 系统目前有实际客户在使用吗？

**标准回答口径：** 系统设计深度对接了泛雅教学平台的OpenAPI标准，已实现课程同步、用户同步、课件解析、进度追踪等集成接口。项目以服务外包竞赛为起点，但架构设计考虑了真实高校的教学场景——双引擎降级、泛雅对接适配、权限体系都是真实部署需求的产物。目前处于产品化验证阶段，核心能力已在测试环境中与泛雅平台完成联调验证。

**关联模块索引：** M08校企集成与落地实施 / OpenAPI对接接口 / `compat_openapi.go`

### Q22. 怎么收费？商业模式是什么？

**标准回答口径：** 系统面向高校教学场景设计，核心价值是帮助教师将传统课件快速转化为可交互、可问答、可续学的智能教学内容。目前的产品定位是服务外包竞赛项目，商业模式可灵活适配——学校自建部署（提供私有化部署包）、SaaS服务（按课件/学生数计费）、与泛雅等平台集成（按调用量结算）均可。具体收费模式需根据实际落地场景确定。

**关联模块索引：** M08校企集成与落地实施 / M10拓展功能与迭代规划

### Q23. 和竞品（超星学习通、雨课堂等）相比优势在哪？

**标准回答口径：** 超星和雨课堂解决的是"教学管理"问题——签到、作业、考试、资源分发。我们解决的是"教学内容智能化"问题——课件自动生成讲稿、AI苏格拉底式问答、节点级精确续学、学情多维加权分析。我们不是和他们竞争教学管理功能，而是在他们的教学管理之上提供更深度的教学智能化能力。与泛雅平台的OpenAPI对接设计，正好说明了我们定位为"能力补充"而非"平台替代"。

**关联模块索引：** M02总体方案概述 / M10拓展功能与迭代规划

### Q24. 和直接用ChatGPT/文心一言比，你们的价值在哪里？

**标准回答口径：** 通用大模型不知道你的课件内容——学生问问题，ChatGPT用它的训练知识回答，可能偏离课程大纲。我们的系统以"课程知识优先"为第一原则：课件解析后自动构建课程专属知识库，问答时通过三级级联上下文确保回答始终在课程范围内。还有一个关键差异：我们不是一次性的问答工具，而是"五环闭环"——可讲、可学、可问、可续、可评可优，问完之后理解程度判定自动影响续接策略，学情数据回传教师端反哺教学。这不是一个ChatGPT接口，而是一套完整教学场景的智能化引擎。

**关联模块索引：** M02总体方案 / M04算法全流程 / D01五环业务闭环 / D04 RAG+问答流

### Q25. 系统安全性怎么样？学生隐私数据怎么保护的？

**标准回答口径：** 安全性分三个层面保障。传输层：OpenAPI接口使用MD5签名校验（`enc`+`time`双重验证），内部API走JWT鉴权。存储层：敏感字段不在API响应中透出（如Answer模型中的正确答案标记为 `json:"-"`），音频内容使用SHA-256哈希做失效检测。平台层：支持多租户隔离（通过PlatformID/OrgCode区分），角色权限体系（teacher/student）限制数据访问范围。所有大模型调用不存储学生身份信息到第三方API。

**关联模块索引：** M09安全体系与合规保障 / OpenAPI签名机制 / JWT鉴权

### Q26. AI生成内容的质量怎么保证？出错了怎么办？

**标准回答口径：** 质量保证有三层。第一层是架构约束——三阶段管线每阶段有JSON Schema严格校验，输出格式异常自动修复。第二层是业务约束——问答强制苏格拉底式引导、溯源绑定页码，学生可验证答案来源。第三层是人机协同——教师端可以查看和修改讲稿脚本（`SaveScript`接口），AI生成内容经过教师审核后再发布给学生。每轮问答都会记录到数据库（`DialogueTurn`表），出问题可追溯可复盘。

**关联模块索引：** M04算法·质量保障 / M05工程·审核机制 / D03解析管线 / `teacher.go` SaveScript

### Q27. 系统能支持哪些教学场景？不只是PPT播放吗？

**标准回答口径：** 系统覆盖课前、课中、课后全场景。课前：教师上传课件→AI自动解析生成讲稿→教师审核修改→发布。课中：学生端播放学习+随时提问+AI苏格拉底式回答+理解程度判定+自动续接。课后：学情数据回传→教师查看班级共性卡点和个人薄弱知识点→调整教学策略。核心能力不是PPT播放，而是"课件内容智能化"——让静态课件变成可交互、可问答、可追踪、可优化的动态教学内容。

**关联模块索引：** M02总体方案 / D01五环业务闭环 / 教师端 / 学生端

### Q28. 项目的技术门槛高吗？学校能自己部署运维吗？

**标准回答口径：** 部署门槛较低——全部组件Docker容器化，单机 `docker-compose up` 即可启动完整系统。硬件要求不高：Go后端和Python引擎对CPU/内存要求适中，PostgreSQL/Redis/MinIO都可以使用Docker默认配置。运维复杂度主要体现在大模型API的Key管理和Dify工作流的维护上。我们已经将核心参数通过环境变量暴露（模型选择、API地址、超时配置），学校IT人员做简单的环境配置即可运维。后续可提供一键部署脚本和运维监控面板进一步降低门槛。

**关联模块索引：** M07部署运维手册 / M08校企集成与落地实施 / Dockerfile / docker-compose.yml
