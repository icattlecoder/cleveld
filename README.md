cleveld
=======
[toc]

## 简介

cleveld 是golang实现的一个key-value数据库服务，基于leveldb 存储引擎。

## 安装

1. 安装leveldb
2. 安装dleveld
```
git clone git@github.com:icattlecoder/cleveld.git
```

## 协议

### get 查询

请求：

```
g <dbname>\r\n
<key> \r\n
\r\n
```

返回：

```
ok\r\n
<value>\r\n
\r\n
```


### set 增加

请求：

```
s <dbname>\r\n 
<key>\r\n
<value>\r\n
\r\n
```

返回：

```
ok\r\n
\r\n
```

### delete 删除


请求：

```
d <dbname>\r\n
<key>\r\n
\r\n
```

### list 列举

```
l <dbname>\r\n
<limit>\r\n
<key>\r\n
\r\n
``` 

## 客户端API
> 目前仅实现了go语言的客户端
### golang

