import json
import httpx
from typing import Dict, Any, List

class ExtractionAgent:
    """提取大模型 Agent，用于结构化抽取"""

    def __init__(self, api_key: str, base_url: str):
        self.api_key = api_key
        self.base_url = base_url

    async def extract_agricultural_data(self, text: str) -> Dict[str, Any]:
        """
        利用大模型从文本中提取指定的农业指标
        """
        prompt = f"""
        你是一位专业的农业数据处理专家。请阅读以下文本，并提取出包含的统计数据。

        【目标指标】：
        1. 区域名称
        2. 村民委员会数量
        3. 总户数
        4. 总人口
        5. 粮食面积（亩）
        6. 粮食全产（公斤/亩）或总产量（吨）

        【要求】：
        - 必须返回如下严格的 JSON 格式，如果某项找不到值，请填入 null。
        - 确保数据的高准确性，数值部分统一转换为纯数字。
        - 返回的消息内容必须且只能是 JSON 字符串，不要带任何 Markdown 代码块标签。

        JSON 骨架：
        {{
            "region": "区域名称",
            "vil_comm_count": 0,
            "total_households": 0,
            "total_population": 0,
            "grain_area_mu": 0,
            "grain_yield_per_mu": 0,
            "total_grain_yield_ton": 0
        }}

        【待提取文本】：
        {text}
        """

        headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        }

        payload = {
            "model": "qwen-turbo",  # 使用通义千问高效型号
            "messages": [{"role": "user", "content": prompt}],
            "temperature": 0.1
        }

        async with httpx.AsyncClient() as client:
            response = await client.post(self.base_url, headers=headers, data=json.dumps(payload))
            if response.status_code == 200:
                # 假设 API 返回的是标准的 OpenAI 格式
                result = response.json()
                content = result['choices'][0]['message']['content']
                return json.loads(content)
            else:
                raise Exception(f"API Error: {response.text}")

# 使用异步并发处理多个文档提取
# extraction_tasks = [agent.extract_agricultural_data(t) for t in parsed_texts]
# all_results = await asyncio.gather(*extraction_tasks)
