export const meta = {
  name: 'finalize-booklet',
  description: 'Generate appendix, global index, and PPT compliance checklist',
  phases: [
    { title: 'Appendices', detail: 'A1-A4 appendix content' },
    { title: 'Global Index', detail: 'Judge Q&A cross-reference index' },
    { title: 'PPT Audit', detail: 'Compliance checklist vs PPT audit' },
    { title: 'Summary', detail: 'Cross-reference and finalize' },
  ],
}

// ===== A1-A2: 资质+佐证 =====
phase('Appendices')
var a12Agent = await agent(
  '你负责生成小册子附录A1「资质证明材料」和A2「佐证材料」的全部正文。\n' +
  '## 内容结构\n\n' +
  '### A1：资质证明材料\n' +
  '写清楚：\n' +
  '1. 软件著作权申报规划（基于系统自主开发代码量说明）\n' +
  '2. 开源组件合规清单汇总（项目使用的开源组件及许可证说明）\n' +
  '3. 系统自主代码声明\n' +
  '注：无实际软著/专利证书则标注"申报规划中"\n\n' +
  '### A2：佐证材料\n' +
  '写清楚：\n' +
  '1. 试点落地记录索引（如有试点，说明试点情况）\n' +
  '2. 需求对接记录说明\n' +
  '3. 技术验证原始数据索引（测试报告、压测数据等）\n' +
  '注：无实际佐证照片/回执则标注"待补充"\n\n' +
  '## 风格：索引清单风、清晰可查\n' +
  '## 输出格式：{"content": "完整markdown含A1和A2", "summary": "摘要", "sub_items_completed": ["A1","A2"]}',
  { label: 'Appendices-A1A2', phase: 'Appendices',
    schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } }
)

// ===== A3-A4: 技术附件+答辩问答库 =====
phase('Appendices')
var a34Agent = await agent(
  '你负责生成小册子附录A3「技术附件」和A4「答辩备用问答库」的全部正文。\n' +
  '## 内容结构\n\n' +
  '### A3：技术附件\n' +
  '写清楚：\n' +
  '1. 完整图系索引 — 引用 docs-selected/图系资产/README.md 中的所有图表，标注每张图对应模块和用途\n' +
  '2. 关键Schema定义摘要 — 三阶段解析管线的JSON Schema结构说明\n' +
  '3. 接口规范摘要 — 系统OpenAPI的整体分类和约定\n' +
  '4. 量化参数参考表 — 从源码中提取的关键参数速查\n\n' +
  '### A4：答辩备用问答库\n' +
  '写清楚：按维度分类的评委高频质疑问题+标准应答口径：\n' +
  '1. 架构类（why 4 layers? why Go? why not all-in-one? 降级怎么保证？）\n' +
  '2. 算法类（三阶段管线与直接调大模型有什么区别？苏格拉底约束真实有效吗？）\n' +
  '3. 工程类（测试覆盖怎么样？并发能力？容错方案？）\n' +
  '4. 商业类（真的有客户吗？怎么收费？和竞品比优势在哪？）\n' +
  '每个问题需给出：问题 → 标准回答口径(3句话以内) → 关联模块索引\n\n' +
  '## 风格：A3索引风 / A4实战问答风\n' +
  '## 输出格式：{"content": "完整markdown含A3和A4", "summary": "摘要", "sub_items_completed": ["A3","A4"]}',
  { label: 'Appendices-A3A4', phase: 'Appendices',
    schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } }
)

