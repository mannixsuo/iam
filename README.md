# iam

```
context = {
   <UserInfo>
   <resourceInfo>
}

<userInfo> = {
  <key>:<value>
  <key>:<value>
  ...
}
<resourceInfo> = {
  <key>:<value>
  <key>:<value>
  ...
}
<key> = (string)
<value> = (string|number)

```
## 权限判断的基本规则

### deny first：
 任何权限如果存在deny直接deny
 
### 规则2 
 只负责检查权限不负责判断用户或者资源信息的正确性


```
who can (or can’t) do what to which resources.

Zelkova translates policies into precise mathematical language and then uses automated
 reasoning tools to check properties of the policies. 

```

## milestones
### V1
1. 实现 RBAC
   1. 实现 `.`,`a[b]`,`a[b:c]`,`a[*].b`,`a[*][c]` 语法
  
  
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:*",
                "cloudwatch:*",
                "ec2:*"
            ],
            "Resource": "*"
        }
    ]
}
```