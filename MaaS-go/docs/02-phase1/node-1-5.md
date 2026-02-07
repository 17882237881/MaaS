# 节点1.5：数据库层设计

> 📅 **学习时间**：5天  
> 🎯 **目标**：设计并实现数据库层，包括模型定义、Repository模式和迁移系统

## 本节你将学到

1. 数据库设计原则（范式、索引、关系）
2. GORM ORM框架使用
3. Repository设计模式
4. 数据库迁移系统
5. 事务管理和连接池

---

## 技术详解

### 1. 数据库设计原则

**第一范式（1NF）**：原子性
- 每个字段都是不可分割的原子值
- 反例：`hobbies: "reading,swimming,gaming"`
- 正例：拆分到单独的hobbies表

**第二范式（2NF）**：完全依赖
- 非主键字段必须完全依赖于主键
- 消除部分依赖

**第三范式（3NF）**：消除传递依赖
- 非主键字段之间不能相互依赖

**我们的设计**：
- 使用UUID作为主键（分布式友好）
- 适当的反范式（如租户配额嵌入）
- 软删除（DeletedAt）支持数据恢复

### 2. GORM简介

GORM是Go语言最流行的ORM库，提供：
- 模型定义和自动迁移
- CRUD操作
- 关联（一对一、一对多、多对多）
- 钩子（BeforeCreate, AfterUpdate等）
- 事务支持

**模型定义示例**：
```go
type Model struct {
    ID        string         `gorm:"type:uuid;primary_key"`
    Name      string         `gorm:"type:varchar(255);not null"`
    CreatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 3. Repository模式

**为什么要用Repository？**
- 业务逻辑与数据访问解耦
- 易于测试（可Mock）
- 支持切换数据库实现

**结构**：
```
Service → Repository → Database
   ↑           ↑           ↑
 业务逻辑    数据访问    具体实现
```

### 4. 数据库迁移

**什么是迁移？**
管理数据库Schema的变更，包括：
- 创建表
- 添加字段
- 修改字段类型
- 创建索引

**GORM AutoMigrate**：
```go
db.AutoMigrate(&Model{}, &User{}, &Tag{})
```

**注意**：AutoMigrate不会删除字段，生产环境建议使用专业的迁移工具。

---

## 实操任务

### 任务1：设计数据库模型

创建以下模型：
1. **User** - 用户表
2. **Tenant** - 租户表
3. **Model** - 模型表
4. **Tag** - 标签表
5. **Metadata** - 元数据表
6. **ModelVersion** - 模型版本表

### 任务2：实现Repository

创建ModelRepository接口和实现：
- Create, GetByID, List, Update, Delete
- AddTags, RemoveTags
- SetMetadata, GetMetadata

### 任务3：数据库迁移

实现AutoMigrate函数，自动创建所有表。

### 任务4：更新Model Registry

更新main.go：
- 连接数据库
- 运行迁移
- 使用Repository操作数据

---

## 代码变更记录

### 提交信息
```
feat(phase1/node1.5): implement database layer design

- Design database schema for MaaS platform
- Implement GORM models with relationships
- Create Repository pattern for data access
- Add database migration system
- Update Model Registry service to use database
```

### 创建的文件

#### 1. model-registry/internal/model/user.go
**新增文件**
用户和租户模型：
- **User**：用户表
  - ID, Username, Email, Password, Role, Status
  - TenantID（外键）
  - 软删除支持

- **Tenant**：租户表
  - ID, Name, Description, Status
  - Quota（嵌入结构体）
  - 包含资源配额限制

#### 2. model-registry/internal/model/model.go
**新增文件**
模型相关表：
- **Model**：模型主表
  - ID, Name, Description, Version, Framework, Status
  - Size, Checksum, StoragePath, DockerImage
  - OwnerID, TenantID（所有权）
  - IsPublic（可见性）
  - Tags, Metadata, Versions（关联）

- **Tag**：标签表
  - ID, Name
  - 多对多关联Model

- **Metadata**：元数据表
  - ID, ModelID, Key, Value
  - 键值对存储

- **ModelVersion**：版本表
  - ID, ModelID, Version, Status
  - ChangeLog（变更说明）

#### 3. model-registry/internal/repository/model_repository.go
**新增文件**
ModelRepository实现：
- **接口定义**：ModelRepository
- **实现**：GormModelRepository
- **方法**：
  - Create：创建模型（检查重复）
  - GetByID：按ID查询（预加载Tags和Metadata）
  - GetByNameAndVersion：按名称版本查询
  - List：列表查询（支持过滤和分页）
  - Update：更新模型
  - Delete：软删除
  - UpdateStatus：更新状态
  - AddTags/RemoveTags：标签管理
  - SetMetadata/GetMetadata：元数据管理

#### 4. model-registry/internal/repository/database.go
**新增文件**
数据库连接和迁移：
- **NewDatabase**：创建数据库连接
- **AutoMigrate**：自动迁移所有表
- 使用PostgreSQL驱动

#### 5. model-registry/internal/service/model_service.go
**新增文件**
业务逻辑层：
- **ModelService接口**：定义业务方法
- **modelService实现**：
  - CreateModel：创建模型（验证框架、添加标签元数据）
  - GetModel：获取模型
  - ListModels：列表查询
  - UpdateModel：更新模型
  - DeleteModel：删除模型
  - UpdateModelStatus：更新状态

### 修改的文件

#### model-registry/cmd/main.go
**大幅更新**
- 添加数据库连接
- 运行AutoMigrate
- 初始化Repository
- 初始化Service
- 健康检查包含数据库状态

---

## 数据库Schema

### ERD（实体关系图）

```
[User] 1--* [Model] *--* [Tag]
  |          |
  |          *--* [Metadata]
  |          |
  |          *--* [ModelVersion]
  |
