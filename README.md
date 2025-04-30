# 秒杀系统开发文档

## 1. 所要使用的技术

### 编程语言与框架
- **Go语言**：具有高效、并发性能好等特点，适合处理高并发场景。
- **Hertz Web框架**：用于快速搭建Web服务，处理HTTP请求和路由。

### 数据库
- **MySQL**：用于存储商品信息、用户信息、订单信息等持久化数据。
- **Redis**：内存数据库，用于缓存商品库存、秒杀活动状态等高频访问数据，以提高系统性能和响应速度。

### 缓存
- **Go标准库中的缓存机制**：如 `sync.Map` 等，用于在内存中进行简单的缓存。
- **第三方缓存库**：如 `go-cache`，提供更丰富的缓存功能和策略。

### 消息队列
- **Kafka**：用于异步处理秒杀订单，将订单消息发送到消息队列中，由消费者进行异步处理，避免高并发时直接操作数据库导致的性能问题。

### 分布式锁
- **基于Redis的分布式锁**：通过Redis的 `SETNX` 等命令实现分布式锁，确保在同一时刻只有一个请求能够进行商品库存的扣减等操作，防止超卖。

### 限流与熔断
- **Golang的令牌桶算法实现**：用于限制每秒的请求数，防止系统被过多的请求压垮。
- **Hystrix或类似的熔断框架**：当系统出现故障或响应时间过长时，自动熔断相关服务，避免级联故障。

### 分布式系统协调
- **etcd或Consul**：用于服务发现和配置管理，帮助各个微服务实例能够动态地发现彼此，并获取最新的配置信息。

### 日志与监控
- **Hlog**：Hertz的日志库，用于记录系统的运行日志，方便排查问题。
- **Prometheus和Grafana**：用于监控系统的各项指标，如CPU使用率、内存使用率、请求响应时间、每秒请求数等，并通过可视化界面展示。

## 2. 项目结构构建

```plaintext
seckill-backend/
├── cmd/                     # 项目入口（main.go）
│   └── main.go
├── router/                  # Hertz 路由定义
│   └── router.go
├── controller/              # 控制器，处理请求
│   ├── seckill.go
│   ├── order.go
│   └── product.go
├── service/                 # 核心业务逻辑
│   ├── seckill.go
│   ├── order.go
│   └── product.go
├── dao/                     # 数据访问层（封装 GORM 操作）
│   ├── product.go
│   ├── order.go
│   ├── migrate.go
│   └── init.go
├── model/                   # 数据模型定义（GORM）
│   ├── product.go
│   └── order.go
├── middleware/              # 自定义中间件
│   ├── ratelimit.go
│   ├── circuitbreaker.go
│   └── logger.go
├── lock/                    # Redis 分布式锁实现
│   └── redis_lock.go
├── kafka/                   # Kafka 生产者/消费者
│   ├── producer.go
│   └── consumer.go
├── utils/                   # 工具类
│   ├── response.go
│   ├── log.go
│   └── errors.go
├── go.mod
└── go.sum
```
# 3.数据库设计

## 1. `Order` 表（订单信息）

| 字段       | 类型        | 描述                     |
|------------|-------------|--------------------------|
| `OrderID`  | `varchar(64)`| 订单 ID，主键            |
| `UserID`   | `int`       | 用户 ID                  |
| `Orders`   | `[]OrderItem`| 订单项，关联 `OrderItem` 表 |
| `Address`  | `varchar`   | 收货地址                  |
| `Total`    | `float64`   | 订单总额                  |

### 外键关系：
- `Orders` 字段与 `OrderItem` 表的 `OrderID` 关联。

---

## 2. `OrderItem` 表（订单项信息）

| 字段       | 类型        | 描述                     |
|------------|-------------|--------------------------|
| `ID`       | `int`       | 主键                     |
| `OrderID`  | `varchar`   | 订单 ID，外键，关联 `Order` 表 |
| `ProductID`| `int`       | 商品 ID，外键，关联 `Product` 表 |
| `Quantity` | `int`       | 商品数量                  |

### 外键关系：
- `OrderID` 字段与 `Order` 表的 `OrderID` 关联。
- `ProductID` 字段与 `Product` 表的 `ID` 关联。

---

## 3. `OrderMessage` 表（秒杀订单消息）

| 字段       | 类型        | 描述                     |
|------------|-------------|--------------------------|
| `UserID`   | `int`       | 用户 ID                  |
| `ProductID`| `int`       | 商品 ID                  |
| `Quantity` | `int`       | 购买数量                  |

---

## 4. `Product` 表（商品信息）

