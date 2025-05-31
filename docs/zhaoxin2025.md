---
title: zhaoxin2025
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# zhaoxin2025

Base URLs:

# Authentication

# Default

## GET RefreshSession

GET /session

> 返回示例

```json
{
  "success": false,
  "message": "鉴权错误: 您未登录\n",
  "code": 6
}
```

```json
{
  "success": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# admin

## POST Register

POST /register

> Body 请求参数

```json
{
  "netid": "3796546547",
  "name": "说英",
  "password": "llzMpNILxfPIK7U",
  "level": "super"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» netid|body|[netid](#schemanetid)| 是 |none|
|» name|body|string| 是 |none|
|» password|body|string| 是 |none|
|» level|body|[level](#schemalevel)| 是 |none|

#### 枚举值

|属性|值|
|---|---|
|» level|super|
|» level|normal|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET GetStu

GET /stu

> Body 请求参数

```json
{
  "netid": "1234567890"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» netid|body|[netid](#schemanetid)| 否 |none|
|» name|body|string| 否 |none|
|» phone|body|[phone](#schemaphone)| 否 |none|
|» school|body|string| 否 |none|
|» first|body|[department](#schemadepartment)| 否 |none|
|» second|body|[department](#schemadepartment)| 否 |none|
|» interviewer|body|string| 否 |none|
|» star|body|[star](#schemastar)| 否 |none|

#### 枚举值

|属性|值|
|---|---|
|» first|tech|
|» first|video|
|» first|art|
|» second|tech|
|» second|video|
|» second|art|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PUT UpdateStu

PUT /stu

> Body 请求参数

```json
{
  "netid": "1331097662",
  "name": "保奕泽",
  "phone": "50165534150",
  "school": "Lorem dolore deserunt consectetur sed",
  "mastered": "proident enim ex cillum dolor",
  "tomaster": "Lorem sit nostrud do",
  "first": "art",
  "second": "art",
  "que_id": 54,
  "que_time": "2025-02-22 04:49:30",
  "interv": false,
  "interviewer": "aute officia ipsum fugiat",
  "evaluation": "dolore enim adipisicing in cupidatat",
  "star": 3,
  "message": 1,
  "pass": 0,
  "id": 85
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» netid|body|[netid](#schemanetid)| 否 |none|
|» name|body|string| 否 |none|
|» phone|body|[phone](#schemaphone)| 否 |none|
|» school|body|string| 否 |none|
|» mastered|body|string| 否 |none|
|» tomaster|body|string| 否 |none|
|» first|body|[department](#schemadepartment)| 否 |none|
|» second|body|[department](#schemadepartment)| 否 |none|
|» que_id|body|integer| 否 |none|
|» que_time|body|string| 否 |none|
|» interv|body|boolean| 否 |none|
|» interviewer|body|string| 否 |none|
|» evaluation|body|string| 否 |none|
|» star|body|[star](#schemastar)| 否 |none|
|» message|body|integer| 否 |none|
|» pass|body|[bool](#schemabool)| 否 |none|
|» id|body|integer| 是 |ID 编号|

#### 枚举值

|属性|值|
|---|---|
|» first|tech|
|» first|video|
|» first|art|
|» second|tech|
|» second|video|
|» second|art|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET Excelize

GET /excel

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET Stat

GET /stat

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST SetTime

POST /settime

> Body 请求参数

```json
{
  "time": "2025-11-21T17:31:30Z"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» time|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "success": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|

## GET SendResultMessage

GET /send

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# stu

## PUT UpdateMessage

PUT /message

> Body 请求参数

```json
{
  "subscribe": 1,
  "intervtime": 1,
  "intervres": 1
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» subscribe|body|[bool](#schemabool)| 是 |none|
|» intervtime|body|[bool](#schemabool)| 是 |none|
|» intervres|body|[bool](#schemabool)| 是 |none|

> 返回示例

```json
{
  "success": true
}
```

```json
{
  "success": false,
  "message": "鉴权错误: 您未登录\n",
  "code": 6
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|

## DELETE CancelInterv

DELETE /interv{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|integer| 是 |none|

> 返回示例

```json
{
  "success": false,
  "message": "记录不存在: 没找到\n",
  "code": 8
}
```

```json
{
  "success": false,
  "message": "鉴权错误: 面试时间在半小时内或已经错过，无法取消预约\n",
  "code": 6
}
```

```json
{
  "success": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|
|» message|string|true|none||none|
|» code|integer|true|none||none|

## POST AppointInterv

POST /interv{id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|path|integer| 是 |none|

> 返回示例

```json
{
  "success": true
}
```

```json
{
  "success": false,
  "message": "记录不存在: 该面试不存在\n",
  "code": 8
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|
|» message|string|true|none||none|
|» code|integer|true|none||none|

## GET GetResult

GET /result

返回信息：已登录学生对应的面试、该部门QQ群二维码图片URL与群号

> 返回示例

> 200 Response

```json
{
  "success": true,
  "data": {
    "data": {
      "id": 21,
      "NetID": "1234567890",
      "Time": "2025-10-27T10:00:00+08:00",
      "Interviewer": "",
      "Department": "tech",
      "QueID": null,
      "Que": null,
      "Star": 0,
      "Evaluation": "",
      "Pass": 0
    },
    "url": "https://shameless-godfather.info/",
    "qqgroup": "59"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET GetInterv

GET /interv

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|date|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET GetIntervDate

GET /date

> 返回示例

> 200 Response

```json
{
  "success": true,
  "data": {}
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|
|» data|object|true|none||none|

# que

## PUT LuckyDog

PUT /lucky

> Body 请求参数

```json
{
  "netid": "1234567890",
  "queid": 2
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» netid|body|[netid](#schemanetid)| 是 |none|
|» queid|body|integer| 是 |none|

> 返回示例

```json
{
  "success": false,
  "message": "记录不存在: 禁止虚空索敌\n",
  "code": 8
}
```

```json
{
  "success": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# interv

## GET Get

GET /

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|integer| 否 |ID 编号|
|department|query|string| 否 |none|
|interviewer|query|string| 否 |none|
|pass|query|integer| 否 |none|
|date|query|string| 否 |none|
|page|query|integer| 是 |none|
|limit|query|integer| 是 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST New

POST /

> Body 请求参数

```json
{
  "timerange": {
    "start": "2025-10-27T11:00:00+08:00",
    "end": "2025-10-27T16:00:00+08:00"
  },
  "interval": 30
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» timerange|body|object| 是 |none|
|»» start|body|string| 是 |none|
|»» end|body|string| 是 |none|
|» interval|body|integer| 是 |none|

> 返回示例

> 200 Response

```json
{
  "success": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PUT Update

PUT /

> Body 请求参数

```json
{
  "id": 1,
  "netid": "4529691599",
  "department": "video",
  "star": 1,
  "pass": "true",
  "evaluation": "ex velit occaecat elit anim"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|integer| 是 |ID 编号|
|» time|body|string| 否 |none|
|» netid|body|[netid](#schemanetid)| 否 |none|
|» department|body|[department](#schemadepartment)| 否 |none|
|» star|body|[star](#schemastar)| 否 |none|
|» pass|body|[bool](#schemabool)| 否 |none|
|» evaluation|body|string| 否 |none|

#### 枚举值

|属性|值|
|---|---|
|» department|tech|
|» department|video|
|» department|art|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## DELETE Delete

DELETE /

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|array[string]| 是 |ID 编号|

> 返回示例

```json
{
  "success": true
}
```

```json
{
  "success": false,
  "message": "记录不存在: 部分面试记录不存在\n",
  "code": 8
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PUT BlockAndRecover

PUT /block

> Body 请求参数

```json
{
  "timerange": {
    "start": "2022-07-27T10:00:00+08:00",
    "end": "2025-08-27T10:00:00+08:00"
  },
  "block": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» timerange|body|object| 是 |none|
|»» start|body|string| 否 |none|
|»» end|body|string| 否 |none|
|» block|body|integer| 是 |none|

> 返回示例

> 200 Response

```json
{
  "success": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PUT QQGroup

PUT /group

> Body 请求参数

```json
{
  "url": "https://shameless-godfather.info/",
  "qqgroup": "13",
  "department": "video"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» url|body|string| 是 |none|
|» qqgroup|body|string| 是 |none|
|» department|body|[department](#schemadepartment)| 是 |none|

#### 枚举值

|属性|值|
|---|---|
|» department|tech|
|» department|video|
|» department|art|

> 返回示例

> 200 Response

```json
{
  "success": true,
  "message": "string",
  "code": 0
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|
|» message|string|false|none||none|
|» code|integer|false|none||none|

## GET GetQue

GET /que

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|department|query|string| 是 |none|

#### 枚举值

|属性|值|
|---|---|
|department|tech|
|department|video|
|department|art|

> 返回示例

> 200 Response

```json
{
  "success": true,
  "data": {
    "id": 2,
    "Question": "non dolore",
    "Department": "tech",
    "Url": "https://stupendous-travel.com/",
    "Times": 0
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» success|boolean|true|none||none|
|» data|object|true|none||none|
|»» id|integer|true|none||none|
|»» Question|string|true|none||none|
|»» Department|string|true|none||none|
|»» Url|string|true|none||none|
|»» Times|integer|true|none||none|

# 数据模型

<h2 id="tocS_netid">netid</h2>

<a id="schemanetid"></a>
<a id="schema_netid"></a>
<a id="tocSnetid"></a>
<a id="tocsnetid"></a>

```json
"stringstri"

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|string|false|none||none|

<h2 id="tocS_department">department</h2>

<a id="schemadepartment"></a>
<a id="schema_department"></a>
<a id="tocSdepartment"></a>
<a id="tocsdepartment"></a>

```json
"tech"

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|string|false|none||none|

#### 枚举值

|属性|值|
|---|---|
|*anonymous*|tech|
|*anonymous*|video|
|*anonymous*|art|

<h2 id="tocS_level">level</h2>

<a id="schemalevel"></a>
<a id="schema_level"></a>
<a id="tocSlevel"></a>
<a id="tocslevel"></a>

```json
"super"

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|string|false|none||none|

#### 枚举值

|属性|值|
|---|---|
|*anonymous*|super|
|*anonymous*|normal|

<h2 id="tocS_phone">phone</h2>

<a id="schemaphone"></a>
<a id="schema_phone"></a>
<a id="tocSphone"></a>
<a id="tocsphone"></a>

```json
"stringstrin"

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|string|false|none||none|

<h2 id="tocS_star">star</h2>

<a id="schemastar"></a>
<a id="schema_star"></a>
<a id="tocSstar"></a>
<a id="tocsstar"></a>

```json
1

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|integer|false|none||none|

<h2 id="tocS_que">que</h2>

<a id="schemaque"></a>
<a id="schema_que"></a>
<a id="tocSque"></a>
<a id="tocsque"></a>

```json
{
  "question": "string",
  "department": "tech",
  "url": "string",
  "times": 0
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|question|string|false|none||none|
|department|[department](#schemadepartment)|false|none||none|
|url|string|false|none||none|
|times|integer|false|none||none|

<h2 id="tocS_page">page</h2>

<a id="schemapage"></a>
<a id="schema_page"></a>
<a id="tocSpage"></a>
<a id="tocspage"></a>

```json
1

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|integer|false|none||none|

<h2 id="tocS_limit">limit</h2>

<a id="schemalimit"></a>
<a id="schema_limit"></a>
<a id="tocSlimit"></a>
<a id="tocslimit"></a>

```json
1

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|integer|false|none||none|

<h2 id="tocS_bool">bool</h2>

<a id="schemabool"></a>
<a id="schema_bool"></a>
<a id="tocSbool"></a>
<a id="tocsbool"></a>

```json
1

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|integer|false|none||none|

