# 图系引用素材：算法量化锚点索引

> 本文件记录所有算法类图表的量化数据来源。锚点数据集中管理，一处修改全局同步。
> 每项标注：锚点名 → 参数值 → 来源路径 → 使用图表

---

## 掌握度模型锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| 掌握度初始值 | 0.5 | knowledge_map.go UpdateMasteryScore | D08 |
| 答对增量 | +0.08 | knowledge_map.go | D08 |
| 答错减量 | -0.06 | knowledge_map.go | D08 |
| 快速响应加分阈值 | ≤5000ms | knowledge_map.go | D08 |
| 快速响应加分值 | +0.02 | knowledge_map.go | D08 |
| 掌握度取值范围 | [0, 1] | knowledge_map.go clampScore | D08 |
| 学情维度数 | 4（提问频次/需重讲/会话数/停留时长） | teacher.go roundTeacherCardIndex | D08 |

## 解析管线锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| 解析管线阶段数 | 3 | generator.py generate_from_markdown | D03 |
| Stage1输出结构 | normalized_markdown + key_points | schema.py build_stage1_markdown_schema | D03 |
| Stage2节点属性 | node_id/title/summary/source_span/prerequisites | schema.py build_stage2_node_tree_schema | D03 |
| Stage3脚本段类型 | 6种(opening/explanation/example/interaction/transition/summary) | schema.py build_enhanced_script_schema | D03 |
| Schema校验模式 | additionalProperties:false（严格模式） | schema.py 所有Schema定义 | D03 |
| LLM失败回退 | 预置模板fallback | generator.py _fallback_pipeline | D03 |
| 节点引用修复 | normalize_stage3_scripts自动映射 | generator.py + schema.py | D03 |
| 引用一致性修复 | align_node_segment_mapping双向对齐 | schema.py | D03 |

## 问答系统锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| 上下文级联层数 | 3（节点→页面→全页） | dialogue_support.go buildNodeScopedContext | D04 |
| L1上下文范围 | 目标节点ScriptSegments | buildNodeScopedContext | D04 |
| L2上下文范围 | 当前页全部TeachingNode | buildNodeScopedContext | D04 |
| L3上下文范围 | CoursePage.SourceText拼接 | buildNodeScopedContext | D04 |
| 苏格拉底约束 | 严禁直接给答案 | qa.py 系统提示词 | D04 |
| LLM理解度判定 | 综合语义分析 | QAStream askQuestionWithFallback | D04, D05 |
| 关键词降级判定 | 中文教育关键词匹配 | qa.py RETEACH_KEYWORDS | D04, D05 |
| SSE事件数 | 3（token/sentence/final） | student.go QAStream writeEvent | D04 |
| final事件字段 | session_id/understanding_level/source_page/resume_page/resume_node_id/resume_sec/follow_up | QAStream final事件 | D04, D05 |

## 续接算法锚点

| 锚点名 | 参数值 | 来源 | 使用图表 |
|--------|:------:|------|:--------:|
| 断点续接精度 | 节点级（node_id） | QAStream final事件resume字段 | D05 |
| 续接状态三元组 | resumePage/resumeNodeID/resumeSec | QAStream + audio_support.go | D05 |
| 续接模式数 | 3（续讲/补充/重讲） | 理解程度三态映射 | D05 |
| none→续接策略 | 回退重讲（ReteachScript） | resolveResumeNodeIDByCourse | D05 |
| partial→续接策略 | 补充讲解（ScriptText定位卡点） | resolveResumeNodeIDByCourse | D05 |
| full→续接策略 | 跳过继续（下一节点） | resolveResumeNodeIDByCourse | D05 |
