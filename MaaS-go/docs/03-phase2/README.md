# 阶段2：核心功能开发

> 🚀 **阶段目标**：实现MaaS平台的核心业务流程，包括服务通信、缓存、认证、模型管理和推理服务

## 阶段总览

### 学习时间
**5周**（约25个工作日）

### 核心目标
1. 掌握gRPC服务间通信
2. 学会Redis缓存策略和设计
3. 实现JWT认证和RBAC权限控制
4. 处理大文件分片上传
5. 构建模型推理服务
6. 理解版本控制和依赖管理

### 最终产出
- ✅ API Gateway通过gRPC调用后端服务
- ✅ 完整的用户认证体系（JWT + RBAC）
- ✅ Redis多级缓存支持
- ✅ 模型文件上传和版本管理
- ✅ 可调用模型的推理接口

## 技术栈详解

### 1. gRPC与Protocol Buffers

**什么是gRPC？**
gRPC是Google开发的高性能RPC框架，基于HTTP/2和Protocol Buffers。

**为什么使用gRPC？**
- **性能高**：基于HTTP/2，支持多路复用、头部压缩
- **跨语言**：支持Go/Java/Python/C++等10+语言
- **强类型**：使用Protocol Buffers定义接口，自动生成代码
- **流式通信**：支持客户端流、服务端流、双向流

**Protocol Buffers（protobuf）**
是一种语言中立、平台中立的数据序列化格式，比JSON更小更快。

**示例对比**：
```protobuf
// protobuf定义
message User {
  string id = 1;
  string name = 2;
  string email = 3;
}
```

```json
// 相同数据的JSON表示（约80字节）
{
  "id": "123",
  "name": "John",
  "email": "john@example.com"
}

// protobuf二进制表示（约20字节）
// 更小的体积 = 更快的传输 = 更低的带宽
```

### 2. Redis缓存

**什么是Redis？**
Redis是内存中的数据结构存储系统，支持String、Hash、List、Set、Sorted Set等数据结构。

**为什么要用缓存？**
- **减少数据库压力**：热点数据从内存读取，减少DB查询
- **提高响应速度**：内存访问比磁盘快1000倍
- **降低系统成本**：减少数据库服务器配置

**缓存模式**：
1. **Cache-Aside（旁路缓存）**：应用先查缓存，没有则查DB并写入缓存
2. **Read-Through**：缓存服务自己处理DB查询
3. **Write-Through**：写数据时同时更新缓存和DB
4. **Write-Behind**：异步写DB，先写缓存

**常见问题**：
- **缓存穿透**：查询不存在的数据，每次都要查DB
  - 解决：缓存空值、布隆过滤器
- **缓存击穿**：热点key过期瞬间，大量请求打到DB
  - 解决：互斥锁、逻辑过期
- **缓存雪崩**：大量key同时过期
  - 解决：随机过期时间、多级缓存

### 3. JWT认证

**什么是JWT？**
JSON Web Token是一种开放标准（RFC 7519），用于在网络应用间安全传输信息。

**JWT结构**：
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
│                        │                             │
│                        │                             │
└── Header（头部）        └── Payload（负载）           └── Signature（签名）
    {                        {
      "alg": "HS256",          "sub": "1234567890",
      "typ": "JWT"             "name": "John Doe",
    }                          "iat": 1516239022
                             }
```

**JWT vs Session**：
| 特性 | JWT | Session |
|------|-----|---------|
| 存储位置 | 客户端 | 服务端 |
| 可扩展性 | 好（无状态） | 差（需要共享session） |
| 安全性 | 中（不能撤销） | 高（可随时删除） |
| 适用场景 | 微服务、移动端 | 传统Web应用 |

### 4. MinIO对象存储

**什么是对象存储？**
对象存储是一种数据存储架构，将数据作为对象管理，而不是文件系统或块存储。

**为什么用MinIO？**
- 开源的S3兼容对象存储
- 高性能（Go语言编写）
- 易于部署（单二进制文件）
- 支持分布式部署

**适用场景**：
- 存储大文件（模型文件、图片、视频）
- 支持断点续传、分片上传
- 需要高可用、可扩展的存储

### 5. 模型版本管理

**Semantic Versioning（语义化版本）**：
```
版本格式：主版本号.次版本号.修订号（MAJOR.MINOR.PATCH）

