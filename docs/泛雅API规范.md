【A12】基于泛雅平台的AI互动智课生成与实时问
答系统【超星集团】
开放API设计规范与示例
一、API 设计总则
1.1 设计目标
遵循命题文件中“开放式框架、标准化接口、可扩展适配 ”核心要求，构建支持“课件解析-智课生成-实时问答-进度适配 ”全流程的开放API体系，确保后续可无缝对接各类主流教育平台，同时满足师生双端用户对智能互动教学的核心需求。
1.2 设计原则
.     开放性：采用松耦合模块化设计，接口参数与返回格式标准化，支持跨平台、跨终端集成；
.     专业性：贴合教育场景特性，接口功能覆盖教学全流程，确保生成内容与交互响应的教育适配性；
.     安全性：遵循《教育数据安全指南》与《个人信息保护法》，所有接口均需签名验证，敏感数据加密传输；
.     可扩展性：预留功能扩展字段与接口版本控制机制，支持后续新增多模态交互、数字人讲授等功能；
.     易用性：接口命名简洁明了，参数设计精简必要项，返回结果结构化，降低对接开发成本。
1.3 通用规范
.     接口协议：采用 HTTP/HTTPS 协议，推荐 HTTPS 加密传输；
.     请求方法：GET（查询类）、POST（提交/生成类）、PUT（更新类）、DELETE（删除类）；
.     数据格式：请求与返回数据统一采用 JSON 格式，编码为 UTF–8；
.     签名机制：采用 MD5 签名验证，确保请求合法性与数据完整性；
.     响应状态：通过状态码统一标识请求结果，200 为成功，4xx 为客户端错误，5xx为服务端错误；
.     版本控制：接口 URL  中包含版本号（如/v1），支持多版本并行兼容。
1.4 签名验证规则
1.   签名参数：enc（必填），通过指定算法计算得出；
2.   签名算法：enc=MD5(参数有序拼接+staticKey+time)，其中：
o  参数有序拼接：按参数名 ASCII 升序排列所有非空请求参数，拼接格式为“key1value1key2value2...”；
o  staticKey：对接双方协商确定的固定密钥；

o  time：当前时间，格式为“yyyy–MM–ddHH:mm:ss”；
 
