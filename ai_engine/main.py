from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import shutil
import os
import uuid
from pathlib import Path
from typing import List, Optional

# 导入你之前写的 AI 核心逻辑
try:
    from .parser import DocumentParser
    from .generator import LessonGenerator, GenerationConfig
    from .qa import QAResponder, QAConfig
except ImportError:
    from parser import DocumentParser
    from generator import LessonGenerator, GenerationConfig
    from qa import QAResponder, QAConfig

app = FastAPI(title="泛雅 AI 智课系统后端", description="为前端提供解析、生成、问答、重讲等核心 AI 接口")

# 启用跨域支持，让前端 Vue/React 能访问
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# 临时文件存放目录
UPLOAD_DIR = Path("uploads")
UPLOAD_DIR.mkdir(exist_ok=True)

# 内存中的简单“状态机”：真实环境下这里会接数据库
# 存储 key: doc_id, value: { "parsed": ..., "generated": ... }
STORAGE = {}

# --- 数据模型定义 ---

class ChatRequest(BaseModel):
    doc_id: str
    question: str
    page: int
    mode: str = "llm"  # 默认改为 llm

# --- 核心接口实现 ---

@app.post("/upload")
async def upload_document(file: UploadFile = File(...)):
    """
    1. 上传接口：接收 PDF/PPTX，并自动触发解析
    对应：智课生成模块 - 文档解析
    """
    file_id = f"doc_{uuid.uuid4().hex[:8]}"
    file_ext = Path(file.filename).suffix
    file_path = UPLOAD_DIR / f"{file_id}{file_ext}"

    with file_path.open("wb") as buffer:
        shutil.copyfileobj(file.file, buffer)

    try:
        # 实时解析文档
        parser = DocumentParser(str(file_path))
        parsed_data = parser.parse()
        
        # 记录到内存（实际开发中这里会存入数据库）
        STORAGE[file_id] = {
            "parsed": parsed_data,
            "file_path": str(file_path)
        }
        
        return {
            "doc_id": file_id,
            "doc_name": file.filename,
            "total_pages": parsed_data["total_pages"],
            "message": "文件解析成功"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"解析失败: {str(e)}")


@app.get("/lessons/{doc_id}")
async def get_lessons(doc_id: str, mode: str = "llm"):
    """
    2. 智课生成接口：获取生成的讲稿和思维导图
    对应：创新点1 (导图) & 智课生成模块 (脚本)
    """
    if doc_id not in STORAGE:
        raise HTTPException(status_code=404, detail="文档未找到")
    
    # 如果还没生成过，则现在生成
    if "generated" not in STORAGE[doc_id]:
        generator = LessonGenerator(GenerationConfig(mode=mode))
        generated_data = generator.generate(STORAGE[doc_id]["parsed"])
        STORAGE[doc_id]["generated"] = generated_data

    return STORAGE[doc_id]["generated"]


@app.post("/chat")
async def chat_with_ai(req: ChatRequest):
    """
    3. 实时问答与案例重讲接口
    对应：创新点2 (溯源) & 创新点3 (重讲案例)
    """
    if req.doc_id not in STORAGE:
        # 如果 doc_id 不存在且 STORAGE 非空，尝试匹配最近的一个文档（方便调试）
        if STORAGE:
            last_id = list(STORAGE.keys())[-1]
            doc_data = STORAGE[last_id]
        else:
            raise HTTPException(status_code=404, detail="未找到任何已解析文档，请先上传。")
    else:
        doc_data = STORAGE[req.doc_id]
    
    parsed_doc = doc_data["parsed"]
    
    # 强制尝试获取环境变量中的模式，如果请求是 mock 也会被改为环境变量设定的模式
    # 这样用户在 Postman 里如果不改参数，也能生效
    env_mode = os.getenv("AI_GEN_MODE", "mock")
    final_mode = env_mode if env_mode == "llm" else req.mode

    # 调用你写的 QA 响应器
    responder = QAResponder(parsed_doc, config=QAConfig(mode=final_mode))
    result = responder.answer(req.question, req.page)
    
    return result


@app.get("/health")
def health_check():
    """健康检查接口，供前端联调测试"""
    return {"status": "running", "version": "v1.0.0"}

if __name__ == "__main__":
    import uvicorn
    # 为了演示，直接运行 main.py 启动
    uvicorn.run(app, host="0.0.0.0", port=8000)
