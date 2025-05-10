# 面试/预约 API (Interv)

**目录**
- [获取面试/预约信息](#获取面试预约信息)
- [新建面试/预约](#新建面试预约)
- [删除面试/预约](#删除面试预约)
- [更新面试/预约信息](#更新面试预约信息)

## 获取面试/预约信息

- **GET** `/api/interv`
- **描述:** 获取面试或预约的相关信息列表 (可能需要权限验证)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Query Parameters:**
    - `status` (string, optional): 按状态筛选 (e.g., "pending", "scheduled", "completed")
    - `date` (string, optional): 按日期筛选 (e.g., "2025-05-10")
    - `page` (integer, optional): 页码
    - `limit` (integer, optional): 每页数量
- **请求示例:** (可带查询参数, e.g., `/api/interv?status=pending&date=2025-05-10`)
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "获取成功",
    "data": {
      "interviews": [
        { "id": 1, "applicant_id": 101, "interviewer_id": 201, "time": "2025-05-15T10:00:00Z", "status": "scheduled" },
        { "id": 2, "applicant_id": 102, "interviewer_id": 202, "time": "2025-05-16T14:00:00Z", "status": "pending_confirmation" }
      ],
      "total": 2
    }
  }
  ```

## 新建面试/预约

- **POST** `/api/interv`
- **描述:** 创建一个新的面试或预约记录
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `applicant_id` (integer, required): 申请人ID
    - `interviewer_id` (integer, optional): 面试官ID (如果已指定)
    - `requested_time` (string, required): 请求的面试时间 (ISO8601格式, e.g., "2025-05-20T09:00:00Z")
    - `notes` (string, optional): 备注信息
- **请求示例:**
  ```json
  {
    "applicant_id": 103,
    "interviewer_id": 201, // 或根据业务逻辑自动分配
    "requested_time": "2025-05-20T09:00:00Z",
    "notes": "申请初级软件工程师职位"
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "创建成功",
    "data": {
      "interview_id": 3,
      "status": "pending_confirmation" // 或直接 scheduled
    }
  }
  ```
- **响应示例 (失败 - 时间冲突):**
  ```json
  {
    "code": 3001, // 示例错误码
    "message": "请求的时间段已有安排"
  }
  ```

## 删除面试/预约

- **DELETE** `/api/interv`
- **描述:** 删除一个面试或预约记录 (通常需要指定ID)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body (或者作为 Path Parameter `/api/interv/{id}`):**
    - `interview_id` (integer, required): 需要删除的面试/预约ID
- **请求示例:** (通常ID在URL中或请求体中)
  ```json
  // 假设通过请求体传递ID
  {
    "interview_id": 3
  }
  // 或者 DELETE /api/interv/3 (如果路由如此设计)
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "删除成功"
  }
  ```
- **响应示例 (失败 - 记录不存在):**
  ```json
  {
    "code": 3002, // 示例错误码
    "message": "记录未找到或无权限删除"
  }
  ```

## 更新面试/预约信息

- **PUT** `/api/interv`
- **描述:** 更新一个已存在的面试或预约信息 (例如: 更改时间, 状态, 参与人等)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `interview_id` (integer, required): 需要更新的面试/预约ID
    - `new_time` (string, optional): 新的面试时间 (ISO8601格式)
    - `status` (string, optional): 新的状态 (e.g., "rescheduled", "confirmed", "cancelled")
    - `notes` (string, optional): 更新的备注信息
    - `interviewer_id` (integer, optional): 更新的面试官ID
- **请求示例:**
  ```json
  {
    "interview_id": 1,
    "new_time": "2025-05-15T11:00:00Z",
    "status": "rescheduled",
    "notes": "用户请求调整时间"
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
    "code": 3003, // 示例错误码
    "message": "无法更新已完成的面试"
  }
  ```
