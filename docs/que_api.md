# 问题/队列 API (Que)

**目录**
- [获取问题/队列信息](#获取问题队列信息)
- [新建问题/队列条目](#新建问题队列条目)
- [删除问题/队列条目](#删除问题队列条目)
- [更新问题/队列条目](#更新问题队列条目)

## 获取问题/队列信息

- **GET** `/api/que`
- **描述:** 获取问题或队列列表 (可能根据用户角色或权限进行过滤)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证 (如果需要认证)。
  - **Query Parameters:**
    - `type` (string, optional): 类型筛选 (e.g., "faq", "ticket")
    - `category` (string, optional): 分类筛选 (e.g., "general", "technical")
    - `status` (string, optional): 状态筛选 (e.g., "open", "closed", "pending")
    - `page` (integer, optional): 页码
    - `limit` (integer, optional): 每页数量
- **请求示例:** (可带查询参数, e.g., `/api/que?type=faq&category=general`)
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "获取成功",
    "data": {
      "questions": [
        { "id": 1, "title": "如何重置密码?", "category": "account", "answer_preview": "您可以通过..." },
        { "id": 2, "title": "服务如何收费?", "category": "billing", "answer_preview": "我们的收费标准是..." }
      ],
      "total": 2
    }
  }
  ```

## 新建问题/队列条目

- **POST** `/api/que`
- **描述:** 创建一个新的问题或队列条目
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证 (如果需要认证)。
  - **Request Body:**
    - `title` (string, required): 问题或条目的标题
    - `category` (string, optional): 分类
    - `details` (string, required): 详细描述
    - `user_id` (integer, optional): 提交问题的用户ID (如果系统未自动关联)
- **请求示例:**
  ```json
  {
    "title": "关于API集成的疑问",
    "category": "technical_support",
    "details": "我无法成功调用 /api/data 接口...",
    "user_id": 105 // 提交问题的用户ID
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "创建成功",
    "data": {
      "question_id": 3,
      "status": "open"
    }
  }
  ```
- **响应示例 (失败 - 输入无效):**
  ```json
  {
    "code": 4001, // 示例错误码
    "message": "标题和内容不能为空"
  }
  ```

## 删除问题/队列条目

- **DELETE** `/api/que`
- **描述:** 删除一个问题或队列条目 (通常需要指定ID)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body (或者作为 Path Parameter `/api/que/{id}`):**
    - `question_id` (integer, required): 需要删除的问题/条目ID
- **请求示例:** (通常ID在URL中或请求体中)
  ```json
  // 假设通过请求体传递ID
  {
    "question_id": 3
  }
  // 或者 DELETE /api/que/3 (如果路由如此设计)
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "删除成功"
  }
  ```
- **响应示例 (失败 - 条目不存在或无权限):**
  ```json
  {
    "code": 4002, // 示例错误码
    "message": "条目未找到或无权限删除"
  }
  ```

## 更新问题/队列条目

- **PUT** `/api/que`
- **描述:** 更新一个已存在的问题或队列条目 (例如: 修改内容, 状态, 分配给某人等)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `question_id` (integer, required): 需要更新的问题/条目ID
    - `title` (string, optional): 新标题
    - `category` (string, optional): 新分类
    - `details` (string, optional): 更新后的详细描述
    - `answer` (string, optional): 对问题的回答 (如果适用)
    - `status` (string, optional): 新状态 (e.g., "answered", "closed", "in_progress")
    - `assigned_to_admin_id` (integer, optional): 分配给处理该问题的管理员ID
- **请求示例:**
  ```json
  {
    "question_id": 1,
    "answer": "您可以通过访问 xxx.com/reset-password 页面来重置您的密码。",
    "status": "answered",
    "assigned_to_admin_id": 5 // 分配给管理员处理
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "更新成功"
  }
  ```
- **响应示例 (失败 - 无效操作):**
  ```json
  {
    "code": 4003, // 示例错误码
    "message": "无法更新已关闭的问题"
  }
  ```
