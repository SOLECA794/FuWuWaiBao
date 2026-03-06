from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import shutil
import os
import uuid
import json
from pathlib import Path
from typing import List, Optional, Any

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


class GenerateScriptRequest(BaseModel):
    page: int
    content: str
    course_name: Optional[str] = None
    mode: str = "llm"


class AskWithContextRequest(BaseModel):
    question: str
    current_page: int
    context: str = ""
    mode: str = "llm"


class ParseKnowledgeRequest(BaseModel):
    text: str
    mode: str = "llm"

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
    env_mode = os.getenv("AI_GEN_MODE", "llm")
    final_mode = env_mode if env_mode == "llm" else req.mode

    # 调用你写的 QA 响应器
    responder = QAResponder(parsed_doc, config=QAConfig(mode=final_mode))
    result = responder.answer(req.question, req.page)
    
    return result


@app.post("/generate-script")
async def generate_script(req: GenerateScriptRequest):
    """基于页内容直接生成讲稿与导图（供 Go 后端调用）。"""
    try:
        parsed_document = {
            "doc_id": f"inline_{uuid.uuid4().hex[:8]}",
            "doc_name": req.course_name or "inline_course",
            "parsed_pages": [
                {
                    "page": req.page,
                    "content": req.content.strip() or f"第{req.page}页内容",
                    "content_length": len(req.content.strip()),
                }
            ],
        }
        generator = LessonGenerator(GenerationConfig(mode=req.mode))
        result = generator.generate(parsed_document)
        lesson = (result.get("lessons") or [{}])[0]
        return {
            "page": req.page,
            "script": lesson.get("script", ""),
            "mindmap_markdown": lesson.get("mindmap_markdown", ""),
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"讲稿生成失败: {str(e)}")


@app.post("/ask-with-context")
async def ask_with_context(req: AskWithContextRequest):
    """基于上下文直接问答（供 Go 后端调用）。"""
    try:
        parsed_document = {
            "doc_id": f"inline_{uuid.uuid4().hex[:8]}",
            "doc_name": "inline_context",
            "parsed_pages": [
                {
                    "page": req.current_page,
                    "content": req.context.strip() or "当前页面暂无文本上下文",
                    "content_length": len(req.context.strip()),
                }
            ],
        }
        responder = QAResponder(parsed_document, config=QAConfig(mode=req.mode))
        return responder.answer(req.question, req.current_page)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"问答失败: {str(e)}")


def _llm_parse_knowledge(text: str, mode: str) -> dict[str, Any]:
    if mode != "llm":
        raise RuntimeError("parse-knowledge 仅支持 llm 模式")

    try:
        from openai import OpenAI
    except ImportError:
        raise RuntimeError("缺少 openai 依赖，请安装 requirements.txt")

    api_key = os.getenv("AI_API_KEY")
    if not api_key:
        raise RuntimeError("缺少环境变量 AI_API_KEY，无法调用大模型")

    base_url = os.getenv("AI_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1")
    model = os.getenv("AI_MODEL", "qwen-plus")
    client = OpenAI(api_key=api_key, base_url=base_url)

    prompt = (
        "请把输入内容拆解为知识点树，输出严格 JSON，格式为: "
        "{\"structure\":[{\"name\":\"章节\",\"children\":[{\"name\":\"知识点\",\"children\":[]}]}]}。"
        "不要输出额外解释。"
    )

    response = client.chat.completions.create(
        model=model,
        temperature=0.2,
        messages=[
            {"role": "system", "content": "你是课程知识点结构化助手。"},
            {"role": "user", "content": prompt + "\n\n输入内容:\n" + text[:5000]},
        ],
        response_format={"type": "json_object"},
    )
    content = response.choices[0].message.content
    data = json.loads(content)
    if "structure" not in data or not isinstance(data.get("structure"), list):
        raise RuntimeError("模型未返回合法 structure 字段")
    return {"structure": data["structure"]}


@app.post("/parse-knowledge")
async def parse_knowledge(req: ParseKnowledgeRequest):
    """将文本拆解为知识点树（供 Go 后端调用）。"""
    text = (req.text or "").strip()
    if not text:
        raise HTTPException(status_code=400, detail="text 不能为空")
    try:
        return _llm_parse_knowledge(text, req.mode)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"知识点解析失败: {e}")


@app.get("/health")
def health_check():
    """健康检查接口，供前端联调测试"""
    return {"status": "running", "version": "v1.0.0"}

if __name__ == "__main__":
    import uvicorn
    # 为了演示，直接运行 main.py 启动
    uvicorn.run(app, host="0.0.0.0", port=8000)
