# 节点2.2：Redis 缓存层设计

> 📅 **学习时间**：4天  
> 🎯 **目标**：为 Model Registry 添加缓存层，降低数据库压力并提升响应速度

## 本节你将学到

1. Redis 基础数据类型与适用场景
2. Cache-Aside 缓存模式
3. 缓存穿透/击穿/雪崩的应对策略
4. 多级缓存（可选）设计

---

## 技术详解

### 1. Redis 是什么？

Redis 是内存中的数据结构存储系统，支持 String/Hash/List/Set/Sorted Set 等数据结构。

**适用场景**：
- 热点数据缓存
- 计数器、排行榜
- 分布式锁

### 2. 缓存模式（Cache-Aside）

```
请求 → 缓存命中 → 返回
       缓存未命中 → 查数据库 → 回填缓存 → 返回
```

### 3. 常见问题

- **缓存穿透**：请求不存在的数据
  - 解决：缓存空值、布隆过滤器
- **缓存击穿**：热点 key 失效瞬间
  - 解决：互斥锁、逻辑过期
- **缓存雪崩**：大量 key 同时过期
  - 解决：随机过期时间、分级缓存

---

## 实操任务（建议实现顺序）

1. 安装并配置 Redis
2. 编写 Redis 连接封装（`redis-py`）
3. 为 `GetModel` 添加缓存逻辑
4. 为 `ListModels` 添加缓存（带分页 / filter）
5. 设计缓存 key 规则：
   - 单条：`model:{id}`
   - 列表：`models:{page}:{limit}:{framework}:{status}`
6. 增加缓存失效策略（更新/删除模型时清理）

---

## Python 示例（Cache-Aside）

```python
import json
from redis.asyncio import Redis

async def get_model_cached(redis: Redis, repo, model_id: str):
    key = f"model:{model_id}"

    cached = await redis.get(key)
    if cached:
        return json.loads(cached)

    model = await repo.get_by_id(model_id)
    if model is None:
        await redis.set(key, "null", ex=60)  # 缓存空值
        return None

    await redis.set(key, json.dumps(model), ex=3600)
    return model
```

---

## 检查清单

- [ ] Redis 连接正常
- [ ] 缓存命中逻辑正确
- [ ] 缓存过期策略合理
- [ ] 更新/删除后缓存能正确失效
- [ ] 缓存穿透/击穿处理完成

---

## 下一步

完成本节点后，继续进入：

**节点2.3：认证与授权** → [回到阶段2目录](./README.md)
