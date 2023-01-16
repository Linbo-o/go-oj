# Go-Oj

## 一、项目总览

### 1.1 目录结构

```bash
├─app                           // 程序具体逻辑代码
│  ├─http                           // http 请求处理逻辑
│  │  ├─contorllers                    //api控制器
│  │  │  └─v1
│  │  │      ├─auth
│  │  │      └─problem
│  │  └─middlewares                    //中间件
│  ├─models                         //存放各模块的对象以及对象的各种方法和相关函数
│  │  ├─category
│  │  ├─problem
│  │  ├─problem-category
│  │  ├─submit
│  │  ├─testcase
│  │  └─user
│  └─requests                      //存放各种请求的对应结构，包含请求数据校验的逻辑
│      └─validators
├─bootstrap                    //存放各个大的模块的启动函数，包括配置、日志、数据库、路由的初始化
├─config                       //存放各种配置信息
├─pkg                          //内置辅助包，对各种第三方包进行封装，为逻辑的具体实现提供一些应用级的工具函数
│  ├─app
│  ├─auth
│  ├─captcha
│  ├─config
│  ├─database
│  ├─helpers
│  ├─jwt
│  ├─logger
│  ├─redis
│  ├─response
│  ├─sms
│  └─verifycode
├─routes                      //路由
├─sql                         //创建数据库使用的sql
├─storage                     //存放日志文件和用户提交的代码
│    ├─code
│    └─logs
├─ .env                       //配置文件

```

### 1.2 代码行数

![](E:\GoProject\src\go-oj\images\1.png)

### 1.3 所有API

| 请求方法 | API 地址                      | 说明                               |
| -------- | ----------------------------- | ---------------------------------- |
| POST     | /v1/auth/signup/phone/exist   | 手机号是否已注册                   |
| POST     | /v1/auth/signup/email/exist   | 邮箱是否已经注册                   |
| POST     | /v1/auth/verify/captcha       | 发送图案验证码                     |
| POST     | /v1/auth/verify/captcha       | 发送短信验证码                     |
| POST     | /v1/auth/signup/using-phone   | 使用手机号注册                     |
| POST     | /v1/auth/login/using-phone    | 使用手机验证码登录                 |
| POST     | /v1/auth/login/using-password | 使用 用户名/手机号/邮箱 + 密码登录 |
| POST     | /v1/problem/create            | 创建题目                           |
| POST     | /v1/problem/modify            | 修改题目                           |
| POST     | /v1/problem/problem-list      | 获取题目列表                       |
| POST     | /v1/problem/problem-detail    | 获取题目详尽信息                   |
| POST     | /v1/problem/problem-judge     | 提交题目，进行评测                 |
| POST     | /v1/problem/problem-status    | 查看评测结果                       |



| 请求方法 | API地址  | 说明                                                         |
| -------- | -------- | ------------------------------------------------------------ |
| GET      | /v1/test | 测试用，通过jwt获取用户信息，方便看调用其他api时的用户信息变化 |

### 1.4 第三方依赖

使用到的开源库：

* `gin`—— 路由、路由组、中间件

* `zap`—— 高性能日志方案

* `viper`—— 配置信息

* `cast`——类型转换

* `redis`——数据缓存

* `base64Captcha`——图案验证码

* `jwt`——JWT操作

* `govalidator`——验证信息

* `aliyun-communicate`——阿里云短信服务

  

* base64Captcha——图片验证码

## 二、API设计

### 1. 1 用户注册

手机号注册用户会调用四个 API ：

1. 调用 `signup/phone/exist` 检查手机号是否已被注册；
2. 调用 `verify-codes/captcha` 获取图片验证码，验证后才有发`数字验证码`权限；
3. 调用 `verify-codes/phone` 先检验图片验证码是否正确，若正确发送短信验证码；
4. 调用 `signup/using-phone` 检验验证码，然后再次校验手机号是否已经注册，注册用户,注册成功后会返回一个token。

注：2、3、4的调用顺序不可变，也不能缺少，否则都会注册失败

### 1.2 用户登录

手机验证码登录会调用3个API：

1. 调用 `verify-codes/captcha` 获取图片验证码，验证后才有发『数字验证码』权限；
2. 调用 `verify-codes/phone` 先检验图片验证码是否正确，然后再次校验手机号是否已经注册，若正确发送短信验证码；
3. 调用 `login/using-phone` 登录用户,登录成功后会返回一个token。

### 1.3 问题相关模块

问题相关模块API的调用都会经过jwt认证中间件，认证是否为用户，区分用户权限( 普通用户和管理员 )

#### 1.3.1 查看题目

1、首先通过jwt认证中间件

2、调用`problem-list`查看题目列表，获取题目的唯一标识identity

3、调用`problem-detail`查看题目的详细信息

#### 1.3.2  创建/修改题目

1、首先通过jwt认证中间件

2、调用`problem/create`或者是`problem/modify`接口，**判断为管理员后**可根据提交的信息对题目进行响应的操作

#### 1.3.3 提交代码评测

1、首先通过jwt认证中间件

2、调用`problem-judge`接口，提交代码，返回题目的identity(唯一标识)

3、调用`problem-status`，从数据库中获取评测状态和结果

[详情](http://192.168.1.8:10393/shareDoc?issue=00ade2eed528a8a4bf6abd3b9ef579a7&target_id=0f9839ea-beb2-41e0-a791-846f18d3dbd4)

