# A12-v4-final — Design Spec

> Competition defense PPT for 服创大赛 National Finals
> Design principle: 品宣式标题 + 专业式内容

## I. Project Information

| Item | Value |
| ---- | ----- |
| **Project Name** | A12-v4-final — 泛雅AI互动智课系统 |
| **Canvas Format** | PPT 16:9 (1280x720) |
| **Page Count** | 17 |
| **Design Style** | B) General Consulting + Tech/Education — data clarity first, logical persuasion |
| **Target Audience** | National finals judges (professors, industry experts, enterprise representatives) |
| **Use Case** | 7-12 min defense presentation; 17 slides covering pain points, solution, architecture, algorithms, team |
| **Created Date** | 2026-06-02 |

---

## II. Canvas Specification

| Property | Value |
| -------- | ----- |
| **Format** | PPT 16:9 |
| **Dimensions** | 1280x720 |
| **viewBox** | `0 0 1280 720` |
| **Margins** | left/right 60px, top/bottom 50px |
| **Content Area** | 1160x620 |

---

## III. Visual Theme

### Theme Style
- **Style**: Competition defense — professional, tech-forward, deep blue contemporary
- **Theme**: Light theme with deep blue accent
- **Tone**: Professional, authoritative, technically rigorous

### Color Scheme
| Role | HEX | Purpose |
| ---- | --- | ------- |
| **Background** | `#FFFFFF` | Page background |
| **Secondary bg** | `#F0F4F8` | Card background, section background |
| **Primary** | `#1565C0` | Title decorations, key sections, icons — deep tech blue |
| **Accent** | `#FF6F00` | Data highlights, key information — warm orange contrast |
| **Secondary accent** | `#0D47A1` | Secondary emphasis, gradient transitions |
| **Body text** | `#1A1A2E` | Main body text |
| **Secondary text** | `#6B7280` | Captions, annotations |
| **Tertiary text** | `#9CA3AF` | Supplementary info, footers |
| **Border/divider** | `#E5E7EB` | Card borders, divider lines |
| **Success** | `#10B981` | Positive indicators (green) |
| **Warning** | `#EF4444` | Issue markers (red) |

### Gradient Scheme
```xml
<!-- Title gradient -->
<linearGradient id="titleGradient" x1="0%" y1="0%" x2="100%" y2="100%">
  <stop offset="0%" stop-color="#1565C0"/>
  <stop offset="100%" stop-color="#0D47A1"/>
</linearGradient>

<!-- Background decorative gradient -->
<radialGradient id="bgDecor" cx="80%" cy="20%" r="50%">
  <stop offset="0%" stop-color="#1565C0" stop-opacity="0.08"/>
  <stop offset="100%" stop-color="#1565C0" stop-opacity="0"/>
</radialGradient>

<!-- Card background gradient -->
<linearGradient id="cardBg" x1="0%" y1="0%" x2="0%" y2="100%">
  <stop offset="0%" stop-color="#FFFFFF"/>
  <stop offset="100%" stop-color="#F8FAFC"/>
</linearGradient>
```

---

## IV. Typography System

**Typography direction**: modern CJK sans — contrast pairing (KaiTi title + Microsoft YaHei body)

| Role | Chinese | English | Fallback tail |
| ---- | ------- | ------- | ------------- |
| **Title** | `KaiTi` | `Georgia` | `serif` |
| **Body** | `"Microsoft YaHei", "PingFang SC"` | `Arial` | `sans-serif` |
| **Emphasis** | `"Microsoft YaHei"` | `Georgia` | `serif` |
| **Code** | — | `Consolas, "Courier New"` | `monospace` |

**Per-role font stacks**:
- Title: `Georgia, KaiTi, serif` — Latin in Georgia (elegant serif), CJK in KaiTi (academic feel)
- Body: `"Microsoft YaHei", Arial, sans-serif` — CJK-primary, Latin secondary
- Emphasis: `Georgia, "Microsoft YaHei", serif`
- Code: `Consolas, "Courier New", monospace`

### Font Size Hierarchy
**Baseline**: Body font size = 20px