1.0.0
│ │ │
│ │ └── Patch：Bug修复，向后兼容
│ └──── Minor：新功能，向后兼容
└────── Major：重大变更，可能不兼容
```

**版本管理策略**：
- **latest标签**：始终指向最新版本
- **版本别名**：v1指向最新的v1.x.x
- **预发布版本**：v1.0.0-alpha, v1.0.0-beta

## 节点详解

### 节点2.1：服务间通信（5天）

**学习目标**：
- 理解RPC vs HTTP REST
- 掌握Protocol Buffers语法
- 实现gRPC服务端和客户端
- 配置服务发现和负载均衡

**技术详解**：

**RPC vs REST对比**：

| 特性 | REST | RPC（gRPC） |
|------|------|-------------|
| 协议 | HTTP/1.1 | HTTP/2 |
| 格式 | JSON/XML | Protocol Buffers |
| 可读性 | 好（文本） | 差（二进制） |
| 性能 | 一般 | 高（5-10倍） |
| 流式 | 不支持 | 支持 |
| 浏览器 | 原生支持 | 需要gRPC-Web |

**Protobuf定义示例**：
```protobuf
syntax = "proto3";

package model;

option go_package = "github.com/17882237881/MaaS/shared/proto/model";

// 模型服务
service ModelService {
  rpc CreateModel(CreateModelRequest) returns (CreateModelResponse);
  rpc GetModel(GetModelRequest) returns (Model);
  rpc ListModels(ListModelsRequest) returns (ListModelsResponse);
  rpc DeleteModel(DeleteModelRequest) returns (DeleteModelResponse);
}

message CreateModelRequest {
  string name = 1;
  string description = 2;
  string version = 3;
  string framework = 4;
  repeated string tags = 5;
  map<string, string> metadata = 6;
}

message CreateModelResponse {
  string id = 1;
  string message = 2;
}

message Model {
  string id = 1;
  string name = 2;
  string description = 3;
  string version = 4;
  string framework = 5;
  string status = 6;
  int64 size = 7;
  repeated string tags = 8;
  string created_at = 9;
}
```

**实操任务**：
1. 安装Protocol Buffers编译器
2. 创建.proto文件定义服务接口
3. 生成Go代码：`protoc --go_out=. --go-grpc_out=. *.proto`
4. 实现gRPC服务端（Model Registry）
5. 实现gRPC客户端（API Gateway调用）
6. 测试服务间调用

**检查点**：
- [ ] .proto文件定义完整
- [ ] 生成的Go代码可用
- [ ] gRPC服务端能启动
- [ ] API Gateway能成功调用Model Registry
- [ ] 支持错误处理和超时控制

---

### 节点2.2：缓存层设计（4天）

**学习目标**：
- 掌握Redis基础数据类型
- 实现多级缓存架构
- 处理缓存常见问题
- 实现缓存预热和刷新

**技术详解**：

**Redis数据类型**：

**1. String（字符串）**
```go
// 基础操作
SET key value          // 设置
GET key                // 获取
SETEX key 60 value     // 设置并指定过期时间（秒）
INCR counter           // 原子递增

// 使用场景：缓存对象、计数器、分布式锁
```

**2. Hash（哈希）**
```go
// 存储对象
HSET user:1 name "John" age 25
HGET user:1 name       // 获取单个字段
HGETALL user:1         // 获取所有字段

// 使用场景：存储对象（比JSON更节省空间）
```

**3. List（列表）**
```go
LPUSH queue job1       // 从左侧插入
RPUSH queue job2       // 从右侧插入
LPOP queue             // 从左侧弹出
LRANGE queue 0 -1      // 获取所有元素

