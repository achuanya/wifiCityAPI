# WiFi城市API接口文档 - 第三部分

## 扫码日志相关API

### 1. 记录用户扫码连接日志

记录用户扫描二维码连接WiFi的日志信息。

**请求**:

- 方法: `POST`
- 路径: `/scan-logs`
- 内容类型: `application/json`

**请求参数**:

| 参数名             | 类型    | 必填  | 描述                           |
|-------------------|---------|------|--------------------------------|
| store_id          | integer | 是   | 门店ID                         |
| user_union_id     | string  | 是   | 用户微信UnionID                |
| device_info       | string  | 否   | 用户设备信息                    |
| ip_address        | string  | 否   | 用户IP地址                      |
| network_type      | string  | 否   | 网络类型：WIFI/5G/4G/3G/2G/UNKNOWN |
| location_lat      | decimal | 否   | 用户扫码纬度                    |
| location_lng      | decimal | 否   | 用户扫码经度                    |
| mini_program_version | string | 否 | 小程序版本号                   |
| success_flag      | boolean | 否   | 是否成功连接WiFi，默认false      |
| fail_reason_code  | string  | 否   | 连接失败错误码                  |
| fail_reason_message | string | 否  | 连接失败详细信息                |
| wifi_ssid         | string  | 否   | 连接的WiFi名称                  |
| wifi_mac          | string  | 否   | 连接WiFi的MAC地址               |
| wifi_signal       | integer | 否   | WiFi信号强度                    |
| qr_code_type      | string  | 否   | 二维码类型：STORE/EVENT/POSTER/DESK/OTHER |
| qr_code_id        | string  | 否   | 二维码ID                        |
| system_info       | string  | 否   | 操作系统信息                    |
| brand             | string  | 否   | 设备品牌                        |
| model             | string  | 否   | 设备型号                        |
| page_path         | string  | 否   | 扫码来源页路径                  |
| referer           | string  | 否   | 扫码来源URL或分享来源           |
| remark            | string  | 否   | 备注信息                        |

**请求示例**:

```json
{
  "store_id": 1000001,
  "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
  "device_info": "iPhone14,2",
  "ip_address": "192.168.1.100",
  "network_type": "4G",
  "location_lat": 30.246001,
  "location_lng": 120.127002,
  "mini_program_version": "1.0.2",
  "success_flag": true,
  "wifi_ssid": "Starbucks-Guest",
  "wifi_mac": "00:11:22:33:44:55",
  "wifi_signal": 75,
  "qr_code_type": "STORE",
  "qr_code_id": "SBX-WL-1000001",
  "system_info": "iOS 16.0",
  "brand": "Apple",
  "model": "iPhone 14 Pro",
  "page_path": "pages/scan/index",
  "referer": "pages/index/index"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log_id": 1000000001,
    "store_id": 1000001,
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "scan_time": "2023-07-16T16:00:00Z",
    "success_flag": true,
    "wifi_ssid": "Starbucks-Guest"
  }
}
```

### 2. 查询扫码日志列表

获取扫码日志列表，支持多条件筛选。

**请求**:

- 方法: `GET`
- 路径: `/scan-logs`

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| store_id       | integer | 否   | 按门店ID筛选                    |
| user_union_id  | string  | 否   | 按用户UnionID筛选               |
| start_time     | string  | 否   | 开始时间，ISO8601格式            |
| end_time       | string  | 否   | 结束时间，ISO8601格式            |
| success_flag   | boolean | 否   | 按连接结果筛选                   |
| wifi_ssid      | string  | 否   | 按WiFi名称筛选                  |
| qr_code_type   | string  | 否   | 按二维码类型筛选                 |
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
        "log_id": 1000000001,
        "store_id": 1000001,
        "store_name": "星巴克西湖旗舰店",
        "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
        "user_nickname": "张三",
        "scan_time": "2023-07-16T16:00:00Z",
        "network_type": "4G",
        "success_flag": true,
        "wifi_ssid": "Starbucks-Guest",
        "device_info": "iPhone 14 Pro, iOS 16.0",
        "qr_code_type": "STORE"
      },
      // 更多扫码日志...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 135,
      "total_pages": 7
    }
  }
}
```

### 3. 查询指定门店的每日扫码量

获取指定门店的每日扫码统计数据。

**请求**:

- 方法: `GET`
- 路径: `/stores/{store_id}/scan-stats/daily`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| start_date     | string  | 否   | 开始日期，格式YYYY-MM-DD，默认30天前 |
| end_date       | string  | 否   | 结束日期，格式YYYY-MM-DD，默认今天  |
| success_only   | boolean | 否   | 是否仅统计成功连接，默认false      |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "store_name": "星巴克西湖旗舰店",
    "stats": [
      {
        "date": "2023-07-16",
        "total_scans": 45,
        "success_scans": 42,
        "unique_users": 38,
        "new_users": 5
      },
      {
        "date": "2023-07-15",
        "total_scans": 52,
        "success_scans": 48,
        "unique_users": 45,
        "new_users": 7
      },
      // 更多日期...
    ],
    "summary": {
      "total_period_scans": 1463,
      "average_daily_scans": 48.77,
      "success_rate": 92.3,
      "total_unique_users": 834,
      "total_new_users": 121
    }
  }
}
```

