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

[小程序api]()

[后台api]()

## 通用api

### 获取登录状态

GET {baseurl}/api/

> 返回示例

未登录
```json
{
  "success": true,
  "data": "未登录"
}
```

已登录
```json
{
  "success": true,
  "data": {
    "netid": "0474135450"
  }
}
```

### 刷新session登录状态

GET {baseurl}/api/session

> 返回示例

失败
```json
{
  "success": false,
  "message": "鉴权错误: 您未登录\n",
  "code": 6
}
```

成功
```json
{
  "success": true
}
```
