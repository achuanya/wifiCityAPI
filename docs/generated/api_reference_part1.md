# WiFi城市API接口文档 - 第一部分

## API概述

### 基本信息

- **基础URL**: `https://api.wificity.example.com/v1`
- **内容类型**: `application/json`
- **字符编码**: `UTF-8`

### 认证方式

系统使用基于Token的认证机制：

1. 通过微信授权接口获取临时凭证
2. 使用临时凭证获取访问Token
3. 在后续请求的Header中添加Token:
   ```
   Authorization: Bearer {access_token}
   ```

### 响应格式

所有API响应均使用JSON格式，基本结构如下：

```json
{
  "code": 200,           // 状态码，200表示成功，非200表示失败
  "message": "success",  // 状态消息
  "data": {              // 响应数据，失败时可能为null或包含错误详情
    // 具体的返回数据
  }
}
```

### 常见状态码

| 状态码 | 说明 |
|--------|------|
| 200    | 请求成功 |
| 400    | 请求参数错误 |
| 401    | 未授权或授权失败 |
| 403    | 权限不足 |
| 404    | 资源不存在 |
| 500    | 服务器内部错误 |

### 分页机制

对于返回多条记录的接口，支持分页查询，分页参数：

- `page`: 页码，从1开始，默认为1
- `page_size`: 每页记录数，默认为20，最大为100

