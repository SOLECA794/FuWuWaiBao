import docx
import pandas as pd
import os
import asyncio
from typing import Dict, Any

class DocumentParser:
    """提供多种格式文档的解析能力"""

    @staticmethod
    def parse_docx(file_path: str) -> str:
        doc = docx.Document(file_path)
        return "\n".join([para.text for para in doc.paragraphs])

    @staticmethod
    def parse_excel(file_path: str) -> str:
        df = pd.read_excel(file_path)
        return df.to_csv(index=False)

    @staticmethod
    def parse_markdown(file_path: str) -> str:
        with open(file_path, 'r', encoding='utf-8') as f:
            return f.read()

    @staticmethod
    def parse_txt(file_path: str) -> str:
        with open(file_path, 'r', encoding='utf-8') as f:
            return f.read()

    @classmethod
    def parse_any(cls, file_path: str) -> str:
        ext = os.path.splitext(file_path)[1].lower()
        if ext == '.docx':
            return cls.parse_docx(file_path)
        elif ext in ['.xlsx', '.xls']:
            return cls.parse_excel(file_path)
        elif ext == '.md':
            return cls.parse_markdown(file_path)
        elif ext == '.txt':
            return cls.parse_txt(file_path)
        else:
            return f"Unsupported format: {ext}"

async def run_batch_parse(file_paths: list[str]) -> Dict[str, str]:
    """异步并发解析"""
    loop = asyncio.get_event_loop()
    results = {}
    for path in file_paths:
        # 在独立的线程池中解析，避免阻塞异步链路
        results[os.path.basename(path)] = await loop.run_in_executor(None, DocumentParser.parse_any, path)
    return results

if __name__ == "__main__":
    # 测试代码
    # test_path = "path_to_your_test_file"
    # print(DocumentParser.parse_any(test_path))
    pass