### 4. 查询指定用户的扫码历史记录

获取指定用户的扫码历史记录。

**请求**:

- 方法: `GET`
- 路径: `/users/{user_union_id}/scan-logs`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| store_id       | integer | 否   | 按门店ID筛选                    |
| start_time     | string  | 否   | 开始时间，ISO8601格式            |
| end_time       | string  | 否   | 结束时间，ISO8601格式            |
| success_flag   | boolean | 否   | 按连接结果筛选                   |
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
        "log_id": 1000000001,
        "store_id": 1000001,
        "store_name": "星巴克西湖旗舰店",
        "scan_time": "2023-07-16T16:00:00Z",
        "success_flag": true,
        "wifi_ssid": "Starbucks-Guest",
        "device_info": "iPhone 14 Pro, iOS 16.0"
      },
      {
        "log_id": 1000000025,
        "store_id": 1000023,
        "store_name": "星巴克龙井路店",
        "scan_time": "2023-07-15T14:30:00Z",
        "success_flag": true,
        "wifi_ssid": "Starbucks-Guest",
        "device_info": "iPhone 14 Pro, iOS 16.0"
      }
      // 更多扫码记录...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 23,
      "total_pages": 2
    }
  }
}
```

### 5. 查询扫码连接失败日志

获取连接失败的扫码日志记录。

**请求**:

- 方法: `GET`
- 路径: `/scan-logs/failures`

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| store_id       | integer | 否   | 按门店ID筛选                    |
| start_time     | string  | 否   | 开始时间，ISO8601格式            |
| end_time       | string  | 否   | 结束时间，ISO8601格式            |
| fail_reason_code | string | 否  | 按失败原因码筛选                 |
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
        "log_id": 1000000002,
        "store_id": 1000001,
        "store_name": "星巴克西湖旗舰店",
        "user_union_id": "o6_bmanotheruser_2sgVt7hMZ",
        "user_nickname": "李四",
        "scan_time": "2023-07-16T16:05:00Z",
        "success_flag": false,
        "fail_reason_code": "PASSWORD_ERROR",
        "fail_reason_message": "WiFi密码不正确",
        "wifi_ssid": "Starbucks-Guest",
        "device_info": "Xiaomi 12, Android 12"
      },
      // 更多失败记录...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 43,
      "total_pages": 3
    },
    "failure_stats": {
      "password_error": 18,
      "connection_timeout": 12,
      "signal_weak": 5,
      "device_not_supported": 3,
      "other": 5
    }
  }
}
```

### 6. 更新扫码日志连接结果

更新已有扫码日志的连接结果信息。

**请求**:

- 方法: `PATCH`
- 路径: `/scan-logs/{log_id}/result`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| log_id     | integer | 是   | 日志ID                         |

**请求参数**:

| 参数名               | 类型    | 必填  | 描述                           |
|---------------------|---------|------|--------------------------------|
| success_flag        | boolean | 是   | 是否成功连接WiFi               |
| fail_reason_code    | string  | 否   | 连接失败错误码，success_flag为false时需提供 |
| fail_reason_message | string  | 否   | 连接失败详细信息               |
| wifi_signal         | integer | 否   | WiFi信号强度                   |

**请求示例**:

```json
{
  "success_flag": false,
  "fail_reason_code": "CONNECTION_TIMEOUT",
  "fail_reason_message": "连接超时，请重试",
  "wifi_signal": 45
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "log_id": 1000000003,
    "success_flag": false,
    "fail_reason_code": "CONNECTION_TIMEOUT",
    "fail_reason_message": "连接超时，请重试",
    "wifi_signal": 45,
    "updated_at": "2023-07-16T16:15:00Z"
  }
}
```

## 优惠券相关API

### 1. 创建优惠券

创建新的优惠券。

**请求**:

- 方法: `POST`
- 路径: `/coupons`
- 内容类型: `application/json`

**请求参数**:

| 参数名                | 类型     | 必填  | 描述                           |
|----------------------|----------|------|--------------------------------|
| coupon_name          | string   | 是   | 优惠券名称                      |
| coupon_code          | string   | 否   | 优惠券兑换码，可为空            |
| coupon_type          | string   | 是   | 优惠券类型：DISCOUNT/CASH/GIFT/SHIPPING |
| value                | decimal  | 是   | 优惠券面值                      |
| min_purchase_amount  | decimal  | 否   | 最低消费金额，默认0            |
| usage_limit_per_user | integer  | 否   | 每用户可领取数量，默认1，0表示不限 |
| total_quantity       | integer  | 否   | 总发行量，默认0表示不限         |
| start_time           | string   | 是   | 生效时间，ISO8601格式           |
| end_time             | string   | 是   | 过期时间，ISO8601格式           |
| validity_days        | integer  | 否   | 领取后有效天数                  |
| store_id             | integer  | 否   | 适用门店ID，为空表示全平台通用   |
| description          | string   | 否   | 优惠券详细描述                  |
| status               | integer  | 否   | 状态：1启用，0禁用，默认1       |

