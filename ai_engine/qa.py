import json
import os
import re
from dataclasses import dataclass
from typing import Any
from dotenv import load_dotenv

load_dotenv()


RETEACH_KEYWORDS = [
    "听不懂",
    "不懂",
    "太难",
    "没明白",
    "再讲",
    "换个说法",
    "举例",
    "不太会",
]


@dataclass
class QAConfig:
    mode: str = "llm"  # 默认使用真实模型
    model: str = "qwen-plus"
    temperature: float = 0.2


class QAResponder:
    """问答溯源 + 案例重讲响应器。"""

    def __init__(self, parsed_document: dict[str, Any], config: QAConfig | None = None):
        self.parsed_document = parsed_document
        self.config = config or QAConfig(
            mode=os.getenv("AI_GEN_MODE", "llm"),
            model=os.getenv("AI_MODEL", "qwen-plus"),
            temperature=float(os.getenv("AI_TEMPERATURE", "0.2")),
        )
        self.page_map = {
            item["page"]: item.get("content", "")
            for item in self.parsed_document.get("parsed_pages", [])
            if "page" in item
        }

    def answer(self, question: str, current_page: int) -> dict[str, Any]:
        source_page = self._resolve_source_page(current_page)
        source_content = self.page_map.get(source_page, "")
        need_reteach = self._need_reteach(question)

        # 始终使用 LLM 进行回答
        answer_text = self._llm_answer(question, source_page, source_content, need_reteach)

        return {
            "question": question,
            "source_page": source_page,
            "source_excerpt": source_content[:160].strip(),
            "intent": {
                "need_reteach": need_reteach,
                "reason": "keyword_trigger" if need_reteach else "normal_qa",
            },
            "answer": answer_text,
            "resume_page": source_page,
            "follow_up_suggestion": self._make_followup_suggestion(question, need_reteach, source_page),
        }

    def _resolve_source_page(self, current_page: int) -> int:
        if current_page in self.page_map:
            return current_page
        if not self.page_map:
            return 1
        nearest = min(self.page_map.keys(), key=lambda page: abs(page - current_page))
        return nearest

    @staticmethod
    def _need_reteach(question: str) -> bool:
        q = question.strip().lower()
        return any(keyword in q for keyword in RETEACH_KEYWORDS)

    def _llm_answer(self, question: str, source_page: int, source_content: str, need_reteach: bool) -> str:
        try:
            from openai import OpenAI
        except ImportError as error:
            return self._fallback_answer(question, source_content, need_reteach, f"missing_openai:{error}")

        api_key = os.getenv("AI_API_KEY")
        if not api_key:
            return self._fallback_answer(question, source_content, need_reteach, "missing_api_key")

        base_url = os.getenv("AI_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1")
        client = OpenAI(api_key=api_key, base_url=base_url)

        role_desc = "高校专业课老师" if not need_reteach else "善于举例的耐心导师"
        system_prompt = (
            f"你是{role_desc}。当前正在讲解一份文档。"
            "请基于参考内容回答学生问题。如果内容不足以回答，请根据专业知识点拨。"
            "如果学生表示听不懂（need_reteach=True），请务必用一个更生活化、更通俗的案例来类比解析。"
        )
        
        user_prompt = (
            f"学生问题: {question}\n"
            f"当前文档第 {source_page} 页参考内容: {source_content}\n"
            f"是否需要重讲(Reteach): {need_reteach}\n"
            "请给出专业且易懂的回答。"
        )

        try:
            response = client.chat.completions.create(
                model=self.config.model,
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_prompt}
                ],
                temperature=self.config.temperature
            )
            return response.choices[0].message.content
        except Exception as e:
            return self._fallback_answer(question, source_content, need_reteach, str(e))

    def _fallback_answer(self, question: str, source_content: str, need_reteach: bool, reason: str) -> str:
        lines = [line.strip() for line in re.split(r"[\n。；;]", source_content or "") if line.strip()]
        preview = "；".join(lines[:3]) if lines else "当前页缺少足够上下文，建议教师补充讲稿或原文。"
        if need_reteach:
            return (
                f"我先换一种更容易理解的方式说明。当前这部分主要在讲：{preview}。"
                "你可以把它理解为先明确核心概念，再看它在例子里是怎么工作的。"
                f"如果还不清楚，我们可以继续围绕“{question}”拆成更细的步骤。"
            )
        return (
            f"结合当前课件内容，这个问题可以先抓住这些信息：{preview}。"
            f"如果你想继续追问“{question}”，我可以再展开原理、例子或推导过程。"
            f"（当前回答使用本地兜底模式，原因：{reason}）"
        )

    def _make_followup_suggestion(self, question: str, need_reteach: bool, source_page: int) -> str:
        """生成更贴近实际的后续操作建议（用于前端展示给学生选择）。"""
        if need_reteach:
            return (
                "我可以：1) 用更生活化的例子重新解释；"
                "2) 按步骤分解并配上简短练习；"
                "3) 给出一个可操作的小实验/图示帮助理解。你想先听哪个？"
            )

        # 对于正常提问，给出更具体的后续选项
        return (
            "我可以：1) 展开详细原理与推导；"
            "2) 给出一个相关的真实例子或应用场景；"
            "3) 把相关页的要点串成一份简短笔记。请选择你想要的方式。"
        )
        

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