// 使用场景：队列、最新消息列表
```

**4. Set（集合）**
```go
SADD tags:model1 tag1 tag2
SMEMBERS tags:model1   // 获取所有成员
SISMEMBER tags:model1 tag1  // 判断是否包含

// 使用场景：标签、去重
```

**多级缓存架构**：
```
请求 → L1本地缓存（BigCache）→ 未命中 → 
      L2分布式缓存（Redis）→ 未命中 → 
      L3数据库（PostgreSQL）→ 回填缓存
```

**Cache-Aside模式代码示例**：
```go
func (s *Service) GetModel(ctx context.Context, id string) (*Model, error) {
    // 1. 查本地缓存
    if model, ok := s.localCache.Get(id); ok {
        return model, nil
    }
    
    // 2. 查Redis
    key := fmt.Sprintf("model:%s", id)
    if data, err := s.redis.Get(ctx, key).Bytes(); err == nil {
        var model Model
        if err := json.Unmarshal(data, &model); err == nil {
            // 回填本地缓存
            s.localCache.Set(id, &model, time.Minute)
            return &model, nil
        }
    }
    
    // 3. 查数据库
    model, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 4. 回填缓存
    if data, err := json.Marshal(model); err == nil {
        s.redis.Set(ctx, key, data, time.Hour)
    }
    s.localCache.Set(id, model, time.Minute)
    
    return model, nil
}
```

**实操任务**：
1. 安装并配置Redis
2. 实现基础Redis操作封装
3. 实现模型查询缓存
4. 添加本地缓存层（可选）
5. 实现缓存穿透保护（布隆过滤器）
6. 编写缓存性能测试

**检查点**：
- [ ] Redis连接正常
- [ ] 缓存命中和未命中逻辑正确
- [ ] 缓存过期策略合理
- [ ] 缓存穿透/击穿问题已处理
- [ ] 缓存性能提升明显（QPS提升2倍以上）

---

### 节点2.3：认证与授权（5天）

**学习目标**：
- 实现JWT认证流程
- 设计RBAC权限模型
- 实现API密钥管理
- 添加请求鉴权中间件

**技术详解**：

**认证（Authentication）vs 授权（Authorization）**：
- **认证**：验证"你是谁"（检查凭证）
- **授权**：验证"你能做什么"（检查权限）

**JWT认证流程**：
```
1. 用户登录
   Client → POST /auth/login (username, password)
          → Server验证 → 生成JWT → 返回Token

2. 访问受保护资源
   Client → GET /api/users/me
          → Header: Authorization: Bearer <token>
          → Server验证JWT签名和过期时间
          → 返回资源

3. Token刷新
   Client → POST /auth/refresh
          → Header: Authorization: Bearer <refresh_token>
          → Server验证 → 生成新的Access Token
```

**RBAC模型（Role-Based Access Control）**：
```
User（用户） → UserRole（用户角色关联） → Role（角色） → RolePermission（角色权限关联） → Permission（权限）

示例：
- 用户：user1
- 角色：developer（开发者）
- 权限：model:create, model:read, inference:execute

管理员角色可能有：model:*, user:*, billing:*（所有权限）
```

**Casbin权限库使用**：
```go
// 定义权限模型（model.conf）
[request_definition]
r = sub, obj, act  // 请求：主体（用户），对象（资源），动作（操作）

[policy_definition]
p = sub, obj, act  // 策略：角色，资源，操作

[policy_effect]
e = some(where (p.eft == allow))  // 任一策略允许则通过

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act

// 策略文件（policy.csv）
p, admin, model, *
p, admin, user, *
p, developer, model, read
p, developer, model, create
p, developer, inference, execute
```

**实操任务**：
1. 实现用户注册和登录接口
2. 集成JWT生成和验证
3. 设计RBAC权限模型
4. 集成Casbin权限控制
5. 实现Token刷新机制
6. 添加API密钥支持

**检查点**：
- [ ] 用户能注册和登录
- [ ] JWT Token正确生成和验证
- [ ] 未认证请求返回401
- [ ] 无权限请求返回403
- [ ] Token过期能正确处理
- [ ] RBAC权限控制生效

---

### 节点2.4：模型上传与存储（5天）

**学习目标**：
- 实现大文件分片上传
- 集成MinIO对象存储
- 实现断点续传
- 添加文件校验和病毒扫描（可选）

**技术详解**：

**大文件上传方案**：

**1. 分片上传流程**：
```
1. 初始化上传
   Client → POST /upload/init
          → Server返回upload_id

