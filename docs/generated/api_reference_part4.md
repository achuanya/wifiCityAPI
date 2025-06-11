# WiFi城市API接口文档 - 第四部分

## 优惠券相关API（续）

### 6. 更新优惠券使用限制

更新指定优惠券的使用限制设置。

**请求**:

- 方法: `PATCH`
- 路径: `/coupons/{coupon_id}/limits`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名                | 类型     | 必填  | 描述                           |
|----------------------|----------|------|--------------------------------|
| min_purchase_amount  | decimal  | 否   | 最低消费金额                    |
| usage_limit_per_user | integer  | 否   | 每用户可领取数量，0表示不限     |

**请求示例**:

```json
{
  "min_purchase_amount": 20.00,
  "usage_limit_per_user": 2
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "min_purchase_amount": 20.00,
    "usage_limit_per_user": 2,
    "updated_at": "2023-07-16T19:00:00Z"
  }
}
```

### 7. 更新优惠券发行量

更新指定优惠券的发行量设置。

**请求**:

- 方法: `PATCH`
- 路径: `/coupons/{coupon_id}/quantity`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| total_quantity | integer | 是   | 总发行量，0表示不限             |

**请求示例**:

```json
{
  "total_quantity": 2000
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "total_quantity": 2000,
    "issued_quantity": 152,
    "remaining_quantity": 1848,
    "updated_at": "2023-07-16T19:30:00Z"
  }
}
```

### 8. 更新优惠券适用门店

更新指定优惠券的适用门店。

**请求**:

- 方法: `PATCH`
- 路径: `/coupons/{coupon_id}/store`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名    | 类型    | 必填  | 描述                           |
|----------|---------|------|--------------------------------|
| store_id | integer | 否   | 适用门店ID，null表示全平台通用   |

**请求示例**:

```json
{
  "store_id": 1000001
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "store_id": 1000001,
    "store_name": "星巴克西湖旗舰店",
    "updated_at": "2023-07-16T20:00:00Z"
  }
}
```

### 9. 更新优惠券状态

更新指定优惠券的状态（启用/禁用）。

**请求**:

- 方法: `PATCH`
- 路径: `/coupons/{coupon_id}/status`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名   | 类型    | 必填  | 描述                           |
|---------|---------|------|--------------------------------|
| status  | integer | 是   | 状态：1启用，0禁用              |

**请求示例**:

```json
{
  "status": 0
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "status": 0,
    "updated_at": "2023-07-16T20:30:00Z"
  }
}
```

### 10. 删除优惠券

删除指定的优惠券。

**请求**:

- 方法: `DELETE`
- 路径: `/coupons/{coupon_id}`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "deleted_at": "2023-07-16T21:00:00Z"
  }
}
```

### 11. 批量创建优惠券

批量创建多个优惠券。

**请求**:

- 方法: `POST`
- 路径: `/coupons/batch`
- 内容类型: `application/json`

**请求参数**:

| 参数名    | 类型    | 必填  | 描述                           |
|----------|---------|------|--------------------------------|
| coupons  | array   | 是   | 优惠券数组                      |

每个优惠券的参数与"创建优惠券"接口相同。

**请求示例**:

```json
{
  "coupons": [
    {
      "coupon_name": "折扣券8折",
      "coupon_type": "DISCOUNT",
      "value": 0.8,
      "min_purchase_amount": 50.00,
      "usage_limit_per_user": 1,
      "total_quantity": 500,
      "start_time": "2023-07-17T00:00:00Z",
      "end_time": "2023-08-17T23:59:59Z",
      "description": "8折优惠券，满50元可使用"
    },
    {
      "coupon_name": "运费券",
      "coupon_type": "SHIPPING",
      "value": 5.00,
      "min_purchase_amount": 0.00,
      "usage_limit_per_user": 1,
      "total_quantity": 500,
      "start_time": "2023-07-17T00:00:00Z",
      "end_time": "2023-08-17T23:59:59Z",
      "description": "免邮券，任意消费可使用"
    }
  ]
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupons": [
      {
        "coupon_id": 1002,
        "coupon_name": "折扣券8折",
        "coupon_type": "DISCOUNT",
        "value": 0.8,
        "created_at": "2023-07-16T21:30:00Z"
      },
      {
        "coupon_id": 1003,
        "coupon_name": "运费券",
        "coupon_type": "SHIPPING",
        "value": 5.00,
        "created_at": "2023-07-16T21:30:00Z"
      }
    ],
    "total_created": 2
  }
}
```

### 12. 查询门店可用优惠券列表

获取指定门店可用的优惠券列表。

**请求**:

- 方法: `GET`
- 路径: `/stores/{store_id}/coupons`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| coupon_type    | string  | 否   | 按优惠券类型筛选                |
| page           | integer | 否   | 页码，默认1                     |
| page_size      | integer | 否   | 每页记录数，默认20              |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "coupon_id": 1001,
        "coupon_name": "新用户专享立减10元",
        "coupon_type": "CASH",
        "value": 10.00,
        "min_purchase_amount": 20.00,
        "start_time": "2023-07-17T00:00:00Z",
        "end_time": "2023-09-17T23:59:59Z",
        "total_quantity": 2000,
        "remaining_quantity": 1848,
        "description": "新用户首次扫码连接WiFi即可领取，满20元可使用，有效期30天"
      },
      // 更多优惠券...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

### 13. 查询用户可领取的优惠券列表

获取指定用户可领取的优惠券列表。

**请求**:

- 方法: `GET`
- 路径: `/users/{user_union_id}/available-coupons`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| store_id       | integer | 否   | 按门店ID筛选                    |
| coupon_type    | string  | 否   | 按优惠券类型筛选                |
| page           | integer | 否   | 页码，默认1                     |
| page_size      | integer | 否   | 每页记录数，默认20              |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "coupon_id": 1002,
        "coupon_name": "折扣券8折",
        "coupon_type": "DISCOUNT",
        "value": 0.8,
        "min_purchase_amount": 50.00,
        "start_time": "2023-07-17T00:00:00Z",
        "end_time": "2023-08-17T23:59:59Z",
        "description": "8折优惠券，满50元可使用",
        "store_id": null,
        "store_name": "全平台通用"
      },
      // 更多优惠券...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 2,
      "total_pages": 1
    }
  }
}
```

## 优惠券日志相关API

### 1. 用户领取优惠券

记录用户领取优惠券的操作。

**请求**:

- 方法: `POST`
- 路径: `/users/{user_union_id}/coupons/{coupon_id}/receive`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |
| coupon_id      | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名    | 类型    | 必填  | 描述                           |
|----------|---------|------|--------------------------------|
| store_id | integer | 否   | 领取门店ID                      |
| remark   | string  | 否   | 备注信息                        |

**请求示例**:

```json
{
  "store_id": 1000001,
  "remark": "首次连接WiFi赠送"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log_id": 10001,
    "coupon_id": 1002,
    "coupon_name": "折扣券8折",
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "store_id": 1000001,
    "action_type": "RECEIVE",
    "action_time": "2023-07-17T09:00:00Z",
    "valid_until": "2023-08-17T23:59:59Z",
    "status": 1
  }
}
```

### 2. 用户使用优惠券

记录用户使用优惠券的操作。

**请求**:

- 方法: `POST`
- 路径: `/users/{user_union_id}/coupons/{coupon_id}/use`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |
| coupon_id      | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名            | 类型    | 必填  | 描述                           |
|------------------|---------|------|--------------------------------|
| store_id         | integer | 是   | 使用门店ID                      |
| order_id         | string  | 是   | 关联订单ID                      |
| amount_deducted  | decimal | 是   | 优惠券抵扣金额                  |
| remark           | string  | 否   | 备注信息                        |

**请求示例**:

```json
{
  "store_id": 1000001,
  "order_id": "ORD2023071700001",
  "amount_deducted": 10.00,
  "remark": "购买咖啡使用"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log_id": 10002,
    "coupon_id": 1001,
    "coupon_name": "新用户专享立减10元",
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "store_id": 1000001,
    "action_type": "USE",
    "action_time": "2023-07-17T10:00:00Z",
    "order_id": "ORD2023071700001",
    "amount_deducted": 10.00,
    "status": 1
  }
}
```

### 3. 记录优惠券发放

系统管理员手动发放优惠券的操作记录。

**请求**:

- 方法: `POST`
- 路径: `/coupons/{coupon_id}/issue`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |
| store_id       | integer | 否   | 发放门店ID                      |
| remark         | string  | 否   | 备注信息                        |

**请求示例**:

```json
{
  "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
  "store_id": 1000001,
  "remark": "客户投诉补偿"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log_id": 10003,
    "coupon_id": 1003,
    "coupon_name": "运费券",
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "store_id": 1000001,
    "action_type": "ISSUE",
    "action_time": "2023-07-17T11:00:00Z",
    "status": 1,
    "remark": "客户投诉补偿"
  }
}
```

### 4. 记录优惠券过期

标记优惠券为已过期状态。

**请求**:

- 方法: `POST`
- 路径: `/users/{user_union_id}/coupons/{coupon_id}/expire`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |
| coupon_id      | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名    | 类型    | 必填  | 描述                           |
|----------|---------|------|--------------------------------|
| remark   | string  | 否   | 备注信息                        |

**请求示例**:

```json
{
  "remark": "用户未在有效期内使用"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log_id": 10004,
    "coupon_id": 1003,
    "coupon_name": "运费券",
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "action_type": "EXPIRE",
    "action_time": "2023-07-17T12:00:00Z",
    "status": 1
  }
}
```

### 5. 记录优惠券退回

记录优惠券退回操作。

**请求**:

- 方法: `POST`
- 路径: `/users/{user_union_id}/coupons/{coupon_id}/refund`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |
| coupon_id      | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名    | 类型    | 必填  | 描述                           |
|----------|---------|------|--------------------------------|
| order_id | string  | 是   | 关联订单ID                      |
| remark   | string  | 否   | 备注信息                        |

**请求示例**:

```json
{
  "order_id": "ORD2023071700001",
  "remark": "订单取消，退回优惠券"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log_id": 10005,
    "coupon_id": 1001,
    "coupon_name": "新用户专享立减10元",
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "action_type": "REFUND",
    "action_time": "2023-07-17T13:00:00Z",
    "order_id": "ORD2023071700001",
    "status": 1,
    "remark": "订单取消，退回优惠券"
  }
}
```

### 6. 查询用户已领取优惠券列表

获取用户已领取的优惠券列表。

**请求**:

- 方法: `GET`
- 路径: `/users/{user_union_id}/coupons`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| status     | string  | 否   | 状态筛选：AVAILABLE/USED/EXPIRED |
| page       | integer | 否   | 页码，默认1                     |
| page_size  | integer | 否   | 每页记录数，默认20              |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "coupon_id": 1002,
        "coupon_name": "折扣券8折",
        "coupon_type": "DISCOUNT",
        "value": 0.8,
        "min_purchase_amount": 50.00,
        "received_time": "2023-07-17T09:00:00Z",
        "valid_until": "2023-08-17T23:59:59Z",
        "status": "AVAILABLE",
        "store_id": null,
        "store_name": "全平台通用"
      },
      {
        "coupon_id": 1001,
        "coupon_name": "新用户专享立减10元",
        "coupon_type": "CASH",
        "value": 10.00,
        "min_purchase_amount": 20.00,
        "received_time": "2023-07-16T15:00:00Z",
        "used_time": "2023-07-17T10:00:00Z",
        "status": "USED",
        "store_id": 1000001,
        "store_name": "星巴克西湖旗舰店",
        "order_id": "ORD2023071700001"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 2,
      "total_pages": 1
    },
    "summary": {
      "available": 1,
      "used": 1,
      "expired": 0,
      "total": 2
    }
  }
}
```

### 7. 查询指定优惠券的领取详情

获取指定优惠券的领取详情记录。

**请求**:

- 方法: `GET`
- 路径: `/coupons/{coupon_id}/receive-logs`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| start_time     | string  | 否   | 开始时间，ISO8601格式            |
| end_time       | string  | 否   | 结束时间，ISO8601格式            |
| page           | integer | 否   | 页码，默认1                     |
| page_size      | integer | 否   | 每页记录数，默认20              |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1002,
    "coupon_name": "折扣券8折",
    "list": [
      {
        "log_id": 10001,
        "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
        "user_nickname": "张三",
        "store_id": 1000001,
        "store_name": "星巴克西湖旗舰店",
        "action_type": "RECEIVE",
        "action_time": "2023-07-17T09:00:00Z",
        "status": 1
      },
      // 更多记录...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 45,
      "total_pages": 3
    },
    "summary": {
      "total_received": 45,
      "total_used": 12,
      "total_expired": 5,
      "total_refunded": 0,
      "total_available": 28
    }
  }
}
```

### 8. 查询指定优惠券的使用详情

获取指定优惠券的使用详情记录。

**请求**:

- 方法: `GET`
- 路径: `/coupons/{coupon_id}/use-logs`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| store_id       | integer | 否   | 按门店ID筛选                    |
| start_time     | string  | 否   | 开始时间，ISO8601格式            |
| end_time       | string  | 否   | 结束时间，ISO8601格式            |
| page           | integer | 否   | 页码，默认1                     |
| page_size      | integer | 否   | 每页记录数，默认20              |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "coupon_name": "新用户专享立减10元",
    "list": [
      {
        "log_id": 10002,
        "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
        "user_nickname": "张三",
        "store_id": 1000001,
        "store_name": "星巴克西湖旗舰店",
        "action_type": "USE",
        "action_time": "2023-07-17T10:00:00Z",
        "order_id": "ORD2023071700001",
        "amount_deducted": 10.00,
        "status": 1
      },
      // 更多记录...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 32,
      "total_pages": 2
    },
    "summary": {
      "total_used_count": 32,
      "total_deducted_amount": 320.00,
      "average_deducted_amount": 10.00
    }
  }
}
```

### 9. 查询门店优惠券核销记录

获取指定门店的优惠券核销记录。

**请求**:

- 方法: `GET`
- 路径: `/stores/{store_id}/coupon-usage`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| coupon_id      | integer | 否   | 按优惠券ID筛选                  |
| coupon_type    | string  | 否   | 按优惠券类型筛选                |
| start_time     | string  | 否   | 开始时间，ISO8601格式            |
| end_time       | string  | 否   | 结束时间，ISO8601格式            |
| page           | integer | 否   | 页码，默认1                     |
| page_size      | integer | 否   | 每页记录数，默认20              |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "store_name": "星巴克西湖旗舰店",
    "list": [
      {
        "log_id": 10002,
        "coupon_id": 1001,
        "coupon_name": "新用户专享立减10元",
        "coupon_type": "CASH",
        "user_nickname": "张三",
        "action_time": "2023-07-17T10:00:00Z",
        "order_id": "ORD2023071700001",
        "amount_deducted": 10.00
      },
      // 更多记录...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 85,
      "total_pages": 5
    },
    "summary": {
      "total_used_count": 85,
      "total_deducted_amount": 712.50,
      "by_coupon_type": {
        "CASH": {
          "count": 45,
          "amount": 450.00
        },
        "DISCOUNT": {
          "count": 32,
          "amount": 237.50
        },
        "SHIPPING": {
          "count": 8,
          "amount": 25.00
        }
      }
    }
  }
}
```

## 数据统计与报表API

### 1. 门店统计

#### 1.1 统计门店总数

获取系统中的门店总数统计信息。

**请求**:

- 方法: `GET`
- 路径: `/stats/stores/count`

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| status     | integer | 否   | 按状态筛选：1正常，0停用        |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "total_stores": 128,
    "active_stores": 120,
    "inactive_stores": 8
  }
}
```

#### 1.2 统计不同省份/城市门店数量

获取不同地区的门店分布统计。

**请求**:

- 方法: `GET`
- 路径: `/stats/stores/distribution`

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| level      | string  | 否   | 统计级别：province/city/district，默认province |
| status     | integer | 否   | 按状态筛选：1正常，0停用        |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "level": "province",
    "distribution": [
      {
        "name": "浙江省",
        "count": 42,
        "percentage": 32.8
      },
      {
        "name": "上海市",
        "count": 35,
        "percentage": 27.3
      },
      {
        "name": "江苏省",
        "count": 30,
        "percentage": 23.4
      },
      // 更多省份...
    ],
    "total": 128
  }
}
```

#### 1.3 按日/周/月统计新增门店数

获取指定时间段内新增门店的统计数据。

**请求**:

- 方法: `GET`
- 路径: `/stats/stores/growth`

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| start_date | string  | 否   | 开始日期，格式YYYY-MM-DD，默认30天前 |
| end_date   | string  | 否   | 结束日期，格式YYYY-MM-DD，默认今天  |
| group_by   | string  | 否   | 分组方式：day/week/month，默认day |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "group_by": "day",
    "time_range": {
      "start": "2023-06-17",
      "end": "2023-07-17"
    },
    "trend": [
      {
        "period": "2023-06-17",
        "new_stores": 2
      },
      {
        "period": "2023-06-18",
        "new_stores": 0
      },
      // 更多日期...
    ],
    "summary": {
      "total_new_stores": 32,
      "average_daily_growth": 1.07
    }
  }
}
``` 