| Purpose | Ratio | @20px | Weight |
| ------- | ----- | ----- | ------ |
| Cover title | 3-4x | 60-80px | Bold |
| Page title | 1.8-2x | 36-40px | Bold |
| Subtitle | 1.2-1.4x | 24-28px | SemiBold |
| **Body** | **1x** | **20px** | Regular |
| Annotation | 0.75x | 15px | Regular |
| Page number | 0.55x | 11px | Regular |

### Formula Rendering Policy
`text-only` — No formula rendering needed for this deck.

---

## V. Layout Principles

### Page Structure
- **Header area**: 0-60px from top — page title or none (cover/chapter pages have no header bar)
- **Content area**: 70-650px — main content
- **Footer area**: 650-720px — page number, team name

### Layout Pattern Library
| Pattern | Usage |
| ------- | ----- |
| **Single column centered** | Covers (P01, P17), conclusion |
| **Asymmetric split (3:7)** | Content pages — left title/nav, right detailed content |
| **Three/four column cards** | Feature lists (P06, P07), team intro (P16) |
| **Top-bottom split** | Architecture pages (P09, P10) — diagram top, annotations bottom |
| **Matrix grid (2x2)** | Decision framework (P10) — 4 design decisions |
| **Full-bleed + floating text** | Cover (P01) |
| **Center-radiating** | Architecture overview (P09) — layered boxes radiating from center |

### Spacing Specification
| Element | Value |
| ------- | ----- |
| Safe margin | 60px |
| Content block gap | 32px |
| Card gap | 24px |
| Card padding | 24px |
| Card border radius | 12px |

---

## VI. Icon Usage Specification

### Source
- **Library**: `tabler-filled` — rounded filled style, professional for education/tech
- **Stroke width**: N/A (filled style)

### Recommended Icon List
| Purpose | Icon Path | Page |
| ------- | --------- | ---- |
| Teacher | `tabler-filled/user` | P04, P06 |
| Student | `tabler-filled/users` | P04, P07 |
| Document | `tabler-filled/file-description` | P03, P06 |
| Chat/QA | `tabler-filled/message` | P07 |
| Play | `tabler-filled/player-play` | P05 |
| Settings | `tabler-filled/adjustments` | P09 |
| Shield | `tabler-filled/shield` | P12 |
| Database | `tabler-filled/database` | P09 |
| Cloud | `tabler-filled/cloud` | P09 |
| Check | `tabler-filled/check` | P14 |
| Book | `tabler-filled/book` | P03 |
| Bulb | `tabler-filled/bulb` | P11 |
| Chart | `tabler-filled/chart-bar` | P06, P15 |
| Code | `tabler-filled/code-circle` | P11 |
| Download | `tabler-filled/download` | P06 |
| Search | `tabler-filled/search` | P04 |
| Star | `tabler-filled/star` | P17 |
| Team | `tabler-filled/users-group` | P16 |
| Building | `tabler-filled/building` | P15 |
| School | `tabler-filled/building-bridge-2` | P15 |

---

## VII. Visualization Reference List

No chart templates needed. All diagrams are custom SVG hand-drawn.

---

## VIII. Image Resource List

No external images needed. All visual elements (architecture diagrams, flow charts, cards) are hand-crafted SVG within each slide.

---

## IX. Content Outline

### Part 1: 项目总览与痛点 (P01-P04)