3.   验证流程：服务端接收请求后，按相同规则计算签名，与请求参数中的 enc 比
对，一致则通过验证，否则返回403 错误。
1.5 通用响应格式
json
{

"code": 200,                // 状态码：200 成功，4xx 客户端错误，

5xx 服务端错误

"msg": "操作成功",           // 状态描述

"data": {},                // 业务数据（成功时返回，格式随接口
变化）
"requestId": "req20240520001" // 请求唯一标识（用于问题排查）}
二、核心模块API 设计规范与示例
2.1 智课智能生成模块API
2.1.1 课件上传与解析接口
.     接口功能：接收 PPT/PDF 格式课件文件，解析知识点层级、公式图表、重点标注等结构化信息；
.     接口地址：/api/v1/lesson/parse
.     请求方法：POST
.     请求参数：
表格

参数名	类型	是否
必填	说明	示例

schoolId	String	是	学校 ID（对接时提供）	"sch10001"


userId	
String	
是	用户 ID（教师工号/平台用户唯一标识）	
"tea20001"

courseId	String	是	课程 ID（关联课程体系）	"cou30001"

fileType	String	是	课件文件类型	"ppt"、"pdf"


"http://xxx.com/course/ ppt/123.pdf"



 

参数名	类型	是否
必填	说明	示例
isExtractKeyPoint	Boolean	否	是否自动提取重点（默认 true）	true
enc	String	是	签名信息	"C4C859FB8E0E2035DB8C69 212E89838A"
.     返回数据示例：
json
{
"code": 200,

"msg": "课件解析成功",

"data": {

"parseId": "parse20240520001", // 解析任务 ID（用于查询生成结果）
"fileInfo": {

"fileName": "材料力学-梁弯曲理论.pptx",

"fileSize": 2048000, // 文件大小（字节）

"pageCount": 25     // 页数

},

"structurePreview": { // 知识点结构预览

"chapters": [
{
"chapter Id": "chap001",

"chapterName": "梁弯曲理论基础",

"subChapters": [
{
"subChapter Id": "sub001",

"subChapterName": "平面假设的定义",
 
"isKeyPoint": true,

"pageRange": "3-5" // 对应课件页数

}
]
}
]
},

"taskStatus": "processing" // 任务状态：processing（处理中）、 completed（完成）、failed（失败）
},
"requestId": "req20240520001"}
2.1.2 智课脚本生成接口
.     接口功能：基于课件解析结果，生成符合教学逻辑的结构化讲授脚本；
.     接口地址：/api/v1/lesson/generateScript
.     请求方法：POST
.     请求参数：
表格
参数名	类型	是否
必填	说明	示例
parseId	String	是	课件解析任务 ID	"parse20240520001"

teachingStyle	
String	
否	讲授风格（默认"standard"）	"standard"（标准）、
"detailed"（详细）、"concise" （简洁）

speechSpeed	
String	
否	语速适配（用于语音合成，默认
"normal"）	
"slow"、"normal"、"fast"
customOpening	String	否	自定义开场白	"同学们好，今天我们学习梁弯曲理论的核心知识点"
enc	String	是	签名信息	"D7E3F9A2B4C6D8E0F2A3B5C7D 9E1F3A5"
.     返回数据示例：
json
{
"code": 200,
 
"msg": "脚本生成成功",

"data": {

"scriptId": "script20240520001", // 脚本 ID（用于编辑、语音合成）
"scriptStructure": [
{
"section Id": "sec001",

"sectionName": "开场白",

"content": "同学们好，今天我们将深入学习材料力学中的梁弯曲理论，这部分内容是后续工程结构设计的重要基础，首先我们从核心假设——平面假设开始讲起。",
"duration": 15, // 预计讲授时长（秒）

"relatedChapter Id": ""
},
{
"section Id": "sec002",

"sectionName": "平面假设的定义",

"content": "平面假设是梁弯曲理论的基本假设，指梁变形前垂直于轴线的平面截面，变形后仍保持为平面且垂直于变形后的轴线 ...",
"duration": 45,
"relatedChapter Id": "sub001",
"relatedPage": "3-5",

"keyPoints": ["平面假设的核心内涵", "变形前后截面特性", "

假设的工程意义"]

}
],
 
"editUrl":
"http://xxx.com/script/edit?scriptId=script20240520001", // 脚本编辑地址
"audioGenerateUrl":
"http://xxx.com/api/v1/lesson/generateAudio" // 语音合成接口地址
},
"requestId": "req20240520002"}
2.1.3 语音合成接口
.     接口功能：将结构化脚本转换为语音音频，支持对接通用语音合成工具；
.     接口地址：/api/v1/lesson/generateAudio
.     请求方法：POST
.     请求参数：
表格
参数名	类型	是否
必填	说明	示例
scriptId	String	是	脚本 ID	"script20240520001"
voiceType	String	否	语音类型（默认
"female_standard"）	"female_standard"、 "male_professional"
audioForm at	String	否	音频格式（默认"mp3"）	"mp3"、"wav"
sectionId s	Array<Strin g>	否	指定合成的章节 ID （默认全部）	["sec001","sec002"]
enc	String	是	签名信息	"E8F9A3B5C7D9E1F3A5B7C9 D1E3F5A7B9"
.     返回数据示例：
json
{
"code": 200,

"msg": "语音合成成功",

"data": {
"audioId": "audio20240520001",
"audioUrl":
"http://xxx.com/audio/lesson/20240520001.mp3", // 音频文件
URL
 
"audioInfo": {

"totalDuration": 600, // 总时长（秒）

"fileSize": 9600000,  // 文件大小（字节）

"format": "mp3",

"bitRate": 128000   // 比特率

},

"sectionAudios": [ // 分章节音频信息

{
"section Id": "sec001",
"audioUrl":
"http://xxx.com/audio/section/sec001.mp3",
"duration": 15
}
]
},
"requestId": "req20240520003"}
2.2 多模态实时问答模块API
2.2.1 问答交互接口
.     接口功能：接收学生文字/语音提问，结合课程上下文生成精准解答，支持多轮交互；
.     接口地址：/api/v1/qa/interact
.     请求方法：POST
.     请求参数：
表格
参数名	类型	是否
必填	说明	示例
schoolId	String	是	学校 ID	"sch10001"
userId	String	是	学生学号/用户ID	"stu20001"
courseId	String	是	课程 ID	"cou30001"
lessonId	String	是	智课 ID（关联当前学习的智课）	"lesson20240520001"
sessionId	String	是	会话 ID（多轮交	"ses20240520001"

 


参数名	类型	是否
必填	说明	示例

			互标识，首次请求可生成）	

questionType	String	是	提问类型	"text"（文字）、"voice" （语音）

"平面假设为什么能简化梁弯曲问题？"、
"http://xxx.com/questio

n/voice/123.wav"


currentSection Id	String	是	当前学习章节 ID	"sec002"

[{"question":"什么是平面假设？","answer":"平面假设
是...","timestamp":"202 4-05-2010:00:00"}]

enc	String	是	签名信息	"F9A2B4C6D8E0F2A3B5C7D9 E1F3A5B7C9"

.     返回数据示例：
json
{
"code": 200,

"msg": "问答交互成功",

"data": {
"answer Id": "ans20240520001",

"answerContent": "平面假设之所以能简化梁弯曲问题，核心原因是它忽略了剪切变形对截面形状的影响，使得梁弯曲时的正应力沿截面高度呈线性分布 ...我们可以通过一个简单的例子理解：想象一根矩形截面梁，变形前截面是平面，变形后仍保持平面，这样就可以用几何关系直接推导正应力公式，无需考虑复杂的剪切变形影响。",
"answerType": "text", // 回答类型：text（文字）、mixed（图
文混合）
"relatedKnowledge": { // 关联知识点
 
"knowledgeId": "know001",

"knowledgeName": "平面假设的工程简化意义",

"relatedSection Id": "sec002"
},

"suggestions": [ // 追问建议（支持多轮交互）

"想了解平面假设的适用范围吗？",

"需要结合具体例题理解正应力分布吗？"

],

"understandingLevel": "partial" // 学生理解程度：none（未理解）、partial（部分理解）、full（完全理解）
},
"requestId": "req20240520004"}
2.2.2 语音提问识别接口
.     接口功能：将学生语音提问转换为文字，用于后续问答处理；
.     接口地址：/api/v1/qa/voiceToText
.     请求方法：POST
.     请求参数：
表格
参数名	类型	是否
必填	说明	示例
voiceUrl	String	是	语音文件 URL	"http://xxx.com/question/ voice/123.wav"

voiceDuration	Int	否	语音时长（秒）	10
language	String	否	语言类型（默认"zh-CN"）	"zh-CN"、"en-US"
enc	String	是	签名信息	"A3B5C7D9E1F3A5B7C9D1E3F5 A7B9C1D3"
.     返回数据示例：
json
{
"code": 200,
 
"msg": "语音识别成功",

"data": {

"text": "平面假设为什么能简化梁弯曲问题？",

"confidence": 0.98, // 识别置信度

"timestamp": "2024-05-20 10:05:00"
},
"requestId": "req20240520005"}
2.3 学习进度智能适配模块API
2.3.1 学习进度追踪接口
.     接口功能：记录学生学习进度、问答交互记录，为适配调整提供数据支撑；
.     接口地址：/api/v1/progress/track
.     请求方法：POST
.     请求参数：
表格
参数名	类型	是否
必填	说明	示例
schoolId	String	是	学校 ID	"sch10001"
userId	String	是	学生学号/用户ID	"stu20001"
courseId	String	是	课程 ID	"cou30001"
lessonId	String	是	智课 ID	"lesson20240520001"
currentSectionId	String	是	当前学习章节 ID	"sec002"
progressPercent	Float	是	章节学习进度（0-100）	60.5
lastOperateTime	String	是	最后操作时间	"2024-05-2010:10:00"
qaRecordId	String	否	最近问答记录 ID （如有）	"ans20240520001"
enc	String	是	签名信息	"B5C7D9E1F3A5B7C9D1E3F5A7 B9C1D3E5"
.     返回数据示例：
json
{
"code": 200,

"msg": "进度追踪成功",
 
"data": {
"trackId": "track20240520001",

"totalProgress": 45.2, // 智课总学习进度（0-100）

"nextSectionSuggest": "sec002" // 建议后续学习章节（基于问答结果适配）
},
"requestId": "req20240520006"}
2.3.2 学习节奏调整接口
.     接口功能：基于学生理解程度与学习进度，调整后续讲授节奏；
.     接口地址：/api/v1/progress/adjust
.     请求方法：POST
.     请求参数：
表格
参数名	类型	是否
必填	说明	示例
userId	String	是	学生学号/用户ID	"stu20001"
lessonId	String	是	智课 ID	"lesson20240520001"
currentSectionId	String	是	当前章节 ID	"sec002"
understandingLevel	String	是	理解程度（来自问答结果）	"partial"
qaRecordId	String	是	问答记录 ID	"ans20240520001"
enc	String	是	签名信息	"C7D9E1F3A5B7C9D1E3F5A7 B9C1D3E5F7"
.     返回数据示例：
json
{
"code": 200,

"msg": "节奏调整成功",

"data": {
"adjustPlan": {

"continueSection Id": "sec002", // 续讲章节 ID（定位原讲解节点）
 
"adjustType": "supplement", // 调整类型：supplement（补充讲解）、accelerate（加速）、normal（正常）
"supplementContent": { // 补充讲解内容（理解程度为 partial
时返回）
"content": "为了进一步理解平面假设的简化作用，我们以矩形截面梁为例 ...",
"duration": 30,

"relatedExample": "工程中常见的简支梁弯曲问题，均基于平面假设推导正应力公式"
},

"nextSections": [ // 后续章节调整建议

{
"section Id": "sec002",

"adjustedDuration": 75, // 调整后时长（秒）

"isKeyPointStrengthen": true // 是否强化重点讲解

},
{
"section Id": "sec003",
"adjustedDuration": 40,
"isKeyPointStrengthen": false
}
]
}
},
"requestId": "req20240520007"}
三、平台对接预留接口规范
3.1 课程信息同步接口
 
.     接口功能：与外部教育平台同步课程基础信息，支持智课关联课程体系；
.     接口地址：/api/v1/platform/syncCourse
.     请求方法：POST
.     请求参数（外部平台传入）：
表格
参数名	类型	是否必填	说明	示例
platformId	String	是	外部平台 ID（对接时分配）	"plat001"
courseInfo	Object	是	课程信息	详见下方示例
enc	String	是	签名信息	按双方协商规则生成
.     课程信息示例：
json
{
"courseId": "plat_cou001",

"courseName": "材料力学（上册）",

"schoolId": "sch10001",

"schoolName": "某某大学",

"teacher Info": [{"teacher Id": "plat_tea001", "teacherName": "张教授"}],
"term": "20242",
"credit": 3.0,
"period": 48,
"courseCover": "http://xxx.com/course/cover/001.jpg"} .     返回数据示例：
json
{
"code": 200,

"msg": "课程同步成功",

"data": {

"internalCourseId": "cou30001", // 系统内部课程 ID

"syncStatus": "success",
"syncTime": "2024-05-20 11:00:00"
 
},
"requestId": "req20240520008"}
3.2 用户信息同步接口
.     接口功能：同步外部平台用户信息（教师/学生），支持权限校验与身份识别；
.     接口地址：/api/v1/platform/syncUser
.     请求方法：POST
.     请求参数（外部平台传入）：
表格
参数名	类型	是否必填	说明	示例
platformId	String	是	外部平台 ID	"plat001"
userInfo	Object	是	用户信息	详见下方示例
enc	String	是	签名信息	按双方协商规则生成
.     用户信息示例：
json
{
"userId": "plat_stu001",

"userName": "李四",

"role": "student", // 角色：student（学生）、teacher（教师）

"schoolId": "sch10001",

"relatedCourseIds": ["plat_cou001"], // 关联课程 ID（外部平台）
"contactInfo": {"phone": "13800138000", "email":
"lisi@xxx.com"}}
.     返回数据示例：
json
{
"code": 200,

"msg": "用户同步成功",

"data": {

"internalUserId": "stu20001", // 系统内部用户 ID
 
"authToken": "eyJhbGciOiJIU zI1NiIsInR5cCI6IkpXVCJ9..." // 身份验证令牌
},
"requestId": "req20240520009"}
四、数据安全与隐私保护规范
4.1 数据传输安全
.     所有 API 接口推荐采用 HTTPS 加密传输，敏感字段（如用户手机号、身份令牌）需额外加密；
.     签名密钥定期更换，避免长期使用同一密钥导致安全风险；
.     接口请求设置超时机制（默认 30 秒），防止恶意请求占用资源。
4.2 数据存储安全
.     课件内容、学习数据、交互记录等敏感信息存储时采用AES–256 加密；
.     学生提问内容、个人信息等隐私数据按最小必要原则收集，不存储无关信息；
.     数据保留期限遵循教育行业规范，超出期限自动脱敏或删除。
4.3 访问权限控制
.     基于角色的权限管理（RBAC），不同角色（学生、教师、平台管理员）分配不同接口访问权限；
.     接口调用需验证用户身份令牌，非法请求直接返回401 错误；
.     记录接口调用日志（含调用方、时间、操作内容），日志保留至少 6 个月，便于安全审计。
五、接口版本管理与兼容说明
5.1 版本控制规则
.     接口版本号包含在 URL  中（如/api/v1/lesson/parse），主版本号（v1）变更表示不兼容的接口调整；
.     次版本迭代（如 v1.1）仅新增功能或参数，不修改原有参数与返回格式，确保向下兼容。
5.2 兼容策略
.     新增参数时均设为非必填项，提供默认值；
.     废弃参数时，接口仍可接收但不处理，返回数据中不再包含废弃字段，并在文档中明确标注；
.     重大版本更新时，保留旧版本接口至少 6 个月过渡期，提供版本迁移指南。
六、接口调用错误码说明
表格
错误码    说明                   处理建议
200   操作成功    正常处理返回数据
400   参数错误    检查参数格式、必填项是否完整
401   未授权      验证用户身份令牌是否有效
403   签名验证失败检查签名参数、密钥、时间戳是否正确
404   资源不存在  确认课件 ID、脚本 ID 等资源标识是否有效
408   请求超时    缩短请求数据量或稍后重试
500   服务端错误  联系技术支持，提供 requestId 排查问题
 
错误码    说明                   处理建议
503   服务暂不可用服务升级或负载过高，稍后重试
