# 管理员 API

**目录**
- [登录](#登录)
- [登出](#登出)
- [获取登录状态](#获取登录状态)
- [更新管理员信息](#更新管理员信息)
- [获取学生列表](#获取学生列表)
- [更新学生信息](#更新学生信息)
- [注册新用户 (管理员或特定权限)](#注册新用户-管理员或特定权限)

## 登录

- **POST** `/api/admin`
- **描述:** 管理员登录
- **参数列表:**
  - **Request Body:**
    - `username` (string, required): 管理员用户名
    - `password` (string, required): 管理员密码
- **请求示例:**
  ```json
  {
    "username": "admin",
    "password": "password"
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "登录成功",
    "data": {
      "token": "some_jwt_token",
      "user_info": {
        "id": 1,
        "username": "admin",
        "role": 2
      }
    }
  }
  ```
- **响应示例 (失败):**
  ```json
  {
    "code": 1001, // 示例错误码
    "message": "用户名或密码错误"
  }
  ```

## 登出

- **DELETE** `/api/admin`
- **描述:** 管理员登出
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
- **请求示例:** (无特定请求体, 依赖Header中的Token)
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "登出成功"
  }
  ```

## 获取登录状态

- **GET** `/api/admin/`
- **描述:** 获取当前管理员登录状态
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
- **请求示例:** (无特定请求体, 依赖Header中的Token)
- **响应示例 (已登录):**
  ```json
  {
    "code": 0,
    "message": "已登录",
    "data": {
      "user_info": {
        "id": 1,
        "username": "admin",
        "role": 2
      }
    }
  }
  ```
- **响应示例 (未登录):**
  ```json
  {
    "code": 1002, // 示例错误码
    "message": "未登录或会话已过期"
  }
  ```

## 更新管理员信息

- **PUT** `/api/admin/`
- **描述:** 更新管理员信息 (例如: 密码)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `old_password` (string, required): 旧密码
    - `new_password` (string, required): 新密码
- **请求示例:**
  ```json
  {
    "old_password": "old_password",
    "new_password": "new_secure_password"
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "信息更新成功"
  }
  ```
- **响应示例 (失败):**
  ```json
  {
    "code": 1003, // 示例错误码
    "message": "旧密码错误"
  }
  ```

## 获取学生列表

- **GET** `/api/admin/stu`
- **描述:** 管理员获取学生用户列表
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Query Parameters:**
    - `page` (integer, optional): 页码，默认为1
    - `limit` (integer, optional): 每页数量，默认为10
- **请求示例:** (可带分页参数, e.g., `/api/admin/stu?page=1&limit=10`)
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "获取成功",
    "data": {
      "students": [
        { "id": 101, "name": "学生A", "student_id": "S001" },
        { "id": 102, "name": "学生B", "student_id": "S002" }
      ],
      "total": 2
    }
  }
  ```

## 更新学生信息

- **PUT** `/api/admin/stu`
- **描述:** 管理员更新特定学生信息
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `user_id` (integer, required): 需要更新的学生用户ID
    - `name` (string, optional): 学生姓名
    - `email` (string, optional): 学生邮箱
    - `...` (any, optional): 其他可更新的学生字段
- **请求示例:**
  ```json
  {
    "user_id": 101,
    "name": "学生A更新",
    "email": "studentA_updated@example.com"
    // 其他可更新字段
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "学生信息更新成功"
  }
  ```

## 注册新用户 (管理员或特定权限)

- **POST** `/api/admin/register`
- **描述:** 管理员或特定权限用户注册新用户 (可能是学生或其他角色)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `username` (string, required): 新用户名
    - `password` (string, required): 新用户密码
    - `role` (integer, required): 用户角色 (例如: 1 代表学生)
    - `email` (string, optional): 用户邮箱
- **请求示例:**
  ```json
  {
    "username": "newuser",
    "password": "password123",
    "role": 1, // 假设1代表学生角色
    "email": "newuser@example.com"
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "用户注册成功",
    "data": {
      "user_id": 103
    }
  }
  ```
- **响应示例 (失败 - 用户名已存在):**
  ```json
  {
    "code": 1004, // 示例错误码
    "message": "用户名已存在"
  }
  ```