2. 上传分片（并行）
   Client → PUT /upload/chunk
          → Header: X-Upload-ID, X-Chunk-Index
          → Body: chunk_data
          
3. 合并分片
   Client → POST /upload/complete
          → Body: {upload_id, chunks: [hash1, hash2, ...]}
          → Server验证并合并 → 返回文件URL
```

**2. 断点续传**：
```
Client → GET /upload/status?upload_id=xxx
       → Server返回已上传的分片索引
       → Client只上传缺失的分片
```

**MinIO集成示例**：
```go
// 初始化MinIO客户端
minioClient, err := minio.New("localhost:9000", &minio.Options{
    Creds:  credentials.NewStaticV4("access-key", "secret-key", ""),
    Secure: false,
})

// 上传文件
uploadInfo, err := minioClient.PutObject(
    ctx,
    "models",                    // bucket名称
    "user1/model-v1.pkl",       // 对象名称
    fileReader,                 // 文件内容
    fileSize,                   // 文件大小
    minio.PutObjectOptions{},
)

// 生成下载URL（1小时有效）
presignedURL, err := minioClient.PresignedGetObject(
    ctx,
    "models",
    "user1/model-v1.pkl",
    time.Hour,
    nil,
)
```

**实操任务**：
1. 部署MinIO服务
2. 实现分片上传接口（init/chunk/complete）
3. 实现断点续传功能
4. 添加文件MD5校验
5. 实现上传进度通知（WebSocket或SSE）
6. 限制文件类型和大小

**检查点**：
- [ ] 小文件上传正常
- [ ] 大文件（>100MB）分片上传正常
- [ ] 断点续传功能工作
- [ ] 文件校验正确
- [ ] 上传进度可查看
- [ ] 文件存储在MinIO中

---

### 节点2.5：模型版本管理（4天）

**学习目标**：
- 实现语义化版本控制
- 设计版本元数据结构
- 实现版本回滚
- 添加版本间差异比较

**技术详解**：

**版本管理策略**：

**1. 版本号生成**：
```go
// 手动指定
version := "1.0.0"