分页响应包含以下额外信息：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      // 实际数据条目
    ],
    "pagination": {
      "page": 1,         // 当前页码
      "page_size": 20,   // 每页记录数
      "total": 125,      // 总记录数
      "total_pages": 7   // 总页数
    }
  }
}
```

## 门店相关API

### 1. 新增门店

创建一个新的门店记录。

**请求**:

- 方法: `POST`
- 路径: `/stores`
- 内容类型: `application/json`

**请求参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| name       | string  | 是   | 门店名称                        |
| country    | string  | 否   | 国家，默认为"中国"               |
| province   | string  | 是   | 省份                           |
| city       | string  | 是   | 城市                           |
| district   | string  | 是   | 区/县                          |
| address    | string  | 是   | 详细地址                        |
| latitude   | decimal | 是   | 门店纬度，精确到小数点后6位       |
| longitude  | decimal | 是   | 门店经度，精确到小数点后6位       |
| phone      | string  | 是   | 联系电话                        |
| status     | integer | 否   | 门店状态，1正常，0停用，默认为1   |

**请求示例**:

```json
{
  "name": "星巴克西湖店",
  "country": "中国",
  "province": "浙江省",
  "city": "杭州市",
  "district": "西湖区",
  "address": "杨公堤1号",
  "latitude": 30.245843,
  "longitude": 120.126843,
  "phone": "0571-88888888",
  "status": 1
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "name": "星巴克西湖店",
    "country": "中国",
    "province": "浙江省",
    "city": "杭州市",
    "district": "西湖区",
    "address": "杨公堤1号",
    "latitude": 30.245843,
    "longitude": 120.126843,
    "phone": "0571-88888888",
    "wifi_count": 0,
    "status": 1,
    "created_at": "2023-07-15T10:30:00Z",
    "updated_at": "2023-07-15T10:30:00Z"
  }
}
```

### 2. 同时新增门店及初始WIFI配置

创建门店并同时配置初始WiFi信息。

**请求**:

- 方法: `POST`
- 路径: `/stores/with-wifi`
- 内容类型: `application/json`

**请求参数**:

门店信息参数与"新增门店"接口相同，增加以下参数：

| 参数名               | 类型    | 必填  | 描述                           |
|---------------------|---------|------|--------------------------------|
| wifi_configs        | array   | 是   | WiFi配置数组                    |

每个WiFi配置包含：

| 参数名              | 类型    | 必填  | 描述                           |
|--------------------|---------|------|--------------------------------|
| ssid               | string  | 是   | WiFi名称                        |
| password           | string  | 是   | WiFi密码（明文，系统会加密存储）  |
| encryption_type    | string  | 是   | 加密类型：WPA2/WPA3/WEP/OPEN     |
| wifi_type          | string  | 是   | WiFi类型：CUSTOMER/STAFF/EVENT  |
| max_connections    | integer | 否   | 最大连接数，默认50              |

**请求示例**:

```json
{
  "name": "星巴克西湖店",
  "country": "中国",
  "province": "浙江省",
  "city": "杭州市",
  "district": "西湖区",
  "address": "杨公堤1号",
  "latitude": 30.245843,
  "longitude": 120.126843,
  "phone": "0571-88888888",
  "status": 1,
  "wifi_configs": [
    {
      "ssid": "Starbucks-Guest",
      "password": "welcome123",
      "encryption_type": "WPA2",
      "wifi_type": "CUSTOMER",
      "max_connections": 100
    },
    {
      "ssid": "Starbucks-Staff",
      "password": "staff2023",
      "encryption_type": "WPA3",
      "wifi_type": "STAFF",
      "max_connections": 20
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
    "store": {
      "store_id": 1000001,
      "name": "星巴克西湖店",
      "country": "中国",
      "province": "浙江省",
      "city": "杭州市",
      "district": "西湖区",
      "address": "杨公堤1号",
      "latitude": 30.245843,
      "longitude": 120.126843,
      "phone": "0571-88888888",
      "wifi_count": 2,
      "status": 1,
      "created_at": "2023-07-15T10:30:00Z",
      "updated_at": "2023-07-15T10:30:00Z"
    },
    "wifi_configs": [
      {
        "wifi_id": 1,
        "store_id": 1000001,
        "ssid": "Starbucks-Guest",
        "encryption_type": "WPA2",
        "wifi_type": "CUSTOMER",
        "max_connections": 100,
        "last_updated": "2023-07-15T10:30:00Z"
      },
      {
        "wifi_id": 2,
        "store_id": 1000001,
        "ssid": "Starbucks-Staff",
        "encryption_type": "WPA3",
        "wifi_type": "STAFF",
        "max_connections": 20,
        "last_updated": "2023-07-15T10:30:00Z"
      }
    ]
  }
}
```

### 3. 查询门店列表

获取门店列表，支持分页和多条件筛选。

**请求**:

- 方法: `GET`
- 路径: `/stores`

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| page       | integer | 否   | 页码，默认1                      |
| page_size  | integer | 否   | 每页记录数，默认20               |
| name       | string  | 否   | 门店名称（模糊匹配）             |
| province   | string  | 否   | 省份                           |
| city       | string  | 否   | 城市                           |
| district   | string  | 否   | 区/县                          |
| status     | integer | 否   | 门店状态：1正常，0停用           |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "store_id": 1000001,
        "name": "星巴克西湖店",
        "country": "中国",
        "province": "浙江省",
        "city": "杭州市",
        "district": "西湖区",
        "address": "杨公堤1号",
        "latitude": 30.245843,
        "longitude": 120.126843,
        "phone": "0571-88888888",
        "wifi_count": 2,
        "status": 1,
        "created_at": "2023-07-15T10:30:00Z",
        "updated_at": "2023-07-15T10:30:00Z"
      },
      // 更多门店记录...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 125,
      "total_pages": 7
    }
  }
}
```

### 4. 查询门店详情

获取指定门店的详细信息。

**请求**:

- 方法: `GET`
- 路径: `/stores/{store_id}`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "name": "星巴克西湖店",
    "country": "中国",
    "province": "浙江省",
    "city": "杭州市",
    "district": "西湖区",
    "address": "杨公堤1号",
    "latitude": 30.245843,
    "longitude": 120.126843,
    "phone": "0571-88888888",
    "wifi_count": 2,
    "status": 1,
    "created_at": "2023-07-15T10:30:00Z",
    "updated_at": "2023-07-15T10:30:00Z"
  }
}
```

### 5. 更新门店基本信息

更新指定门店的基本信息。

**请求**:

- 方法: `PUT`
- 路径: `/stores/{store_id}`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**请求参数**:

与"新增门店"接口相同，但所有参数均为可选，仅更新提供的字段。

**请求示例**:

```json
{
  "name": "星巴克西湖旗舰店",
  "address": "杨公堤1-3号"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "name": "星巴克西湖旗舰店",
    "country": "中国",
    "province": "浙江省",
    "city": "杭州市",
    "district": "西湖区",
    "address": "杨公堤1-3号",
    "latitude": 30.245843,
    "longitude": 120.126843,
    "phone": "0571-88888888",
    "wifi_count": 2,
    "status": 1,
    "created_at": "2023-07-15T10:30:00Z",
    "updated_at": "2023-07-15T14:15:00Z"
  }
}
```

### 6. 更新门店联系电话

单独更新门店的联系电话。

**请求**:

- 方法: `PATCH`
- 路径: `/stores/{store_id}/phone`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**请求参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| phone      | string  | 是   | 新的联系电话                    |

**请求示例**:

```json
{
  "phone": "0571-87654321"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "phone": "0571-87654321",
    "updated_at": "2023-07-15T14:30:00Z"
  }
}
```

### 7. 更新门店地理位置信息

单独更新门店的地理位置信息。

**请求**:

- 方法: `PATCH`
- 路径: `/stores/{store_id}/location`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**请求参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| country    | string  | 否   | 国家                           |
| province   | string  | 否   | 省份                           |
| city       | string  | 否   | 城市                           |
| district   | string  | 否   | 区/县                          |
| address    | string  | 否   | 详细地址                        |
| latitude   | decimal | 否   | 门店纬度                        |
| longitude  | decimal | 否   | 门店经度                        |

**请求示例**:

```json
{
  "latitude": 30.246001,
  "longitude": 120.127002,
  "address": "杨公堤1-3号西湖景区入口"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "country": "中国",
    "province": "浙江省",
    "city": "杭州市",
    "district": "西湖区",
    "address": "杨公堤1-3号西湖景区入口",
    "latitude": 30.246001,
    "longitude": 120.127002,
    "updated_at": "2023-07-15T14:45:00Z"
  }
}
```

### 8. 更新门店状态

单独更新门店的状态（启用/停用）。

**请求**:

- 方法: `PATCH`
- 路径: `/stores/{store_id}/status`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**请求参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| status     | integer | 是   | 门店状态：1正常，0停用           |

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
    "store_id": 1000001,
    "status": 0,
    "updated_at": "2023-07-15T15:00:00Z"
  }
}
```