// ===== 全局交叉索引 =====
phase('Global Index')
var zxAgent = await agent(
  '你负责生成小册子的「全局交叉索引ZX」章节。\n' +
  '这是一个评委提问→模块定位的快速检索表，让评委拿到手册就能找到对应内容。\n\n' +
  '## 内容结构\n\n' +
  '### ZX-1：按提问维度索引表\n' +
  '一个大型表格，列：评委提问 → 对应模块 → 对应PPT页 → 图系编号 → 一句话答案\n\n' +
  '需要覆盖的提问维度（至少40个问题）：\n' +
  '- 架构类：架构怎么分层的？编排层为什么用Go？Dify挂了怎么办？中间件怎么部署的？并发能力？水平扩展？\n' +
  '- 算法类：三阶段解析做什么的？和通用大模型有什么区别？苏格拉底约束是什么意思？怎么判断学生理解程度？续接怎么做到的？掌握度怎么计算？\n' +
  '- 工程类：项目用了什么技术栈？代码规范？接口怎么设计的？数据库怎么设计的？异常处理？\n' +
  '- 测试类：测试了哪些？性能指标？降级成功率？\n' +
  '- 需求类：解决了什么痛点？谁的目标用户？有什么商业价值？怎么收费？\n' +
  '- 落地类：有真实用户吗？集成到什么平台了？部署需要什么环境？\n' +
  '- 安全类：数据安全？知识产权？开源合规？\n' +
  '- 未来类：后续迭代方向？什么时候能商用？\n\n' +
  '格式示例：\n' +
  '| 评委提问 | 对应模块 | PPT页 | 图系 | 一句话答案 |\n' +
  '|---------|---------|:----:|:---:|----------|\n' +
  '| 架构为什么这样分层？ | M03-01 | P09-P10 | D02 | 编排与能力分离——编排层不做AI推理，保证响应延迟可控 |\n\n' +
  '### ZX-2：模块快速定位索引\n' +
  '按模块编号列出每个模块的核心内容和页码索引\n\n' +
  '## 风格：实用检索风、表格为主\n' +
  '## 输出格式：{"content": "完整markdown含ZX-1和ZX-2", "summary": "摘要", "sub_items_completed": ["ZX-1","ZX-2"]}',
  { label: 'Global-Index', phase: 'Global Index',
    schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } }
)

// ===== PPT合规审核 =====
phase('PPT Audit')
var pptAgent = await agent(
  '你负责生成「PPT合规审核修改清单」。\n' +
  '对照自查表，逐项检查17页PPT是否存在缺失或表述不规范的情况，输出整改清单。\n\n' +
  '## 必须读取的文件\n' +
  '1. docs-selected/中国大学生服务外包创新创业大赛国赛答辩项目全维度合规自查审核表.md\n' +
  '2. docs-selected/PPT完整最终修改稿.md\n' +
  '3. docs-selected/PPT正式版修改意见文档V2.0.md\n\n' +
  '## 内容结构\n\n' +
  '### 逐维度检查结果表\n' +
  '对自查表全部9个维度（市场价值、痛点量化、架构设计、交付物、团队分工、创新点、赛题匹配、合规、其他），逐一标注：\n' +
  '- 自查条目编号和内容\n' +
  '- 当前PPT覆盖状态（已覆盖/部分覆盖/未覆盖）\n' +
  '- 对应PPT页码\n' +
  '- 缺失内容说明\n' +
  '- 整改建议（包括在小册子中补充的内容）\n\n' +
  '### PPT页面逐页修改建议\n' +
  '对17页PPT，给出基于自查表的最终修改建议。\n\n' +
  '## 风格：清单风、可执行\n' +
  '## 输出格式：{"content": "完整markdown含审核表和修改建议", "summary": "摘要", "sub_items_completed": ["audit-table","page-by-page"]}',
  { label: 'PPT-Audit', phase: 'PPT Audit',
    schema: { type: 'object', properties: { content: { type: 'string' }, summary: { type: 'string' }, sub_items_completed: { type: 'array', items: { type: 'string' } } }, required: ['content', 'summary', 'sub_items_completed'] } }
)

// ===== Summary =====
phase('Summary')
log('All final agents completed.')
var completed = []
if (a12Agent) { completed.push('Appendices-A1A2') }
if (a34Agent) { completed.push('Appendices-A3A4') }
if (zxAgent) { completed.push('Global-Index-ZX') }
if (pptAgent) { completed.push('PPT-Audit') }
log('Final batch completed: ' + completed.join(', '))

return {
  modules_completed: completed,
  a12_summary: a12Agent ? a12Agent.summary : null,
  a34_summary: a34Agent ? a34Agent.summary : null,
  zx_summary: zxAgent ? zxAgent.summary : null,
  ppt_summary: pptAgent ? pptAgent.summary : null,
}
