package v1

// {project}/{module}/{group}/{user}/{resource_id}:{resource_attr}/{resource_attr}...
// 允许用户查询组内的资源信息
// Effect: Evaluate
// Action: listPortPolicy
// Resources: envc-platform/portPolicy/${user_group}/*/*

// 不允许用户查询组内id为1的资源信息
// Effect: Deny
// Action: listPortPolicy
// Resources: envc-platform/portPolicy/${user.group}/*/1

// context:{
// key:value,
// key:value
//}
//
