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

policy 代表 `能不能` `对哪些资源` `进行哪些操作`

### policy基本元素

|元素名称|描述|
|:----|:----|
|效力（Effect）|	授权效力包括两种：允许（Allow）和拒绝（Deny）。|
|操作（Action）|	操作是指对具体资源的操作。|
|资源（Resource）|	资源是指被授权的具体对象。|
|限制条件（Condition）|	限制条件是指授权生效的限制条件。|

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

代码示例：

```go 
// 定义一个policy 能吃自己的食物
   p := Policy{
		Version: 1,
		Statements: []*Statement{
			{
				Action:   &Action{"food:eat"},
				Resource: &Resource{"{$.requester.name}:food:*"},
				Effect:   Allow,
			}
		},
	}
	ctx := &Context{
		Action:    "food:eat",
		Requester: map[string]interface{}{"name": "tom"},
		Resource:  "tom:food:bread",
	}
    ctx := &Context{
		Action:    "food:eat",
		Requester: map[string]interface{}{"name": "tom"},
		Resource:  "jerry:food:bread",
	}
	allow, match, err := p.Evaluate(ctx)
    fmt.Println(allow, match, err)
    allow, match, err = p.Evaluate(ctx)
    fmt.Println(allow, match, err)

output:
--------------
true true <nil>
false false <nil>

```

这里只是进行了权限判断,权限具体赋予在个人还是角色还是用户组上需要自己定义

## 权限判断的基本规则

### deny 优先

 任何权限如果存在deny直接deny
 
### 信任参数

 只负责检查权限不负责判断用户或者资源信息的正确性
 

## milestones

### V1

1. 实现 `.` 取值,`a[b]`数组,`a[b:c]`切片,`a[*].b`扫描语法
  
