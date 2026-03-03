import json
import os
from dataclasses import dataclass
from typing import Any
from dotenv import load_dotenv

load_dotenv()


@dataclass
class GenerationConfig:
    mode: str = "llm"  # 默认使用真实模型
    model: str = "qwen-plus"
    temperature: float = 0.2


class LessonGenerator:
    """根据解析后的页级文本，生成讲稿与思维导图结构。"""

    def __init__(self, config: GenerationConfig | None = None):
        self.config = config or GenerationConfig(
            mode=os.getenv("AI_GEN_MODE", "llm"),
            model=os.getenv("AI_MODEL", "qwen-plus"),
            temperature=float(os.getenv("AI_TEMPERATURE", "0.2")),
        )

    def generate(self, parsed_document: dict[str, Any]) -> dict[str, Any]:
        pages = parsed_document.get("parsed_pages", [])
        outputs: list[dict[str, Any]] = []

        for page_item in pages:
            page = page_item["page"]
            content = page_item["content"]
            # 始终使用 LLM 生成模式，移除 Mock 逻辑调用
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
        try:
            from openai import OpenAI
        except ImportError as error:
            raise RuntimeError(
                "当前环境缺少 openai 依赖，请先安装 requirements.txt 后重试。"
            ) from error

        api_key = os.getenv("AI_API_KEY")
        if not api_key:
            raise RuntimeError("缺少环境变量 AI_API_KEY，无法调用大模型。")

        base_url = os.getenv("AI_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1")
        client = OpenAI(api_key=api_key, base_url=base_url)

        system_prompt = (
            "你是高校课程助教。"
            "请根据给定页面内容，产出两个字段：script（口语化讲稿）与 mindmap_markdown（思维导图Markdown）。"
            "输出必须是严格 JSON，不要额外解释。"
        )
        user_prompt = (
            f"页面: {page}\n"
            f"内容: {content}\n"
            "请返回 JSON: {\"script\": \"...\", \"mindmap_markdown\": \"...\"}"
        )

        response = client.chat.completions.create(
            model=self.config.model,
            temperature=self.config.temperature,
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_prompt},
            ],
            response_format={"type": "json_object"}
        )

        try:
            raw_content = response.choices[0].message.content
            data = json.loads(raw_content)
            return {
                "page": page,
                "script": data.get("script", "无法生成讲稿"),
                "mindmap_markdown": data.get("mindmap_markdown", ""),
                "source_excerpt": content[:120].strip()
            }
        except Exception as e:
            return {
                "page": page,
                "script": f"生成失败: {str(e)}",
                "mindmap_markdown": "",
                "source_excerpt": content[:120].strip()
            }

    @staticmethod
    def _extract_json(text: str) -> dict[str, Any]:
        cleaned = text.strip()
        if cleaned.startswith("```"):
            cleaned = cleaned.strip("`")
            cleaned = cleaned.replace("json\n", "", 1).strip()
        start = cleaned.find("{")
        end = cleaned.rfind("}")
        if start == -1 or end == -1 or end <= start:
            raise ValueError("模型返回内容中未找到有效 JSON")
        return json.loads(cleaned[start : end + 1])
