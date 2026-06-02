export const meta = {
  name: 'generate-booklet-batch2',
  description: 'Generate remaining 6 booklet modules in parallel',
  phases: [
    { title: 'M02 Solution', detail: 'Overall solution and design' },
    { title: 'M05 Engineering', detail: 'Development engineering standards' },
    { title: 'M07 Deployment', detail: 'Deploy, operations, user manual' },
    { title: 'M08 Integration', detail: 'School integration and pilot deployment' },
    { title: 'M09 Security', detail: 'Security and compliance' },
    { title: 'M10 Roadmap', detail: 'Extensions and iteration planning' },
    { title: 'Save All', detail: 'Cross-reference and save outputs' },
  ],
}

var OUTPUT_DIR = '/mnt/d/Desktop/FuWuWaiBao/docs-selected'

// Helper: simple string builder to avoid template issues
function prompt(text) { return text; }

// ===== M02 =====
phase('M02 Solution')
var m02Prompt = prompt(
  '你负责生成小册子M02模块「总体方案与概要设计」的全部正文。\n' +
  '## 必须读取的源文件\n' +
  '1. docs-selected/2-解决方案与业务功能/启智云双端解决方案PPT大纲.md\n' +
  '2. docs-selected/2-解决方案与业务功能/教师学生联合业务流程图.md\n' +
  '3. docs-selected/PPT完整最终修改稿.md\n\n' +
  '## 内容结构\n' +
  '### 02-01：总体业务方案\n' +
  '写清楚：双端产品设计思路、五环闭环业务流程、产品定位\n' +
  '### 02-02：技术选型论证\n' +
  '写清楚：Go/Dify/Python/PostgreSQL/Redis/MinIO各技术为什么选、不选什么、项目中怎么用\n' +
  '### 02-03：产品原型说明\n' +
  '写清楚：教师端和学生端核心页面交互逻辑\n\n' +
  '## 风格：正式、商务化、方案设计文档风\n' +
  '## 输出格式：{"content": "完整markdown", "summary": "摘要", "sub_items_completed": ["02-01","02-02","02-03"]}'
)
var m02Agent = await agent(m02Prompt, { label: 'M02-Solution', phase: 'M02 Solution', schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } })

// ===== M05 =====
phase('M05 Engineering')
var m05Prompt = prompt(
  '你负责生成小册子M05模块「工程开发设计规范」的全部正文。\n' +
  '## 必须读取的源文件\n' +
  '1. backend/api/main.go 和 backend/internal/ 目录结构\n' +
  '2. ai_engine/ 目录结构\n' +
  '3. backend/internal/model/models.go — 数据模型\n' +
  '4. backend/internal/service/task_scheduler.go — 调度器\n' +
  '5. docs-selected/图系资产/reference/01_architecture_anchors.md\n\n' +
  '## 内容结构\n' +
  '### 05-01：项目工程目录规范 — 前后端分离结构、模块拆分\n' +
  '### 05-02：代码管控规范 — Git分支、commit规范\n' +
  '### 05-03：接口文档 — 核心API分类（按功能描述，不列出具体路径）\n' +
  '### 05-04：数据库概要设计 — 核心表关系、索引思路\n' +
  '### 05-05：统一异常处理与日志 — 异常分层、日志分级\n\n' +
  '## 风格：工程化、可交付风\n' +
  '## 输出格式：{"content": "完整markdown", "summary": "摘要", "sub_items_completed": ["05-01","05-02","05-03","05-04","05-05"]}'
)
var m05Agent = await agent(m05Prompt, { label: 'M05-Engineering', phase: 'M05 Engineering', schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } })

// ===== M07 =====
phase('M07 Deployment')
var m07Prompt = prompt(
  '你负责生成小册子M07模块「部署运维与使用手册」的全部正文。\n' +
  '## 必须读取的源文件\n' +
  '1. 项目中Dockerfile、docker-compose等部署配置文件（如有）\n' +
  '2. backend/internal/service/course.go — 课件上传+三级预览降级\n' +
  '3. backend/internal/service/dify_client.go — 引擎配置\n\n' +
  '## 内容结构\n' +
  '### 07-01：部署手册 — 环境依赖、部署步骤、容器化方案\n' +
  '### 07-02：运维文档 — 日常巡检、故障排查（Dify降级/DB异常/MinIO异常）、日志管理\n' +
  '### 07-03：用户手册（教师端） — 课件管理→脚本编辑→发布→学情查看\n' +
  '### 07-04：用户手册（学生端） — 课程学习→提问→续学→掌握度\n\n' +
  '## 风格：操作手册风、步骤清晰\n' +
  '## 输出格式：{"content": "完整markdown", "summary": "摘要", "sub_items_completed": ["07-01","07-02","07-03","07-04"]}'
)
var m07Agent = await agent(m07Prompt, { label: 'M07-Deployment', phase: 'M07 Deployment', schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } })

// ===== M08 =====
phase('M08 Integration')
var m08Prompt = prompt(
  '你负责生成小册子M08模块「校企集成与落地实施」的全部正文。\n' +
  '## 必须读取的源文件\n' +
  '1. backend/internal/handler/compat_openapi.go — 泛雅OpenAPI集成\n' +
  '2. backend/internal/handler/platform_v1.go — 平台V1接口\n' +
  '3. backend/internal/handler/compat_common.go — 鉴权机制\n' +
  '4. backend/internal/handler/compat_student.go — 学生端兼容\n' +
  '5. docs-selected/4-团队分工项目管理与未来展望/团队分工与进度推广.md\n\n' +
  '## 内容结构\n' +
  '### 08-01：泛雅平台集成方案 — OpenAPI对接、鉴权、键名兼容\n' +
  '### 08-02：校园试点落地文档 — 试点院校、规模（基于文档，不虚构）\n' +
  '### 08-03：师生试用反馈 — 效率提升、体验反馈\n' +
  '### 08-04：落地成果佐证索引 — 可提供的证明材料清单\n\n' +
  '## 风格：正式、落地化、举证化\n' +
  '## 输出格式：{"content": "完整markdown", "summary": "摘要", "sub_items_completed": ["08-01","08-02","08-03","08-04"]}'
)
var m08Agent = await agent(m08Prompt, { label: 'M08-Integration', phase: 'M08 Integration', schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } })

// ===== M09 =====
phase('M09 Security')
var m09Prompt = prompt(
  '你负责生成小册子M09模块「安全体系与合规保障」的全部正文。\n' +
  '## 必须读取的源文件\n' +
  '1. backend/internal/handler/crypto.go — AES-256-GCM加密\n' +
  '2. backend/internal/handler/compat_common.go — Token鉴权、签名\n' +
  '3. backend/internal/model/models.go — 数据模型\n' +
  '4. docs-selected/中国大学生服务外包创新创业大赛国赛答辩项目全维度合规自查审核表.md — 维度八\n\n' +
  '## 内容结构\n' +
  '### 09-01：数据安全方案 — HTTPS、AES-256-GCM字段加密、Token+Session、签名校验\n' +
  '### 09-02：知识产权说明 — 软著/专利进展、代码自主率\n' +
  '### 09-03：开源组件合规审计 — 主要开源组件及License说明\n\n' +
  '## 风格：正式、合规风、可审计\n' +
  '## 输出格式：{"content": "完整markdown", "summary": "摘要", "sub_items_completed": ["09-01","09-02","09-03"]}'
)
var m09Agent = await agent(m09Prompt, { label: 'M09-Security', phase: 'M09 Security', schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } })

// ===== M10 =====
phase('M10 Roadmap')
var m10Prompt = prompt(
  '你负责生成小册子M10模块「拓展功能与迭代规划」的全部正文。\n' +
  '## 必须读取的源文件\n' +
  '1. ai_engine/ — 查看已有AI能力和推荐模块\n' +
  '2. backend/internal/handler/compat_student.go — 学生端拓展功能（笔记/收藏/错题/复习）\n' +
  '3. backend/internal/service/task_scheduler.go — 定时调度\n' +
  '4. docs-selected/PPT完整最终修改稿.md — P16路线图\n' +
  '5. docs-selected/4-团队分工项目管理与未来展望/团队分工与进度推广.md\n\n' +
  '## 内容结构\n' +
  '### 10-01：当前已拓展功能 — TTS语音合成、笔记收藏、错题重练、复习计划、薄弱点检测\n' +
  '### 10-02：版本迭代规划 — 短期/中期/长期规划\n' +
  '### 10-03：多场景复用方案 — K12/职业教育/企业内训适配\n\n' +
  '## 风格：务实、可落地、不虚构\n' +
  '## 输出格式：{"content": "完整markdown", "summary": "摘要", "sub_items_completed": ["10-01","10-02","10-03"]}'
)
var m10Agent = await agent(m10Prompt, { label: 'M10-Roadmap', phase: 'M10 Roadmap', schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } })

// ===== Save =====
phase('Save All')
var completed = []
if (m02Agent) { completed.push('M02-Solution') }
if (m05Agent) { completed.push('M05-Engineering') }
if (m07Agent) { completed.push('M07-Deployment') }
if (m08Agent) { completed.push('M08-Integration') }
if (m09Agent) { completed.push('M09-Security') }
if (m10Agent) { completed.push('M10-Roadmap') }

log('Batch 2 complete. Generated: ' + completed.join(', '))

return {
  modules_completed: completed,
  m02_summary: m02Agent ? m02Agent.summary : null,
  m05_summary: m05Agent ? m05Agent.summary : null,
  m07_summary: m07Agent ? m07Agent.summary : null,
  m08_summary: m08Agent ? m08Agent.summary : null,
  m09_summary: m09Agent ? m09Agent.summary : null,
  m10_summary: m10Agent ? m10Agent.summary : null,
}