#### Slide 01 — Cover
- **Layout**: Full-bleed deep blue gradient background + centered title
- **Title**: 泛雅 AI 互动智课系统
- **Subtitle**: 以节点化内容底座为核心，实现课件解析→智能讲授→上下文问答→断点续接的教学闭环
- **Info**: 服创大赛 A12 · 智能计算赛道 / 命题企业 超星集团 / 参赛团队 DeepCore
- **Design**: Deep blue gradient (#1565C0→#0D47A1) background, white text, decorative geometric circles as background elements

#### Slide 02 — 方案总览
- **Layout**: Top-bottom split. Top: title bar. Bottom: 3-column cards
- **Title**: 一份课件，秒变会讲课、能答疑、可续接的智课
- **3 Cards**: ① 传即生成讲稿 — 三阶段解耦解析，从文本提取到跨页语义聚合 ② 问答秒级响应 — 本页优先检索+溯源生成，回答绑定课件页码 ③ 打断后自动续接 — 理解程度三态判定+节点级续接定位
- **Footer**: 五环闭环图示（可讲→可问→可续→可评→可优）

#### Slide 03 — 背景痛点
- **Layout**: Left title + right content block
- **Title**: 泛雅沉淀海量课件，却"可看不可讲"
- **Content**: 大模型推理、多模态理解、语音合成已具备落地条件，但单点调用无法支撑完整教学闭环——需要可编排、可降级、可续接的系统能力

#### Slide 04 — 痛点分析
- **Layout**: 3-column cards for 3 personas
- **Title**: 三类角色，三道断点
- **Cards**: ① 教师"学情看不见" ② 学生"进度易脱节, 问答不精准" ③ 平台"静态课件智化难"
- **Footer note**: 三个痛点根因是同一个——静态课件只有展示层，没有被系统理解成可定位的教学节点

### Part 2: 方案与双端 (P05-P08)

#### Slide 05 — 方案底座
- **Layout**: Asymmetric split — left 35% concept, right 65% layered architecture diagram (hand-drawn SVG)
- **Title**: 一套底座，两端分工
- **Left**: 业务编排层(Go) + AI能力中台(Dify+Python) + 数据层(PG/Redis/MinIO)
- **Right**: Three-layer architecture illustration with key annotation: "编排层不做AI推理，双引擎自动降级"

#### Slide 06 — 教师端
- **Layout**: 5-step horizontal flow + detail cards below
- **Title**: 教师端：把备课从"逐字撰写"变成"编辑确认"
- **5 Steps**: 智课生成(三阶段解耦) → 脚本编辑(节点粒度) → 发布授课 → 学情查看(多维加权) → 讲稿优化
- **Key tech**: 三阶段解耦管线标注 + 多维加权难度指数（四个维度）

#### Slide 07 — 学生端
- **Layout**: 3-column feature cards + flow diagram bottom
- **Title**: 学生端：边听讲边提问，问完自动接着学
- **3 Features**: ① 上下文问答(三级级联+苏格拉底式引导) ② 断点续学(三态判定) ③ 薄弱点强化(掌握度追踪)
- **Bottom**: SSE三事件流标注（token→sentence→final）

#### Slide 08 — 故事链
- **Layout**: Full-width horizontal flow chain — circular nodes connected by arrows
- **Title**: 一条故事链，把两端串成真实学习闭环
- **Chain**: 教师上传课件 → AI三阶段解析 → 学生播放学习 → 遇疑提问 → SSE流式回答 → 理解程度判定 → 三态续接 → 四维学情回传教师
- **Highlight**: Data flow arrows with "四维数据：提问记录·理解程度·卡点分布·停留时长"

### Part 3: 架构深度 (P09-P10)

#### Slide 09 — 架构组织方式
- **Layout**: Center-radiating layered architecture diagram (hand-drawn SVG)
- **Title**: 核心架构：编排与能力分离，三层协同
- **4 Layers (SVG hand-drawn boxes with arrows)**: 
  - 前端交互层(Vue 3) — 教师端/学生端/泛雅集成
  - 业务编排层(Go) — 学习状态机/SSE流式/上传管道/任务调度/OpenAPI适配
  - ★ 编排层承担状态控制与流转，不承担AI推理
  - AI能力层(Dify+Python) — 5大工作流 + 双引擎降级路径
  - 数据支撑层 — PostgreSQL/Redis/MinIO
- **Annotations**: Dual-engine fallback path arrow (Dify→失败→本地引擎)

#### Slide 10 — 设计决策
- **Layout**: 2x2 matrix — 4 decision cards
- **Title**: 架构设计决策：为什么这样组织
- **4 Cards**: ① 编排与能力分离(Why: AI推理延迟3-5s，编排层必须轻量) ② 能力可降级(Why: 教学场景不能中断) ③ 节点化统一语言(Why: 页不是最小单位，教学节点才是) ④ 课程知识优先(Why: 通用大模型回答易偏离课件)

### Part 4: 算法深度 (P11-P14)

#### Slide 11 — 节点化底座
- **Layout**: Top-bottom — top: concept, bottom: dual-structure illustration
- **Title**: 先节点化，再问答，最后续接
- **Content**: 三阶段解耦：拆页渲染保版面→多模态识别保语义→跨页聚合保脉络。节点双重结构：知识图谱定义（前置依赖、难度系数）+ 脚本片段映射（一段文本可跨节点复用）

#### Slide 12 — 解析链
- **Layout**: Left-right flow diagram (SVG pipeline)
- **Title**: 课件解析链：把静态文件变成可讲授节点
- **Pipeline SVG**: 输入文件 → Stage1 Markdown理解 → Stage2 节点树构建 → Stage3 逐节点脚本生成 → 结构化讲稿
- **Annotations**: VLM渐进增强标注 + 标准化校验+引用一致性修复标注 + Schema约束标注

#### Slide 13 — 问答续接链
- **Layout**: Left-right dual-flow: left=问答流, right=续接流
- **Title**: 问答续接链：回答问题，但不打断学习主线
- **Left (QA flow)**: 学生提问 → 三级级联上下文 → 课程知识检索 → 苏格拉底式回答(禁止直接给答案) → SSE三事件推送
- **Right (Resume flow)**: 理解程度判定(LLM优先→关键词降级) → 三态决策(none→重讲/partial→补充/full→跳过) → 节点级恢复

#### Slide 14 — 模块覆盖表
- **Layout**: Full-width comparison table
- **Title**: 三大模块全覆盖
- **Table columns**: 赛题模块 | 项目实现(机制描述) | 对应算法
- **3 Rows**: 智课生成 / 实时问答 / 进度续接 — each with detailed mechanism implementation

### Part 5: 落地与团队 (P15-P17)

#### Slide 15 — 项目进展
- **Layout**: Left: timeline, Right: bullet achievements
- **Title**: 项目进展清晰，已进入学校推广阶段
- **Left timeline SVG**: 需求→设计→开发→测试→部署→推广
- **Right achievements**: 教师端与学生端功能链路已开发完成，通过学校真实教学环境验证；从课件上传到节点化解析、从上下文问答到三态续接、从学情数据回传到教师端，形成端到端闭环

#### Slide 16 — 团队
- **Layout**: 5-column card layout — one card per team member
- **Title**: DeepCore 团队 · 五人在线，全员攻坚
- **5 Cards**: Each with role title, responsibility description, and defense capability
  - 项目管理 (进度管控+风险管控)
  - 产品与架构 (需求分析+方案设计)
  - 后端开发 (Go+AI中台)
  - 前端开发 (Vue 3+交互)
  - 测试与交付 (质量保障+部署运维)
- **Bottom**: "全员开发保障技术深度……主责锚点保障责任清晰……"

#### Slide 17 — 谢谢观看
- **Layout**: Centered thank you with tech anchor statement
- **Title**: 谢谢观看
- **Subtitle**: "节点化底座 · 编排与能力分离 · 教学闭环"

---

## X. Speaker Notes Requirements

One speaker note file per page, saved to `notes/`:
- **Filename**: matches SVG name (e.g., `01_cover.md`)
- **Content**: script key points, timing cues, transition phrases
- **Reference**: detailed defense speaking points in V3.0 modification document

---

## XI. Technical Constraints Reminder

### SVG Generation Must Follow:
1. viewBox: `0 0 1280 720`
2. Background uses `<rect>` elements
3. Text wrapping uses `<tspan>` (`<foreignObject>` FORBIDDEN)
4. Transparency uses `fill-opacity` / `stroke-opacity`; `rgba()` FORBIDDEN
5. FORBIDDEN: `mask`, `<style>`, `class`, `foreignObject`
6. FORBIDDEN: `textPath`, `animate*`, `script`
7. PPT-safe font stacks ending with pre-installed fonts
8. Icons use `<use data-icon="tabler-filled/icon-name" .../>`
9. All diagrams hand-drawn as SVG shapes (rect, circle, path, line, text) — NO external image dependencies
10. Color values from spec_lock.md only — no invented colors

### Style Consistency Rules:
- Keep consistent with original competition PPT style (deep blue, professional)
- Architecture diagrams: layered boxes with connecting arrows
- Flow charts: horizontal/vertical node→arrow→node
- Cards: rounded rect borders with subtle shadows
- No emoji in formal content areas (only in decorative elements if needed)