**请求示例**:

```json
{
  "coupon_name": "新用户立减10元",
  "coupon_type": "CASH",
  "value": 10.00,
  "min_purchase_amount": 30.00,
  "usage_limit_per_user": 1,
  "total_quantity": 1000,
  "start_time": "2023-07-17T00:00:00Z",
  "end_time": "2023-08-17T23:59:59Z",
  "description": "新用户首次扫码连接WiFi即可领取，满30元可使用",
  "status": 1
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "coupon_name": "新用户立减10元",
    "coupon_code": null,
    "coupon_type": "CASH",
    "value": 10.00,
    "min_purchase_amount": 30.00,
    "usage_limit_per_user": 1,
    "total_quantity": 1000,
    "issued_quantity": 0,
    "start_time": "2023-07-17T00:00:00Z",
    "end_time": "2023-08-17T23:59:59Z",
    "validity_days": null,
    "store_id": null,
    "description": "新用户首次扫码连接WiFi即可领取，满30元可使用",
    "status": 1,
    "created_at": "2023-07-16T17:00:00Z",
    "updated_at": "2023-07-16T17:00:00Z"
  }
}
```

### 2. 查询优惠券列表

获取优惠券列表，支持多条件筛选。

**请求**:

- 方法: `GET`
- 路径: `/coupons`

**查询参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| coupon_name    | string  | 否   | 按优惠券名称筛选（模糊匹配）     |
| coupon_type    | string  | 否   | 按优惠券类型筛选                |
| store_id       | integer | 否   | 按适用门店筛选                  |
| status         | integer | 否   | 按状态筛选：1启用，0禁用        |
| start_time_from | string  | 否   | 开始时间筛选（起）              |
| start_time_to   | string  | 否   | 开始时间筛选（止）              |
| end_time_from   | string  | 否   | 结束时间筛选（起）              |
| end_time_to     | string  | 否   | 结束时间筛选（止）              |
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
        "coupon_name": "新用户立减10元",
        "coupon_type": "CASH",
        "value": 10.00,
        "min_purchase_amount": 30.00,
        "total_quantity": 1000,
        "issued_quantity": 0,
        "start_time": "2023-07-17T00:00:00Z",
        "end_time": "2023-08-17T23:59:59Z",
        "store_id": null,
        "status": 1
      },
      // 更多优惠券...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 8,
      "total_pages": 1
    }
  }
}
```

### 3. 查询优惠券详情

获取指定优惠券的详细信息。

**请求**:

- 方法: `GET`
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
    "coupon_name": "新用户立减10元",
    "coupon_code": null,
    "coupon_type": "CASH",
    "value": 10.00,
    "min_purchase_amount": 30.00,
    "usage_limit_per_user": 1,
    "total_quantity": 1000,
    "issued_quantity": 152,
    "start_time": "2023-07-17T00:00:00Z",
    "end_time": "2023-08-17T23:59:59Z",
    "validity_days": null,
    "store_id": null,
    "store_name": null,
    "description": "新用户首次扫码连接WiFi即可领取，满30元可使用",
    "status": 1,
    "created_at": "2023-07-16T17:00:00Z",
    "updated_at": "2023-07-16T17:00:00Z",
    "usage_stats": {
      "received": 152,
      "used": 43,
      "expired": 5,
      "refunded": 2,
      "available": 102
    }
  }
}
```

### 4. 更新优惠券基本信息

更新指定优惠券的基本信息。

**请求**:

- 方法: `PUT`
- 路径: `/coupons/{coupon_id}`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_name | string  | 否   | 优惠券名称                      |
| coupon_code | string  | 否   | 优惠券兑换码                    |
| description | string  | 否   | 优惠券详细描述                  |

**请求示例**:

```json
{
  "coupon_name": "新用户专享立减10元",
  "description": "新用户首次扫码连接WiFi即可领取，满30元可使用，有效期30天"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "coupon_name": "新用户专享立减10元",
    "coupon_code": null,
    "description": "新用户首次扫码连接WiFi即可领取，满30元可使用，有效期30天",
    "updated_at": "2023-07-16T18:00:00Z"
  }
}
```

### 5. 更新优惠券有效期

更新指定优惠券的有效期设置。

**请求**:

- 方法: `PATCH`
- 路径: `/coupons/{coupon_id}/validity`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| coupon_id  | integer | 是   | 优惠券ID                       |

**请求参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| start_time     | string  | 否   | 生效时间，ISO8601格式           |
| end_time       | string  | 否   | 过期时间，ISO8601格式           |
| validity_days  | integer | 否   | 领取后有效天数                  |

**请求示例**:

```json
{
  "end_time": "2023-09-17T23:59:59Z"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "coupon_id": 1001,
    "start_time": "2023-07-17T00:00:00Z",
    "end_time": "2023-09-17T23:59:59Z",
    "validity_days": null,
    "updated_at": "2023-07-16T18:30:00Z"
  }
}
``` 