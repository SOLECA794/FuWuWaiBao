import asyncio
import os
import json
from parser import DocumentParser
from agent import ExtractionAgent

# 阿里云 DashScope 配置
API_KEY = "sk-ccc7bd9a1e694352b9a70cb6d0620fdd"
BASE_URL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions" # 华北2 兼容模式地址

async def mini_test():
    print("🚀 开始 A23 极简关键过程测试 (使用阿里云通义千问 API)...")
    
    # --- 模拟数据模式 ---
    # 为了确跑通流程，我们直接构造一段模拟的文本，不再强制依赖读取 docx
    mock_text = """
    【南通市兴东镇 2023 年统计概况】
    该镇下辖村民委员会数量共 15 个。全镇总户数达到 8500 户，总常住人口为 26800 人。
    在粮食生产方面，全镇粮食种植面积达 4200 亩，粮食亩产约为 450 公斤，
    初步估算粮食全产累计 1890 吨。
    """
    print("📂 步骤 1: 正在生成模拟农业数据文本...")
    text = mock_text
    print(f"✅ 模拟文本准备完成")
    print("-" * 30)
    print(text.strip())
    print("-" * 30)


    # 项目 2：AI 结构化提取 (Agent)
    print("🤖 步骤 2: 正在调用大模型提取结构化数据 (JSON)...")
    agent = ExtractionAgent(API_KEY, BASE_URL)
    
    try:
        # 注意：这里会产生真实的网络请求，取决于你的 API 配置
        result = await agent.extract_agricultural_data(text)
        print("✅ 提取完成! 得到 JSON 数据:")
        print(json.dumps(result, indent=4, ensure_ascii=False))
        
        # 项目 3：模拟结果保存 (未来会直接写入 Excel)
        output_json = r"d:\Desktop\FuWuWaiBao\A23\data\uploads\result_test.json"
        with open(output_json, "w", encoding="utf-8") as f:
            json.dump(result, f, indent=4, ensure_ascii=False)
        print(f"💾 结果已临时保存至: {output_json}")

    except Exception as e:
        print(f"❌ 提取失败: {e}")
        print("💡 提示：请检查 parser.py 和 agent.py 中的 API 配置是否正确，以及网络是否通畅。")

if __name__ == "__main__":
    asyncio.run(mini_test())
