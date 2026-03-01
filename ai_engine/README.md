# AI Engine - Parser 模块

## 当前能力

- 支持 PDF 按页解析
- 支持 PPTX 按页（幻灯片）解析
- 返回统一结构化 JSON（含 `doc_id`、`total_pages`、`parsed_pages`、`stats`）
- 支持命令行输出到控制台或写入 JSON 文件

## 目录

- `parser.py`：解析入口与命令行运行器
- `schema.py`：统一结构化输出协议

## 快速运行

```powershell
conda activate fuww_ai
python ai_engine\parser.py <你的PDF或PPTX路径>
```

输出到文件：

```powershell
conda activate fuww_ai
python ai_engine\parser.py <你的PDF或PPTX路径> --out ai_engine\result.json
```

## 输出协议（示例字段）

- `doc_id`: 文档唯一 ID
- `doc_name`: 文件名
- `doc_type`: `pdf` 或 `pptx`
- `total_pages`: 总页数
- `parsed_pages`: 页级内容数组（`page`, `content`, `content_length`）
- `stats`: 统计信息（`non_empty_pages`, `empty_pages`, `elapsed_ms`）

## 下一步

- 接入大模型生成：讲稿 + 思维导图
- 增加问答溯源（按页上下文回答）
- 加入评测指标与鲁棒性优化
