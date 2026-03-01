
from fastapi import FastAPI, UploadFile, File, HTTPException, Request
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import os
import shutil
import asyncio
from parser import DocumentParser
from agent import ExtractionAgent

app = FastAPI(title="A23 多源数据融合系统")

# 阿里云 DashScope 配置
API_KEY = "sk-ccc7bd9a1e694352b9a70cb6d0620fdd"
BASE_URL = "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions"

UPLOAD_DIR = r"d:\Desktop\FuWuWaiBao\A23\data\uploads"
os.makedirs(UPLOAD_DIR, exist_ok=True)

# 配置跨域
# 允许来自本地前端的跨域请求（避免使用 '*' 与 allow_credentials 一起）
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:8080", "http://127.0.0.1:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
def read_root():
    return {"status": "ok", "project": "A23 - Document Intelligence and Data Fusion"}
@app.post("/process")
async def process_files(request: Request, files: list[UploadFile] = File(...)):
    if not files:
        raise HTTPException(status_code=400, detail="No files uploaded")
    
    agent = ExtractionAgent(API_KEY, BASE_URL)
    results = []
    
    # 异步并发执行函数
    async def process_single_file(file: UploadFile):
        file_path = os.path.join(UPLOAD_DIR, file.filename)
        # 保存文件
        with open(file_path, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)
        print(f"[Backend] Saved uploaded file: {file.filename}")
        
        try:
            # 1. 解析文本
            text = DocumentParser.parse_any(file_path)
            # 2. AI 提取数据
            data = await agent.extract_agricultural_data(text)
            return {"filename": file.filename, "status": "success", "data": data}
        except Exception as e:
            return {"filename": file.filename, "status": "error", "error": str(e)}

    # 记录请求来源并并发处理所有文件
    client_host = request.client.host if request.client else 'unknown'
    print(f"[Backend] Received /process request from {client_host}, files: {[f.filename for f in files]}")
    # 并发处理所有文件
    tasks = [process_single_file(f) for f in files]
    results = await asyncio.gather(*tasks)
    
    return {"results": results}

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=8001)
