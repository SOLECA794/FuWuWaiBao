import os
import time
import uuid
from typing import Any


def build_document_schema(
    file_path: str,
    document_type: str,
    total_pages: int,
    parsed_pages: list[dict[str, Any]],
    started_at: float,
) -> dict[str, Any]:
    elapsed_ms = int((time.time() - started_at) * 1000)
    non_empty_pages = len(parsed_pages)
    empty_pages = max(total_pages - non_empty_pages, 0)

    return {
        "doc_id": f"doc_{uuid.uuid4().hex[:12]}",
        "doc_name": os.path.basename(file_path),
        "doc_path": file_path,
        "doc_type": document_type,
        "total_pages": total_pages,
        "parsed_pages": parsed_pages,
        "stats": {
            "non_empty_pages": non_empty_pages,
            "empty_pages": empty_pages,
            "elapsed_ms": elapsed_ms,
        },
    }


def build_reconstruction_schema(
    parsed_document: dict[str, Any],
    chapters: list[dict[str, Any]],
    teaching_nodes: list[dict[str, Any]],
    started_at: float,
) -> dict[str, Any]:
    elapsed_ms = int((time.time() - started_at) * 1000)

    return {
        "doc_id": parsed_document.get("doc_id"),
        "doc_name": parsed_document.get("doc_name"),
        "doc_type": parsed_document.get("doc_type"),
        "chapters": chapters,
        "teaching_nodes": teaching_nodes,
        "stats": {
            "chapter_count": len(chapters),
            "node_count": len(teaching_nodes),
            "elapsed_ms": elapsed_ms,
        },
    }


def build_node_script_schema(
    node_id: str,
    title: str,
    script: str,
    mindmap_markdown: str,
    interactive_questions: list[str],
    reteach_script: str,
    transition: str,
) -> dict[str, Any]:
    return {
        "node_id": node_id,
        "title": title,
        "script": script,
        "mindmap_markdown": mindmap_markdown,
        "interactive_questions": interactive_questions,
        "reteach_script": reteach_script,
        "transition": transition,
    }
