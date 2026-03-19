# 订单超时自动关单 - 数据流图 (DFD)

## 订单状态机

```
┌─────────┐    支付成功    ┌─────────┐    完成    ┌─────────┐
│ PENDING │ ─────────────> │  PAID   │ ────────> │ COMPLETE│
│ (待支付) │                │ (已支付) │           │ (已完成) │
└────┬────┘                └─────────┘           └─────────┘
     │                                                  ▲
     │ 超时(15分钟)           取消                       │
     ▼                        ▼                         │
┌──────────┐            ┌───────────┐                    │
│ TIMEOUT  │            │ CANCELLED │ ─────────────────┘
│ (已超时) │            │  (已取消) │
└──────────┘            └───────────┘
```

## Level 1 DFD - 订单超时处理

```
                                    ┌──────────────┐
                                    │   Scheduler  │
                                    │  (定时任务)   │
                                    └──────┬───────┘
                                           │
                    ┌──────────────────────┼──────────────────────┐
                    │                      │                      │
                    ▼                      ▼                      ▼
┌────────┐   创建订单   ┌────────────┐  查询超时   ┌────────────┐ 回滚库存 ┌────────────┐
│  User  │ ──────────> │ Order API  │ ─────────> │Order Service│ ────────> │Inventory API│
└────────┘             └────────────┘            └──────┬───────┘          └────────────┘
                    │                      │           │
                    │                      │           │
                    ▼                      │           ▼
            ┌──────────────┐               │    ┌────────────┐
            │    MySQL     │               │    │   MySQL    │
            │   orders     │               │    │ inventory  │
            │  order_items │               │    └────────────┘
            └──────────────┘               │
                                           │
                                           ▼
                                   ┌──────────────┐
                                   │ order_logs   │
                                   │ (操作日志)   │
                                   └──────────────┘
```

## Level 2 DFD - 超时检查子流程

```
                                    ┌─────────────────────────────┐
                                    │      定时任务 (Scheduler)     │
                                    │   cron: */1 * * * * (每分钟)   │
                                    └──────────────┬──────────────┘
                                                   │
                                                   │ 1. 查询待处理订单
                                                   ▼
┌────────────┐     ┌─────────────────────────────────────────────────────────────────────┐
│            │     │                         Order Service                               │
│            │     │  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐  │
│            │     │  │ 超时查询器    │  │  状态更新器   │  │       库存回滚器           │  │
│            │     │  │              │  │              │  │                           │  │
│            │<────┼──┤ SELECT id    │  │ UPDATE status│  │ UPDATE inventory SET      │  │
│            │     │  │ FROM orders  │─>│ = TIMEOUT    │─>│ stock = stock + quantity   │  │
│            │     │  │ WHERE status  │  │ WHERE id=?   │  │ WHERE product_id=?         │  │
│            │     │  │ = PENDING    │  │ AND version= │  │                             │  │
│            │     │  │ AND expire <= │  │ old_version  │  │ 乐观锁: WHERE version=?    │  │
│            │     │  │ NOW()         │  │              │  │                             │  │
│            │     │  └──────────────┘  └──────────────┘  └──────────────────────────────┘  │
│            │     │            │              │                      │                   │
└────────────┘     └────────────┼──────────────┼──────────────────────┼───────────────────┘
                               │              │                      │
                               ▼              ▼                      ▼
                    ┌─────────────────────────────────────────────────────────────┐
                    │                         MySQL                                  │
                    │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐     │
                    │  │ orders  │  │inventory│  │order_   │  │   order_logs    │     │
                    │  │         │  │         │  │items    │  │                 │     │
                    │  │ id      │  │product_ │  │         │  │ order_id        │     │
                    │  │ status  │  │id       │  │quantity │  │ action: TIMEOUT │     │
                    │  │ expire_ │  │stock    │  │         │  │ reason: 超时关单 │     │
                    │  │ time    │  │version  │  │         │  │ operator: system│     │
                    │  └─────────┘  └─────────┘  └─────────┘  └─────────────────┘     │
                    └─────────────────────────────────────────────────────────────┘
```

## 数据流描述

| 数据流 | 来源 → 目标 | 说明 |
|--------|-------------|------|
| `订单创建请求` | User → Order API | 用户提交购物车，生成待支付订单 |
| `锁定库存请求` | Order API → Inventory API | 扣减库存，记录库存流水 |
| `超时查询请求` | Scheduler → Order Service | 定时任务扫描待支付订单 |
| `更新状态请求` | Order Service → MySQL | 将 PENDING → TIMEOUT |
| `回滚库存请求` | Order Service → Inventory API | 超时/取消时恢复库存 |
| `日志记录` | Order Service → MySQL | 记录操作日志用于审计 |

## 关键表结构变更

### orders 表新增字段

```sql
ALTER TABLE orders ADD COLUMN expire_time DATETIME NOT NULL 
    COMMENT '订单过期时间' AFTER status;

ALTER TABLE orders ADD COLUMN version INT DEFAULT 1 
    COMMENT '乐观锁版本号' AFTER expire_time;

ALTER TABLE orders ADD INDEX idx_expire_status (expire_time, status) 
    COMMENT '超时扫描索引';
```

### order_logs 表（审计日志）

```sql
CREATE TABLE order_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL COMMENT '订单ID',
    action VARCHAR(20) NOT NULL COMMENT '操作类型: TIMEOUT/CANCEL/PAY',
    before_status INT COMMENT '变更前状态',
    after_status INT COMMENT '变更后状态',
    reason VARCHAR(200) COMMENT '变更原因',
    operator VARCHAR(50) DEFAULT 'system' COMMENT '操作者',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
) COMMENT '订单操作日志';
```

## 幂等性保证

```
1. 行锁: SELECT ... FOR UPDATE (防止并发状态变更)
2. 乐观锁: UPDATE ... WHERE version = old_version (防止重复处理)
3. 状态校验: 只处理 PENDING 状态的订单
4. 唯一索引: order_logs 表可加唯一约束防止重复记录
```

## 超时时间配置

```yaml
# config.yaml
order:
  timeout_minutes: 15    # 订单超时时间(分钟)
  scan_interval: 60      # 定时扫描间隔(秒)
```