| 字段          | 类型          | 描述                          |
|---------------|---------------|-------------------------------|
| `ID`          | `int`         | 商品 ID，主键，自增            |
| `Name`        | `varchar(100)`| 商品名称                       |
| `Description` | `text`        | 商品描述                       |
| `Type`        | `varchar(50)` | 商品类型                       |
| `Price`       | `float64`     | 商品价格                       |
| `Stock`       | `int`         | 商品总库存                     |
| `StockSurplus`| `int`         | 剩余库存                       |
| `LimitPerUser`| `int`         | 每个用户的购买限制            |
| `Status`      | `int`         | 商品状态，1: 上架，0: 下架      |
| `StartTime`   | `time.Time`   | 秒杀开始时间（时间戳）         |
| `EndTime`     | `time.Time`   | 秒杀结束时间（时间戳）         |
| `CreatedAt`   | `time.Time`   | 创建时间                       |
| `UpdatedAt`   | `time.Time`   | 更新时间                       |

---

## 外键关系

- **`OrderItem` 表**：
  - `OrderID` 关联 **`Order` 表** 的 `OrderID` 字段。
  - `ProductID` 关联 **`Product` 表** 的 `ID` 字段。

---

## 数据库关系图

- **订单与订单项**：一个订单可以有多个订单项，`Order` 表和 `OrderItem` 表通过 `OrderID` 字段关联。
- **订单项与商品**：每个订单项对应一个商品，`OrderItem` 表的 `ProductID` 字段关联 `Product` 表的 `ID` 字段。

---
# 4.API 接口文档

## 1. 秒杀接口

### 1.1 秒杀请求

**接口路径**：`POST /api/seckill`

**描述**：该接口用于用户参与秒杀活动，提交秒杀请求。

### 请求参数

| 参数        | 类型    | 描述                              | 必填 |
|-------------|---------|-----------------------------------|------|
| `user_id`   | `int`   | 用户 ID                           | 是   |
| `product_id`| `int`   | 商品 ID                           | 是   |

### 请求示例

```json
{
  "user_id": 123,
  "product_id": 456
}
```
### 返回参数
| 参数         | 类型       | 描述     |
|------------|----------|--------|
| `status`   | `int`    | 秒杀请求状态 |
| `message`  | `string` | 状态描述   |

 ### 返回示例
```json
{
  "status": 1,
  "message": "秒杀成功"
}
```
## 2.创建商品接口
### 创建商品
**接口路径**：`POST /api/produc`

**描述**：该接口用于创建新的商品信息，用于秒杀活动中。
###  请求参数
| 参数 	  | 类型	     | 描述	                    | 必填  |
|-------|---------|------------------------|-----|
| name	 | string	 | 商品名称	                  | 是   |
|description	| string	 | 商品描述                   |	否|
|type	| string  | 	商品类型                  |	否|
|price	| float64 | 	商品价格	                 |是|
|stock	| int	    | 商品总库存                  |	是|
|stock_surplus	| int	    | 剩余库存	                  | 是                      |
|limit_per_user	| int	    | 每个用户的购买限制	             |否| 
|status	| int	   |   商品状态：                 1 上架，0 下架	             | 否       |
|start_time	| string	 | 秒杀开始时间（ISO 8601 时间格式）  |	是|
|end_time	| string	 | 秒杀结束时间（ISO 8601 时间格式）	 | 是       |
### 请求示例
```json
{
  "name": "秒杀商品 A",
  "description": "这是秒杀商品 A 的描述",
  "type": "电子产品",
  "price": 199.99,
  "stock": 1000,
  "stock_surplus": 1000,
  "limit_per_user": 1,
  "status": 1,
  "start_time": "2025-05-01T00:00:00Z",
  "end_time": "2025-05-01T01:00:00Z"
}
```
### 返回参数
|参数	| 类型	                  | 描述               |
|----|----------------------|------------------|
|status	| int	| 商品创建状态：1 成功，0 失败 |
|message	| string	              | 状态描述             |
|product_id	|int| 	创建的商品 ID（成功时返回） |
### 返回示例
```json
{
  "status": 1,
  "message": "商品创建成功",
  "product_id": 789
}
```

# 5.日志记录
## 使用了Hertz框架自带的日志记录工具Hlog
### 1.日志中间件:
请求相关信息记录：如记录用户信息、请求时间和请求来源等。部分系统处理信息记录：记录处理耗时和服务器节点信息。
### 2.utils层日志记录:
秒杀过程和结果信息记录，记录商品信息、库存变化、排队信息、秒杀结果和订单信息等。错误信息记录，虽然中间件可以捕获一些通用的错误，如请求格式错误等

# 6.部署说明 ：
使用docker拉取并部署redis kafka zookeeper  

# 7.关于如何处理高并发，防止超卖少卖等问题

1.数据操作采用原子操作，保证数据的一致性，防止秒杀成功却库存不变的情况等

2.redis分布式锁，防止在高并发情况下多个用户同时发送请求导致的超卖。同时可以对高并发进行限流和限制

3.kafka：
缓冲和削峰:在一些业务场景中，请求流量可能会出现高峰和低谷。Kafka可以作为缓冲区，在流量高峰时存储大量消息，然后在流量低谷时让消费者以适当的速度处理这些消息，避免系统因瞬时高并发而崩溃。 

异步处理:Kafka支持异步消息处理，生产者发送消息后不需要等待消费者处理完成，可以继续执行其他任务，提高了系统的整体性能和响应速度

4.使用限流中间价，防止同一用户同时发送大量请求
