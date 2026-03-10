import json
import os
import re
from dataclasses import dataclass
from typing import Any
from dotenv import load_dotenv

try:
    from .schema import build_node_script_schema
except ImportError:
    from schema import build_node_script_schema

load_dotenv()


@dataclass
class GenerationConfig:
    mode: str = "llm"  # 默认使用真实模型
    model: str = "qwen-plus"
    temperature: float = 0.2
    max_content_chars: int = 3500


class LessonGenerator:
    """根据解析后的页级文本，生成讲稿与思维导图结构。"""

    def __init__(self, config: GenerationConfig | None = None):
        self.config = config or GenerationConfig(
            mode=os.getenv("AI_GEN_MODE", "llm"),
            model=os.getenv("AI_MODEL", "qwen-plus"),
            temperature=float(os.getenv("AI_TEMPERATURE", "0.2")),
            max_content_chars=int(os.getenv("AI_MAX_CONTENT_CHARS", "3500")),
        )

    def generate(self, parsed_document: dict[str, Any]) -> dict[str, Any]:
        pages = parsed_document.get("parsed_pages", [])
        outputs: list[dict[str, Any]] = []

        for page_item in pages:
            page = page_item["page"]
            content = page_item["content"]
            generated = self._generate_llm(page, content)

            outputs.append(generated)

        return {
            "doc_id": parsed_document.get("doc_id"),
            "doc_name": parsed_document.get("doc_name"),
            "generator": {
                "mode": self.config.mode,
                "model": self.config.model,
                "temperature": self.config.temperature,
            },
            "lessons": outputs,
        }

    def _generate_llm(self, page: int, content: str) -> dict[str, Any]:
        normalized_content = self._prepare_content(content)
        outline = self._build_outline(normalized_content)

        system_prompt = (
            "你是高校教师讲稿生成助手。"
            "你的任务是把课件页面内容改写成可直接讲授的中文讲稿。"
            "必须准确覆盖原文关键信息，不要编造课件中不存在的结论、数据或公式。"
            "输出只能是 JSON 对象，且只包含 script 与 mindmap_markdown 两个字段。"
        )
        user_prompt = (
            f"页面: {page}\n"
            f"原始内容:\n{normalized_content}\n\n"
            f"提炼要点:\n{outline}\n\n"
            "请生成：\n"
            "1. script: 220-420 字，口语化、按老师上课讲述方式展开，包含开场点题、核心解释、必要举例和本页收束。\n"
            "2. mindmap_markdown: 使用 Markdown 列表，第一行为本页主题，下面给出 3-6 个要点分支。\n"
            "如果原始内容较少，优先解释已有信息，不要用空泛套话补字数。\n"
            "请返回 JSON: {\"script\": \"...\", \"mindmap_markdown\": \"...\"}"
        )

        try:
            client = self._build_llm_client()
            raw_content = self._request_llm_payload(client, system_prompt, user_prompt)
            data = self._extract_json(raw_content)
            return {
                "page": page,
                "script": self._clean_generated_text(data.get("script")) or self._fallback_script(page, normalized_content),
                "mindmap_markdown": self._clean_generated_text(data.get("mindmap_markdown")) or self._fallback_mindmap(page, normalized_content),
                "source_excerpt": normalized_content[:120].strip(),
            }
        except Exception as e:
            return self._fallback_generation(page, normalized_content, str(e))

    def generate_node_script(self, teaching_node: dict[str, Any], course_name: str | None = None) -> dict[str, Any]:
        node_id = str(teaching_node.get("node_id") or "node_unknown")
        title = str(teaching_node.get("title") or "未命名节点")
        core_points = teaching_node.get("core_points") or []
        examples = teaching_node.get("examples") or []
        confusions = teaching_node.get("common_confusions") or []
        summary = str(teaching_node.get("summary") or "")

        assembled_content = "\n".join(
            [
                title,
                summary,
                *[f"要点：{item}" for item in core_points],
                *[f"例子：{item}" for item in examples],
                *[f"易错点：{item}" for item in confusions],
            ]
        ).strip()

        fallback = self._fallback_node_script(node_id, title, summary, core_points, examples, confusions)

        try:
            client = self._build_llm_client()
            system_prompt = (
                "你是高校课程智能讲授编排助手。"
                "请根据讲授节点信息生成适合互动式课堂的讲授内容。"
                "输出严格 JSON，字段必须是 script、mindmap_markdown、interactive_questions、reteach_script、transition。"
            )
            user_prompt = (
                f"课程名: {course_name or '未提供'}\n"
                f"节点标题: {title}\n"
                f"节点摘要: {summary}\n"
                f"核心要点: {json.dumps(core_points, ensure_ascii=False)}\n"
                f"示例: {json.dumps(examples, ensure_ascii=False)}\n"
                f"易错点: {json.dumps(confusions, ensure_ascii=False)}\n"
                f"补充内容: {assembled_content}\n"
                "请输出适合老师直接讲授、并支持后续追问的结果。"
            )
            raw_content = self._request_llm_payload(client, system_prompt, user_prompt)
            data = self._extract_json(raw_content)
            return build_node_script_schema(
                node_id=node_id,
                title=title,
                script=self._clean_generated_text(data.get("script")) or fallback["script"],
                mindmap_markdown=self._clean_generated_text(data.get("mindmap_markdown")) or fallback["mindmap_markdown"],
                interactive_questions=data.get("interactive_questions") or fallback["interactive_questions"],
                reteach_script=self._clean_generated_text(data.get("reteach_script")) or fallback["reteach_script"],
                transition=self._clean_generated_text(data.get("transition")) or fallback["transition"],
            )
        except Exception:
            return fallback

    def _build_llm_client(self) -> Any:
        try:
            from openai import OpenAI
        except ImportError as error:
            raise RuntimeError("当前环境缺少 openai 依赖，请先安装 requirements.txt 后重试。") from error

        api_key = os.getenv("AI_API_KEY")
        if not api_key:
            raise RuntimeError("缺少环境变量 AI_API_KEY，无法调用大模型。")

        base_url = os.getenv("AI_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1")
        return OpenAI(api_key=api_key, base_url=base_url)

    def _request_llm_payload(self, client: Any, system_prompt: str, user_prompt: str) -> str:
        messages = [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_prompt},
        ]

        try:
            response = client.chat.completions.create(
                model=self.config.model,
                temperature=self.config.temperature,
                messages=messages,
                response_format={"type": "json_object"},
            )
        except Exception:
            response = client.chat.completions.create(
                model=self.config.model,
                temperature=self.config.temperature,
                messages=messages,
            )

        return response.choices[0].message.content or ""

    def _prepare_content(self, content: str) -> str:
        normalized = re.sub(r"\n{3,}", "\n\n", (content or "").strip())
        if len(normalized) <= self.config.max_content_chars:
            return normalized

        head = normalized[: int(self.config.max_content_chars * 0.7)].strip()
        tail = normalized[-int(self.config.max_content_chars * 0.2) :].strip()
        return f"{head}\n\n[中间内容已截断，保留首尾关键信息]\n\n{tail}".strip()

    def _build_outline(self, content: str) -> str:
        lines = [line.strip() for line in content.splitlines() if line.strip()]
        if not lines:
            return "- 本页缺少可用文本，需基于标题或截图补充信息"

        outline_lines: list[str] = []
        for line in lines[:6]:
            candidate = re.sub(r"^[\-•*\d.、\s]+", "", line)
            if candidate:
                outline_lines.append(f"- {candidate[:80]}")

        return "\n".join(outline_lines) if outline_lines else "- 本页内容较少，请围绕现有文本讲解"

    @staticmethod
    def _clean_generated_text(value: Any) -> str:
        return re.sub(r"\n{3,}", "\n\n", str(value or "").strip())

    def _fallback_generation(self, page: int, content: str, error_message: str) -> dict[str, Any]:
        return {
            "page": page,
            "script": self._fallback_script(page, content),
            "mindmap_markdown": self._fallback_mindmap(page, content),
            "source_excerpt": content[:120].strip(),
            "generation_error": error_message,
            "used_fallback": True,
        }

    def _fallback_script(self, page: int, content: str) -> str:
        lines = [line.strip() for line in content.splitlines() if line.strip()]
        title = lines[0] if lines else f"第{page}页内容"
        key_points = [re.sub(r"^[\-•*\d.、\s]+", "", line) for line in lines[1:4]]
        key_points = [item for item in key_points if item]

        segments = [f"这一页我们先抓住“{title}”这个主题。"]
        if key_points:
            segments.append("重点有：" + "；".join(key_points) + "。")
        segments.append("讲解时可以先说明概念，再结合课堂中的典型场景或例子帮助学生建立理解。")
        segments.append("最后要把本页结论和上一页、下一页的逻辑关系串起来，方便学生形成完整知识链。")
        return "".join(segments)

    def _fallback_mindmap(self, page: int, content: str) -> str:
        lines = [line.strip() for line in content.splitlines() if line.strip()]
        title = lines[0] if lines else f"第{page}页主题"
        children = lines[1:5] if len(lines) > 1 else ["核心概念", "关键解释", "课堂例子", "本页总结"]
        rows = [f"- {title}"]
        for item in children:
            cleaned = re.sub(r"^[\-•*\d.、\s]+", "", item).strip()
            if cleaned:
                rows.append(f"  - {cleaned}")
        return "\n".join(rows)

    def _fallback_node_script(
        self,
        node_id: str,
        title: str,
        summary: str,
        core_points: list[str],
        examples: list[str],
        confusions: list[str],
    ) -> dict[str, Any]:
        details = "；".join(core_points[:3]) if core_points else summary or title
        example_text = examples[0] if examples else f"可以结合 {title} 的典型应用场景帮助学生理解。"
        reteach_focus = confusions[0] if confusions else f"重新解释 {title} 的核心概念"
        script = (
            f"这一部分我们围绕“{title}”来讲。"
            f"先抓住几个关键点：{details}。"
            f"讲到这里时，可以用“{example_text}”作为课堂例子，让学生把抽象概念和真实情境对应起来。"
            "最后再回到本节点结论，帮助学生形成完整理解。"
        )
        mindmap = "\n".join([f"- {title}"] + [f"  - {item}" for item in (core_points[:4] or [summary or "核心内容"])])
        interactive_questions = [
            f"你觉得“{title}”最关键的点是什么？",
            f"如果把“{title}”放到真实场景里，它会怎么用？",
            f"你现在最不确定的是哪一步？",
        ]
        reteach_script = f"如果你还没完全听懂，我们就换个角度重讲，重点放在：{reteach_focus}。"
        transition = f"理解完“{title}”之后，我们就可以继续进入下一个知识节点。"
        return build_node_script_schema(node_id, title, script, mindmap, interactive_questions, reteach_script, transition)

    @staticmethod
    def _extract_json(text: str) -> dict[str, Any]:
        cleaned = text.strip()
        if cleaned.startswith("```"):
            cleaned = re.sub(r"^```(?:json)?", "", cleaned, count=1).strip()
            cleaned = re.sub(r"```$", "", cleaned).strip()
        start = cleaned.find("{")
        end = cleaned.rfind("}")
        if start == -1 or end == -1 or end <= start:
            raise ValueError("模型返回内容中未找到有效 JSON")
        return json.loads(cleaned[start : end + 1])
