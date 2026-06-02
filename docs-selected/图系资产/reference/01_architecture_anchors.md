# 图系引用素材：架构量化锚点索引

> 本文件记录所有架构类图表的量化数据来源。锚点数据集中管理，一处修改全局同步。
> 每项标注：锚点名 → 参数值 → 来源路径 → 使用图表

---

## 编排层锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| SSE首包响应时间 | <500ms | 编排层流式转发设计 | D02, D04 |
| 对话上下文窗口 | 4轮 | student.go buildDialogueContext | D04 |
| 编排层响应特性 | ms级（不做AI推理） | 架构设计决策 | D02 |
| 学习状态机主状态 | 4态（讲授/问答/判定/续接） | student.go QAStream | D02, D05 |

## AI能力层锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| AI推理典型延迟 | 3-5秒 | 通用模型调用经验 | D02 |
| 三阶段管线阶段数 | 3（理解/建树/生成） | generator.py | D03 |
| 生成温度参数 | 0.2 | generator.py GenerationConfig | D03 |
| 内容截断阈值 | 3500字符 | generator.py max_content_chars | D03 |
| 降级引擎地址 | 127.0.0.1:8000 | dify_client.go postJSON | D06 |
| 降级切换粒度 | 每次调用独立切换 | dify_client.go | D06 |
| 降级对用户透明 | 是（编排层封装） | dify_client.go + 架构设计 | D06 |
| 理解程度三态 | none / partial / full | student.go QAStream final事件 | D04, D05 |

## 数据层锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| 内容哈希失效机制 | SHA-256 | models.go AudioAsset | D07 |
| 预览图降级层数 | 3级 | course.go ensurePreviewFallback | D09 |
| 节点引用修复机制 | 双向对齐 | schema.py align_node_segment_mapping | D03 |

## 任务调度锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| 最大重试次数 | 3次 | task_scheduler.go | D10 |
| 重试退避策略 | 2^retryCount 分钟 | task_scheduler.go processTask | D10 |
| 并发协程数 | 10个 | task_scheduler.go Start | D10 |
| 任务队列容量 | 1000 | task_scheduler.go taskQueue | D10 |
| 调度扫描周期 | 1分钟 | task_scheduler.go dispatchLoop | D10 |
