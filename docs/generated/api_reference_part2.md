# WiFi城市API接口文档 - 第二部分

## WIFI配置相关API

### 1. 新增门店WIFI配置

为指定门店添加新的WiFi配置。

**请求**:

- 方法: `POST`
- 路径: `/stores/{store_id}/wifi-configs`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**请求参数**:

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
  "ssid": "Starbucks-Event",
  "password": "event2023",
  "encryption_type": "WPA2",
  "wifi_type": "EVENT",
  "max_connections": 30
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "wifi_id": 3,
    "store_id": 1000001,
    "ssid": "Starbucks-Event",
    "encryption_type": "WPA2",
    "wifi_type": "EVENT",
    "max_connections": 30,
    "last_updated": "2023-07-16T09:00:00Z"
  }
}
```

### 2. 批量新增门店WIFI配置

为指定门店批量添加多个WiFi配置。

**请求**:

- 方法: `POST`
- 路径: `/stores/{store_id}/wifi-configs/batch`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |

**请求参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| wifi_configs   | array   | 是   | WiFi配置数组                    |

每个WiFi配置包含的参数与"新增门店WIFI配置"接口相同。

**请求示例**:

```json
{
  "wifi_configs": [
    {
      "ssid": "Starbucks-Member",
      "password": "member2023",
      "encryption_type": "WPA3",
      "wifi_type": "CUSTOMER",
      "max_connections": 50
    },
    {
      "ssid": "Starbucks-VIP",
      "password": "vip2023",
      "encryption_type": "WPA3",
      "wifi_type": "CUSTOMER",
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
    "wifi_configs": [
      {
        "wifi_id": 4,
        "store_id": 1000001,
        "ssid": "Starbucks-Member",
        "encryption_type": "WPA3",
        "wifi_type": "CUSTOMER",
        "max_connections": 50,
        "last_updated": "2023-07-16T09:15:00Z"
      },
      {
        "wifi_id": 5,
        "store_id": 1000001,
        "ssid": "Starbucks-VIP",
        "encryption_type": "WPA3",
        "wifi_type": "CUSTOMER",
        "max_connections": 20,
        "last_updated": "2023-07-16T09:15:00Z"
      }
    ],
    "total_added": 2
  }
}
```

### 3. 查询门店所有WIFI配置列表

获取指定门店的所有WiFi配置。

**请求**:

- 方法: `GET`
- 路径: `/stores/{store_id}/wifi-configs`

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
    "store_name": "星巴克西湖旗舰店",
    "wifi_configs": [
      {
        "wifi_id": 1,
        "ssid": "Starbucks-Guest",
        "encryption_type": "WPA2",
        "wifi_type": "CUSTOMER",
        "max_connections": 100,
        "last_updated": "2023-07-15T10:30:00Z"
      },
      {
        "wifi_id": 2,
        "ssid": "Starbucks-Staff",
        "encryption_type": "WPA3",
        "wifi_type": "STAFF",
        "max_connections": 20,
        "last_updated": "2023-07-15T10:30:00Z"
      },
      {
        "wifi_id": 3,
        "ssid": "Starbucks-Event",
        "encryption_type": "WPA2",
        "wifi_type": "EVENT",
        "max_connections": 30,
        "last_updated": "2023-07-16T09:00:00Z"
      },
      {
        "wifi_id": 4,
        "ssid": "Starbucks-Member",
        "encryption_type": "WPA3",
        "wifi_type": "CUSTOMER",
        "max_connections": 50,
        "last_updated": "2023-07-16T09:15:00Z"
      },
      {
        "wifi_id": 5,
        "ssid": "Starbucks-VIP",
        "encryption_type": "WPA3",
        "wifi_type": "CUSTOMER",
        "max_connections": 20,
        "last_updated": "2023-07-16T09:15:00Z"
      }
    ],
    "total_count": 5
  }
}
```

### 4. 查询单个WIFI配置详情

获取指定WiFi配置的详情。

**请求**:

- 方法: `GET`
- 路径: `/wifi-configs/{wifi_id}`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| wifi_id    | integer | 是   | WiFi配置ID                     |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "wifi_id": 1,
    "store_id": 1000001,
    "store_name": "星巴克西湖旗舰店",
    "ssid": "Starbucks-Guest",
    "encryption_type": "WPA2",
    "wifi_type": "CUSTOMER",
    "max_connections": 100,
    "last_updated": "2023-07-15T10:30:00Z"
  }
}
```

### 5. 更新WIFI配置

更新指定的WiFi配置信息。

**请求**:

- 方法: `PUT`
- 路径: `/wifi-configs/{wifi_id}`
- 内容类型: `application/json`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| wifi_id    | integer | 是   | WiFi配置ID                     |

**请求参数**:

| 参数名              | 类型    | 必填  | 描述                           |
|--------------------|---------|------|--------------------------------|
| ssid               | string  | 否   | WiFi名称                        |
| password           | string  | 否   | WiFi密码（明文，系统会加密存储）  |
| encryption_type    | string  | 否   | 加密类型：WPA2/WPA3/WEP/OPEN     |
| wifi_type          | string  | 否   | WiFi类型：CUSTOMER/STAFF/EVENT  |
| max_connections    | integer | 否   | 最大连接数                      |

**请求示例**:

```json
{
  "password": "newpassword2023",
  "encryption_type": "WPA3",
  "max_connections": 120
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "wifi_id": 1,
    "store_id": 1000001,
    "ssid": "Starbucks-Guest",
    "encryption_type": "WPA3",
    "wifi_type": "CUSTOMER",
    "max_connections": 120,
    "last_updated": "2023-07-16T10:00:00Z"
  }
}
```

### 6. 删除WIFI配置

删除指定的WiFi配置。

**请求**:

- 方法: `DELETE`
- 路径: `/wifi-configs/{wifi_id}`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| wifi_id    | integer | 是   | WiFi配置ID                     |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "wifi_id": 1,
    "deleted_at": "2023-07-16T11:00:00Z"
  }
}
```