// 自动递增（基于现有版本）
func NextVersion(current string, changeType string) string {
    parts := strings.Split(current, ".")
    major, _ := strconv.Atoi(parts[0])
    minor, _ := strconv.Atoi(parts[1])
    patch, _ := strconv.Atoi(parts[2])
    
    switch changeType {
    case "major":
        major++
        minor = 0
        patch = 0
    case "minor":
        minor++
        patch = 0
    case "patch":
        patch++
    }
    
    return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
```

**2. 版本元数据**：
```go
type ModelVersion struct {
    ID          string
    ModelID     string      // 关联的模型ID
    Version     string      // 版本号（如 1.0.0）
    Status      string      // 状态：draft/published/deprecated
    ChangeLog   string      // 变更说明
    Size        int64       // 文件大小
    Checksum    string      // 文件校验和
    StoragePath string      // 存储路径
    CreatedBy   string      // 创建者
    CreatedAt   time.Time
    PublishedAt *time.Time  // 发布时间
    IsLatest    bool        // 是否最新版本
}
```

**3. 版本别名**：
```sql
-- latest始终指向最新发布的版本
UPDATE model_versions SET is_latest = false WHERE model_id = ?;
UPDATE model_versions SET is_latest = true WHERE id = ?;

-- 查询时可以用版本号或latest
SELECT * FROM model_versions WHERE model_id = ? AND (version = ? OR (version = 'latest' AND is_latest = true));
```

**实操任务**：
1. 设计版本数据库表结构
2. 实现版本创建和查询接口
3. 实现latest标签自动更新
4. 实现版本回滚功能
5. 添加版本变更日志
6. 实现版本比较（展示差异）

**检查点**：
- [ ] 版本号符合SemVer规范
- [ ] latest标签指向正确
- [ ] 版本历史可查询
- [ ] 版本回滚功能正常
- [ ] 版本元数据完整

---

### 节点2.6：推理服务基础（5天）

**学习目标**：
- 实现同步推理接口
- 设计请求队列和并发控制
- 实现超时和取消机制
- 添加推理结果缓存

**技术详解**：

**推理服务架构**：
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  API请求    │────▶│  推理网关    │────▶│  模型容器   │
│  (Gin)      │     │  (负载均衡)   │     │  (GPU)      │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  连接池      │
                    │  (gRPC/HTTP) │
                    └──────────────┘
```

**并发控制模式**：

**1. 信号量（Semaphore）**：
```go
type InferenceLimiter struct {
    sem chan struct{}
}

func NewInferenceLimiter(maxConcurrent int) *InferenceLimiter {
    return &InferenceLimiter{
        sem: make(chan struct{}, maxConcurrent),
    }
}

func (l *InferenceLimiter) Acquire() {
    l.sem <- struct{}{}  // 获取信号量，满了就阻塞
}

func (l *InferenceLimiter) Release() {
    <-l.sem  // 释放信号量
}

// 使用
limiter := NewInferenceLimiter(100)  // 最多100个并发

func HandleInference(w http.ResponseWriter, r *http.Request) {
    limiter.Acquire()
    defer limiter.Release()
    
    // 执行推理...
}
```

**2. 工作池（Worker Pool）**：
```go
type InferenceWorker struct {
    jobs    chan InferenceJob
    workers int
}

type InferenceJob struct {
    ModelID string
    Input   []float32
    Result  chan InferenceResult
}

func (w *InferenceWorker) Start() {
    for i := 0; i < w.workers; i++ {
        go w.runWorker()
    }
}

func (w *InferenceWorker) runWorker() {
    for job := range w.jobs {
        result := w.process(job)
        job.Result <- result
    }
}
```

**超时控制**：
```go
// 创建带超时的context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// 调用推理服务
result, err := inferenceClient.Predict(ctx, request)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        return nil, errors.New("inference timeout")
    }
    return nil, err
}
```

**实操任务**：
1. 创建推理服务框架（可以是模拟实现）
2. 实现同步推理接口
3. 添加并发控制（限制最大并发数）
4. 实现超时和取消机制
5. 添加推理结果缓存
6. 记录推理指标（延迟、成功率）

**检查点**：
- [ ] 推理接口可调用
- [ ] 并发控制生效（超出限制时排队或拒绝）
- [ ] 超时机制正常工作
- [ ] 推理结果正确返回
- [ ] 性能指标可查看

---

## 阶段2里程碑

### 完成检查清单

**服务间通信**：
- [ ] API Gateway通过gRPC调用Model Registry
- [ ] Protocol Buffers定义完整
- [ ] 支持错误处理和重试

**缓存**：
- [ ] Redis集成完成
- [ ] 模型查询有缓存
- [ ] 缓存命中率>50%

**认证授权**：
- [ ] JWT认证工作正常
- [ ] RBAC权限控制生效
- [ ] 用户能正常注册登录

**模型管理**：
- [ ] 模型文件可上传
- [ ] 版本管理完整
- [ ] 文件存储在MinIO

**推理服务**：
- [ ] 推理接口可调用
- [ ] 并发控制有效
- [ ] 超时机制正常

### 可演示功能
1. 用户注册登录获取JWT
2. 创建模型并上传文件
3. 发布模型版本
4. 调用推理接口
5. 查看缓存命中率

### 下一步
进入**阶段3：企业级特性**，学习：
- 消息队列（Kafka）
- Kubernetes部署
- 限流熔断
- 分布式事务
- 监控告警

---

**继续学习**：[阶段3文档](../04-phase3/README.md)