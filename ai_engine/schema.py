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