[Tenant] 1--*
```

### 表结构

**users表**：
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'developer',
    status VARCHAR(20) DEFAULT 'active',
    tenant_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

**models表**：
```sql
CREATE TABLE models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(50) NOT NULL,
    framework VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    size BIGINT DEFAULT 0,
    checksum VARCHAR(64),
    storage_path VARCHAR(512),
    docker_image VARCHAR(255),
    owner_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
```

---

## 验证步骤

### 1. 数据库连接验证

```bash
# 1. 启动PostgreSQL
docker-compose up -d postgres

# 2. 运行Model Registry
cd model-registry
go run cmd/main.go

# 3. 查看日志，应该显示：
# {"msg":"Running database migrations..."}
# {"msg":"Database migrations completed"}
```

### 2. 表结构验证

```bash
# 连接PostgreSQL
docker exec -it postgres psql -U postgres -d maas_platform

# 查看表
\dt

# 应该看到：
#  models
#  model_metadata
#  model_tags
#  model_versions
#  tags
#  tenants
#  users
```

### 3. API验证

```bash
# 创建模型
curl -X POST http://localhost:8081/api/v1/models \
  -H "Content-Type: application/json" \
  -d '{
    "name": "bert-base",
    "description": "BERT base model",
    "version": "1.0.0",
    "framework": "pytorch",
    "owner_id": "user-123",
    "tenant_id": "tenant-456"
  }'

# 查询模型列表
curl http://localhost:8081/api/v1/models

# 查询单个模型
curl http://localhost:8081/api/v1/models/{id}
```

---

## 检查清单

完成本节点后，请确认：

- [ ] 数据库模型定义完整
- [ ] GORM可以自动创建表
- [ ] Repository接口完整实现
- [ ] CRUD操作正常工作
- [ ] 关联查询（Tags、Metadata）正常
- [ ] 软删除功能正常
- [ ] 分页查询正常
- [ ] 健康检查包含数据库状态

---

## 最佳实践

### 1. 主键选择
使用UUID而非自增ID：
- 分布式系统友好
- 避免暴露数据量
- 可以预生成

### 2. 索引设计
- 外键自动创建索引
- 查询频繁的字段加索引
- 避免过多索引影响写入

### 3. 连接池配置
```go
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(100)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(time.Hour)
```

---

## 下一步

完成本节点后，你已经实现了完整的数据库层。阶段1（基础架构搭建）完成！🎉

**阶段1里程碑**：
- ✅ API Gateway（Gin框架、中间件、日志监控）
- ✅ 配置管理体系
- ✅ Model Registry（数据库层、CRUD操作）

**进入阶段2：核心功能开发** → [查看阶段2文档](../03-phase2/README.md)

在阶段2中，你将学习：
- gRPC服务间通信
- Redis缓存
- JWT认证
- 模型上传和存储

---

## 参考资源

- [GORM官方文档](https://gorm.io/docs/)
- [PostgreSQL官方文档](https://www.postgresql.org/docs/)
- [Repository模式](https://martinfowler.com/eaaCatalog/repository.html)
- [数据库设计范式](https://en.wikipedia.org/wiki/Database_normalization)