### 9. 删除门店

删除指定的门店记录。

**请求**:

- 方法: `DELETE`
- 路径: `/stores/{store_id}`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "deleted_at": "2023-07-15T16:00:00Z"
  }
}
```

### 10. 查询指定区域内的门店

获取指定行政区域内的所有门店。

**请求**:

- 方法: `GET`
- 路径: `/stores/region`

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| province   | string  | 否   | 省份                           |
| city       | string  | 否   | 城市                           |
| district   | string  | 否   | 区/县                          |
| page       | integer | 否   | 页码，默认1                     |
| page_size  | integer | 否   | 每页记录数，默认20              |
| status     | integer | 否   | 门店状态筛选：1正常，0停用       |

**注意**: 必须至少提供一个区域参数(province/city/district)

**响应**:

与"查询门店列表"接口响应格式相同。

### 11. 查询附近门店

获取指定坐标附近的门店。

**请求**:

- 方法: `GET`
- 路径: `/stores/nearby`

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| latitude   | decimal | 是   | 当前位置纬度                    |
| longitude  | decimal | 是   | 当前位置经度                    |
| radius     | decimal | 否   | 搜索半径，单位公里，默认5公里     |
| limit      | integer | 否   | 返回记录数量，默认10，最大50     |
| status     | integer | 否   | 门店状态筛选：1正常，0停用       |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "store_id": 1000001,
        "name": "星巴克西湖旗舰店",
        "address": "杨公堤1-3号西湖景区入口",
        "latitude": 30.246001,
        "longitude": 120.127002,
        "distance": 0.3,
        "wifi_count": 2,
        "status": 1
      },
      {
        "store_id": 1000023,
        "name": "星巴克龙井路店",
        "address": "龙井路18号",
        "latitude": 30.248523,
        "longitude": 120.133568,
        "distance": 0.8,
        "wifi_count": 1,
        "status": 1
      }
      // 更多附近门店...
    ]
  }
}
``` 