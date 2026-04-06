<<<<<<< HEAD
﻿from __future__ import annotations
=======
from __future__ import annotations
>>>>>>> d17b116d297b507f8a5227ba4474640a7e13e8e0

import json
import os
from typing import Any
from urllib.parse import quote_plus

from openai import OpenAI

try:
    from .qa import resolve_llm_base_url
except ImportError:
    from qa import resolve_llm_base_url


def _extract_json_payload(text: str) -> dict[str, Any]:
    cleaned = (text or "").strip()
    if cleaned.startswith("```"):
        cleaned = cleaned.strip("`")
        cleaned = cleaned.replace("json\n", "", 1).strip()
    start = cleaned.find("{")
    end = cleaned.rfind("}")
    if start == -1 or end == -1 or end <= start:
        return {}
    try:
        return json.loads(cleaned[start : end + 1])
    except Exception:
        return {}


def _ensure_https_url(url: str, keyword: str = "") -> str:
    cleaned_url = str(url or "").strip()
    if cleaned_url.startswith("https://"):
        return cleaned_url
    if cleaned_url.startswith("http://"):
        return "https://" + cleaned_url[len("http://") :]

    safe_keyword = quote_plus((keyword or "选择排序").strip() or "选择排序")
    return f"https://search.bilibili.com/all?keyword={safe_keyword}"


def _selection_sort_fallback(keyword: str) -> list[dict[str, Any]]:
    raw_keyword = (keyword or "选择排序").strip() or "选择排序"
    encoded_keyword = quote_plus(raw_keyword)
    return [
        {
            "title": f"B站：{raw_keyword}详解（动画版）",
            "url": f"https://search.bilibili.com/all?keyword={encoded_keyword}%20%E5%8A%A8%E7%94%BB",
            "source": "bilibili",
            "fit_reason": f"围绕{raw_keyword}，动态演示清晰，适合课堂直观讲解与步骤拆解。",
            "type": "网课",
        },
        {
            "title": f"B站：{raw_keyword}手写与复杂度分析",
            "url": f"https://search.bilibili.com/all?keyword={encoded_keyword}%20%E5%A4%8D%E6%9D%82%E5%BA%A6",
            "source": "bilibili",
            "fit_reason": f"兼顾{raw_keyword}实现与复杂度分析，适合板书推导与巩固训练。",
            "type": "网课",
        },
        {
            "title": f"B站：{raw_keyword}专题课程",
            "url": f"https://search.bilibili.com/all?keyword={encoded_keyword}%20%E6%95%99%E5%AD%A6",
            "source": "bilibili",
            "fit_reason": f"内容由浅入深，便于将{raw_keyword}用于课堂导入与复盘总结。",
            "type": "网课",
        },
    ]


def _normalize_items(raw_items: Any) -> list[dict[str, Any]]:
    if not isinstance(raw_items, list):
        return []

    normalized: list[dict[str, Any]] = []
    for item in raw_items:
        if not isinstance(item, dict):
            continue
        title = str(item.get("title") or "").strip()
        url = str(item.get("url") or item.get("link") or "").strip()
        source = str(item.get("source") or "bilibili").strip() or "bilibili"
        fit_reason = str(item.get("fit_reason") or item.get("reason") or "").strip()
        item_type = str(item.get("type") or "网课").strip() or "网课"

        if not title:
            continue
        url = _ensure_https_url(url, title)
        if not fit_reason:
            fit_reason = "面向高中计算机教学，动态演示清晰，便于课堂讲解与练习。"

        normalized.append(
            {
                "title": title,
                "url": url,
                "source": source,
                "fit_reason": fit_reason,
                "type": item_type,
            }
        )

    return normalized


def get_recommendations(payload: dict[str, Any]) -> dict[str, Any]:
    """100% 大模型生成推荐，不查本地数据库；失败时强制返回 3 条兜底资源。"""
    keyword = str((payload or {}).get("keyword") or "").strip() or "选择排序"
    api_key = (os.getenv("AI_API_KEY") or "").strip()
    model = (os.getenv("AI_MODEL") or "qwen-turbo").strip()
    base_url, _ = resolve_llm_base_url(model)

    fallback_items = _selection_sort_fallback(keyword)

    if not api_key:
        return {"recommended_resources": fallback_items}

    try:
        client = OpenAI(api_key=api_key, base_url=base_url)
        system_prompt = (
            "你是智能教研资源推荐专家。"
            "必须只输出 JSON，根键名必须是 recommended_resources。"
            "recommended_resources 必须是数组，且至少 3 条。"
            "每条必须包含 title, url, source, fit_reason, type。"
            "如果没有精准资源，也要根据常识生成 B 站或知名平台可访问链接，绝不允许空数组。"
            f"本次核心搜索词是：{keyword}。请围绕该词生成推荐。"
        )
        user_prompt = (
            f"教师输入: {json.dumps(payload, ensure_ascii=False)}\n"
            f"核心关键词: {keyword}\n"
            "请直接返回 JSON。"
        )

        response = client.chat.completions.create(
            model=model,
            temperature=0.7,
            timeout=12,
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_prompt},
            ],
        )
        raw_content = (response.choices[0].message.content or "").strip()
        print(f"AI Raw Response: {raw_content}", flush=True)
        print(f"DEBUG - AI Keyword: {keyword}", flush=True)

        parsed = _extract_json_payload(raw_content)
        raw_items = parsed.get("recommended_resources")
        if not isinstance(raw_items, list):
            raw_items = parsed.get("resources")
        items = _normalize_items(raw_items)
        if len(items) < 3:
            items.extend(fallback_items[len(items) : 3])

        return {"recommended_resources": items[:3]}
    except Exception as error:
        print(f"AI Raw Response: request_failed:{error}", flush=True)
        print(f"DEBUG - AI Keyword: {keyword}", flush=True)
        return {"recommended_resources": fallback_items}
