# iam

模仿阿里云ram策略的实现

根据传入的context 和 policy 判断是否能执行对应的操作

context 代表某个人想要执行某些操作

比如`tom` 想要对 `tom:food:bread` 执行 `food:eat` 的操作(tom想吃tom的面包)

context 是这样
```json
{
   "Action"   : "food:eat",
   "Resource" : "tom:food:bread",
   "Requester": {
      "name": "tom"
   },
   "Condition": ""
}
```

Resource字符串支持根据context动态计算, `$` 代表context `$.resuester` 就是请求者的属性
```json
{
   "name": "tom"
}
```

policy 代表`什么人` `能不能` `对哪些资源` `进行哪些操作`

比如 定义一个策略，每个人只能吃自己的食物

policy

```json
{
  "Version": 1,
  "Statements": [
    {
    "Action": ["food:eat"],
    "Resource": ["{$.requester.name}:food:*"],
    "Effect": "Allow"
    }
  ]
}
```



## 权限判断的基本规则

### deny 优先

 任何权限如果存在deny直接deny
 
### 信任参数

 只负责检查权限不负责判断用户或者资源信息的正确性
 

## milestones

### V1

1. 实现 `.` 取值,`a[b]`数组,`a[b:c]`切片,`a[*].b`扫描语法
  
