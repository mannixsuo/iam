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

Zelkova translates policies into precise mathematical language and then uses automated reasoning tools to check properties of the policies. 

```

## milestones
### V1
1. 实现 RBAC
  1. 角色 增删改
  2. 用户信息操作
  3. 用户权限判断
    1. 实现 `.`,`a[b]`,`a[b:c]` 语法
  4. 策略(policy)操作
  5. 角色和策略绑定解绑