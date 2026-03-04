🚀 技术栈
组件	技术	说明
框架	Gin	Web框架
数据库	PostgreSQL	关系型数据库
缓存	Redis	会话/进度缓存
对象存储	MinIO	课件文件存储
容器化	Docker	容器部署
语言	Go 1.21+	开发语言
📁 项目结构
text
smart-teaching-backend/
├── cmd/
│   └── api/              # 主程序入口
├── internal/
│   ├── handler/          # HTTP处理器
│   ├── service/          # 业务逻辑
│   ├── repository/       # 数据访问
│   ├── model/            # 数据模型
│   └── middleware/       # 中间件
├── pkg/
│   ├── config/           # 配置管理
│   ├── logger/           # 日志工具
│   ├── oss/              # 对象存储
│   └── ai/               # AI服务（预留）
├── config/               # 配置文件
├── API.md                # 接口文档
├── Dockerfile            # 容器构建文件
├── docker-compose.yml    # 开发环境配置
└── docker-compose.prod.yml # 生产环境配置
✨ 已完成功能
教师端接口
模块	功能	状态
📁 课件管理	上传、列表、删除、发布	✅ 完成
📝 讲稿编辑	获取、保存、AI生成	✅ 完成
📊 学情分析	页面统计、关键词分析	✅ 完成
❓ 提问记录	分页查询学生提问	✅ 完成
🖼️ 课件预览	获取预览图片	✅ 完成
基础服务
✅ PostgreSQL 数据持久化

✅ Redis 缓存

✅ MinIO 对象存储

✅ Docker 容器化部署

📚 接口文档
详细接口文档请查看 API.md

主要接口列表：

方法	路径	功能
GET	/api/teacher/courseware-list	获取课件列表
POST	/api/teacher/upload-courseware	上传课件
DELETE	/api/teacher/courseware/{courseId}	删除课件
POST	/api/teacher/publish-courseware	发布课件
GET	/api/teacher/script/{courseId}/{page}	获取讲稿
POST	/api/teacher/script/save	保存讲稿
POST	/api/teacher/ai-generate-script	AI生成讲稿
GET	/api/teacher/student-stats/{courseId}	学情分析
GET	/api/teacher/question-records/{courseId}	提问记录
GET	/api/courseware/{courseId}/page/{pageNum}	课件预览
🛠️ 环境变量配置
在 config/config.yaml 中配置：

yaml
server:
port: 8080
mode: debug

database:
host: localhost
port: 5432
user: postgres
password: 123456
dbname: teaching

redis:
host: localhost
port: 6379
password: ""

minio:
endpoint: localhost:9000
access_key: minioadmin
secret_key: minioadmin
bucket: courses

📝 更新记录
v1.0.0 (2026-03-03)
✨ 完成所有教师端接口

📚 添加完整接口文档

🐳 支持 Docker 部署

🔧 集成 PostgreSQL/Redis/MinIO