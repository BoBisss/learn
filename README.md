# 个人博客系统介绍

## 项目运行步骤

**当前项目go 版本为1.24.5**

数据库使用mysql ，端口为9002

**数据库连接和端口配置均在personal_blog/config/config.yml，可直接修改**

1. 创建数据库名称为gorm，如果已经被占用则新建一个数据库表名，并将personal_blog/config/config.yml的数据库连接表名修改掉，同时将连接数据库的用户名密码修改
2. cd 到personal_blog目录
3. 执行go run main.go（执行不成功使用go mod tidy 或者有无启用go模块集成）

## 模块

### 注册

post请求：http://localhost:9001/auth/register

```go
{
    "username":"zhangsan",
    "password":"123456",
    "email":"a@qq.com"
}
```

### 登陆

post请求：http://localhost:9001/auth/login

```plain
{
    "username":"zhangsan",
    "password":"123456"
}
```

### 添加文章

post请求：http://localhost:9001/api/posts

```plain
{
    "title":"文章1",
    "content":"文章1内容"
}
```

### 获取文章列表

get请求：http://localhost:9001/api/posts

### 获取文章详情

get请求：http://localhost:9001/api/posts/getById?id=1

### 修改文章

patch请求：http://localhost:9001/api/posts

```plain
{
    "id":2,
    "title":"文章1",
    "content":"文章1内容"
}
```

### 删除文章

delete请求：http://localhost:9001/api/posts

```plain
{
    "id":1
}
```

### 评论文章

post请求：http://localhost:9001/api/comments

```plain
{
    "postId":2,
    "content":"文章1评论"
}
```

### 查询文章所有评论

get请求：http://localhost:9001/api/comments

```plain
{
    "postId":2
}
```
