@echo off
REM Dify 集成验证脚本

echo.
echo ========================================
echo Dify 后端集成验证
echo ========================================
echo.

REM 检查编译文件是否存在
if not exist "api.exe" (
    echo [ERROR] api.exe 不存在，需要先编译（go build ./api）
    exit /b 1
)
echo [✓] api.exe 编译成功（%~z1 bytes）

REM 验证关键配置文件
echo.
echo 关键文件检查：
if exist "..\.env" (
    echo [✓] .env 配置文件存在
) else (
    echo [ERROR] .env 文件不存在
    exit /b 1
)

if exist "pkg\config\config.go" (
    echo [✓] config.go 存在
) else (
    echo [ERROR] config.go 不存在
    exit /b 1
)

if exist "internal\service\dify_client.go" (
    echo [✓] dify_client.go 存在
) else (
    echo [ERROR] dify_client.go 不存在
    exit /b 1
)

REM 检查关键代码存在
echo.
echo 关键代码验证：
findstr /M "NewDifyClient" "internal\service\dify_client.go" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] NewDifyClient 函数存在
) else (
    echo [ERROR] NewDifyClient 函数未找到
    exit /b 1
)

findstr /M "if cfg.AI.UseDify" "api\main.go" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] main.go 中的 Dify 选择逻辑存在
) else (
    echo [ERROR] Dify 选择逻辑未找到
    exit /b 1
)

echo.
echo ========================================
echo 验证成功！Dify 集成已完成
echo ========================================
echo.
echo 启用 Dify 的步骤：
echo 1. 在 .env 中设置：APP_AI_USE_DIFY=true
echo 2. 设置 Dify API Key：APP_AI_DIFY_API_KEY=&lt;key&gt;
echo 3. 运行后端：api.exe
echo 4. POST http://localhost:18080/api/v1/ai/coursewares/test/ask
echo.