### 7. 查询门店特定类型的WIFI配置

获取指定门店特定类型的WiFi配置。

**请求**:

- 方法: `GET`
- 路径: `/stores/{store_id}/wifi-configs/type/{wifi_type}`

**路径参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| store_id   | integer | 是   | 门店ID                         |
| wifi_type  | string  | 是   | WiFi类型：CUSTOMER/STAFF/EVENT |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "store_id": 1000001,
    "store_name": "星巴克西湖旗舰店",
    "wifi_type": "CUSTOMER",
    "wifi_configs": [
      {
        "wifi_id": 4,
        "ssid": "Starbucks-Member",
        "encryption_type": "WPA3",
        "max_connections": 50,
        "last_updated": "2023-07-16T09:15:00Z"
      },
      {
        "wifi_id": 5,
        "ssid": "Starbucks-VIP",
        "encryption_type": "WPA3",
        "max_connections": 20,
        "last_updated": "2023-07-16T09:15:00Z"
      }
    ],
    "total_count": 2
  }
}
```

## 用户信息相关API

### 1. 获取用户详情

根据用户UnionID获取用户详细信息。

**请求**:

- 方法: `GET`
- 路径: `/users/{user_union_id}`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "open_id": "o6_bmjrPTlm6_2sgVt7hMZOPfL",
    "wechat_nickname": "张三",
    "wechat_avatar_url": "https://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
    "phone_number": "13812345678",
    "phone_country_code": "86",
    "gender": 1,
    "language": "zh_CN",
    "country": "中国",
    "province": "浙江",
    "city": "杭州",
    "first_seen": "2023-01-15T08:30:00Z",
    "last_seen": "2023-07-16T12:00:00Z"
  }
}
```

### 2. 根据OpenID获取用户详情

根据微信OpenID获取用户详细信息。

**请求**:

- 方法: `GET`
- 路径: `/users/openid/{open_id}`

**路径参数**:

| 参数名    | 类型    | 必填  | 描述                           |
|----------|---------|------|--------------------------------|
| open_id  | string  | 是   | 用户微信OpenID                 |

**响应**:

与"获取用户详情"接口响应格式相同。

### 3. 根据手机号获取用户详情

根据手机号获取用户详细信息。

**请求**:

- 方法: `GET`
- 路径: `/users/phone/{phone_number}`

**路径参数**:

| 参数名        | 类型    | 必填  | 描述                           |
|--------------|---------|------|--------------------------------|
| phone_number | string  | 是   | 用户手机号                      |

**查询参数**:

| 参数名               | 类型    | 必填  | 描述                           |
|---------------------|---------|------|--------------------------------|
| phone_country_code  | string  | 否   | 手机号国家区号，默认"86"        |

**响应**:

与"获取用户详情"接口响应格式相同。

### 4. 创建/更新用户档案

创建新用户或更新现有用户信息。

**请求**:

- 方法: `POST`
- 路径: `/users`
- 内容类型: `application/json`

**请求参数**:

| 参数名             | 类型    | 必填  | 描述                           |
|-------------------|---------|------|--------------------------------|
| user_union_id     | string  | 是   | 用户微信UnionID                |
| open_id           | string  | 是   | 用户微信OpenID                 |
| wechat_nickname   | string  | 否   | 微信昵称                        |
| wechat_avatar_url | string  | 否   | 微信头像URL                     |
| phone_number      | string  | 否   | 手机号                          |
| phone_country_code| string  | 否   | 手机号国家区号，默认"86"         |
| gender            | integer | 否   | 性别：1男，2女，0未知            |
| language          | string  | 否   | 用户语言，如"zh_CN"              |
| country           | string  | 否   | 用户国家                         |
| province          | string  | 否   | 用户省份                         |
| city              | string  | 否   | 用户城市                         |

**请求示例**:

```json
{
  "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
  "open_id": "o6_bmjrPTlm6_2sgVt7hMZOPfL",
  "wechat_nickname": "张三",
  "wechat_avatar_url": "https://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
  "phone_number": "13812345678",
  "phone_country_code": "86",
  "gender": 1,
  "language": "zh_CN",
  "country": "中国",
  "province": "浙江",
  "city": "杭州"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "open_id": "o6_bmjrPTlm6_2sgVt7hMZOPfL",
    "wechat_nickname": "张三",
    "wechat_avatar_url": "https://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/0",
    "phone_number": "13812345678",
    "phone_country_code": "86",
    "gender": 1,
    "language": "zh_CN",
    "country": "中国",
    "province": "浙江",
    "city": "杭州",
    "first_seen": "2023-07-16T13:00:00Z",
    "last_seen": "2023-07-16T13:00:00Z",
    "is_new_user": true
  }
}
```

### 5. 更新用户手机号

更新指定用户的手机号信息。

**请求**:

- 方法: `PATCH`
- 路径: `/users/{user_union_id}/phone`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**请求参数**:

| 参数名               | 类型    | 必填  | 描述                           |
|---------------------|---------|------|--------------------------------|
| phone_number        | string  | 是   | 手机号                          |
| phone_country_code  | string  | 否   | 手机号国家区号，默认"86"         |

**请求示例**:

```json
{
  "phone_number": "13987654321",
  "phone_country_code": "86"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "phone_number": "13987654321",
    "phone_country_code": "86",
    "updated_at": "2023-07-16T14:00:00Z"
  }
}
```

### 6. 更新用户地理位置信息

更新指定用户的地理位置信息。

**请求**:

- 方法: `PATCH`
- 路径: `/users/{user_union_id}/location`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**请求参数**:

| 参数名    | 类型    | 必填  | 描述                           |
|----------|---------|------|--------------------------------|
| country  | string  | 否   | 用户国家                         |
| province | string  | 否   | 用户省份                         |
| city     | string  | 否   | 用户城市                         |

**请求示例**:

```json
{
  "country": "中国",
  "province": "上海",
  "city": "上海"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "country": "中国",
    "province": "上海",
    "city": "上海",
    "updated_at": "2023-07-16T14:30:00Z"
  }
}
```

### 7. 用户绑定手机号

将手机号与用户UnionID绑定。

**请求**:

- 方法: `POST`
- 路径: `/users/{user_union_id}/bind-phone`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**请求参数**:

| 参数名               | 类型    | 必填  | 描述                           |
|---------------------|---------|------|--------------------------------|
| phone_number        | string  | 是   | 手机号                          |
| phone_country_code  | string  | 否   | 手机号国家区号，默认"86"         |
| verify_code         | string  | 是   | 验证码                          |

**请求示例**:

```json
{
  "phone_number": "13712345678",
  "phone_country_code": "86",
  "verify_code": "123456"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "phone_number": "13712345678",
    "phone_country_code": "86",
    "bind_time": "2023-07-16T15:00:00Z"
  }
}
```

### 8. 用户解绑手机号

解除用户UnionID与手机号的绑定。

**请求**:

- 方法: `POST`
- 路径: `/users/{user_union_id}/unbind-phone`
- 内容类型: `application/json`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**请求参数**:

| 参数名         | 类型    | 必填  | 描述                           |
|---------------|---------|------|--------------------------------|
| verify_code   | string  | 是   | 验证码                          |

**请求示例**:

```json
{
  "verify_code": "123456"
}
```

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "user_union_id": "o6_bmasdasdsad6_2sgVt7hMZOPfL",
    "unbind_time": "2023-07-16T15:30:00Z"
  }
}
```

### 9. 查询用户扫码门店历史

获取用户历史扫码过的门店列表。

**请求**:

- 方法: `GET`
- 路径: `/users/{user_union_id}/scan-stores`

**路径参数**:

| 参数名          | 类型    | 必填  | 描述                           |
|----------------|---------|------|--------------------------------|
| user_union_id  | string  | 是   | 用户微信UnionID                |

**查询参数**:

| 参数名      | 类型    | 必填  | 描述                           |
|------------|---------|------|--------------------------------|
| page       | integer | 否   | 页码，默认1                     |
| page_size  | integer | 否   | 每页记录数，默认20              |
| success    | boolean | 否   | 是否仅查询成功连接的记录         |

**响应**:

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "store_id": 1000001,
        "store_name": "星巴克西湖旗舰店",
        "address": "杨公堤1-3号西湖景区入口",
        "wifi_ssid": "Starbucks-Guest",
        "last_scan_time": "2023-07-16T12:00:00Z",
        "success_flag": true,
        "scan_count": 5
      },
      {
        "store_id": 1000023,
        "store_name": "星巴克龙井路店",
        "address": "龙井路18号",
        "wifi_ssid": "Starbucks-Guest",
        "last_scan_time": "2023-07-15T14:30:00Z",
        "success_flag": true,
        "scan_count": 2
      }
      // 更多门店...
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
``` 