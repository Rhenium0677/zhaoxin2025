# 学生 API

**目录**
- [学生登录](#学生登录)
- [获取学生登录状态](#获取学生登录状态)
- [学生登出](#学生登出)
- [更新学生信息](#更新学生信息)
- [更新/提交额外信息](#更新提交额外信息)

## 学生登录

- **POST** `/api/stu`
- **描述:** 学生用户登录
- **参数列表:**
  - **Request Body:**
    - `student_id` (string, required): 学生学号
    - `password` (string, required): 登录密码
- **请求示例:**
  ```json
  {
    "student_id": "S001",
    "password": "student_password"
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
        "id": 101,
        "student_id": "S001",
        "name": "学生A",
        "role": 1 // 假设1代表学生角色
      }
    }
  }
  ```
- **响应示例 (失败):**
  ```json
  {
    "code": 2001, // 示例错误码
    "message": "学号或密码错误"
  }
  ```

## 获取学生登录状态

- **GET** `/api/stu`
- **描述:** 获取当前学生登录状态
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
        "id": 101,
        "student_id": "S001",
        "name": "学生A",
        "role": 1
      }
    }
  }
  ```
- **响应示例 (未登录):**
  ```json
  {
    "code": 2002, // 示例错误码
    "message": "未登录或会话已过期"
  }
  ```

## 学生登出

- **DELETE** `/api/stu`
- **描述:** 学生用户登出
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

## 更新学生信息

- **PUT** `/api/stu`
- **描述:** 学生更新自己的信息 (例如: 密码, 联系方式等)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `email` (string, optional): 新邮箱地址
    - `phone` (string, optional): 新电话号码
    - `old_password` (string, optional): 旧密码 (如果需要修改密码)
    - `new_password` (string, optional): 新密码 (如果需要修改密码)
    - `...` (any, optional): 其他可更新的学生字段
- **请求示例:**
  ```json
  {
    "email": "studentA_new_email@example.com",
    "phone": "13800138000"
    // 其他可更新字段, 如密码等
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
    "code": 2003, // 示例错误码
    "message": "更新失败，无效的输入"
  }
  ```

## 更新/提交额外信息

- **PUT** `/api/stu/sub`
- **描述:** 学生更新或提交一些特定的额外信息 (具体含义根据业务定义，例如提交作业、报名信息等)
- **参数列表:**
  - **Headers:**
    - `Authorization` (string, required): Bearer Token, 用于身份验证。
  - **Request Body:**
    - `subject` (string, optional): 相关主题或科目
    - `assignment_id` (string, optional): 作业或任务ID
    - `content` (string, optional): 提交的具体内容
    - `...` (any, optional): 其他根据业务定义的特定信息字段
- **请求示例:**
  ```json
  {
    "subject": "计算机科学导论",
    "assignment_id": "A001",
    "content": "这是我的作业内容..."
    // 或其他特定信息字段
  }
  ```
- **响应示例 (成功):**
  ```json
  {
    "code": 0,
    "message": "信息提交成功"
  }
  ```
- **响应示例 (失败):**
  ```json
  {
    "code": 2004, // 示例错误码
    "message": "提交失败，请检查内容"
  }
  ```
