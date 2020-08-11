package policy

import "auth/auth"

//policy  = {
//     <version_block>,
//     <statement_block>
//}
//<version_block> = "Version" : ("1")
//<statement_block> = "Statement" : [ <statement>, <statement>, ... ]
//<statement> = {
//    <effect_block>,
//    <action_block>,
//    <resource_block>,
//    <condition_block?>
//}
//<effect_block> = "Effect" : ("Allow" | "Deny")
//<action_block> = "Action" :
//    ("*" | [<action_string>, <action_string>, ...])
//<resource_block> = "Resource" :
//    ("*" | [<resource_string>, <resource_string>, ...])
//<condition_block> = "Condition" : <condition_map>
//<condition_map> = {
//  <condition_type_string> : {
//      <condition_key_string> : <condition_value_list>,
//      <condition_key_string> : <condition_value_list>,
//      ...
//  },
//  <condition_type_string> : {
//      <condition_key_string> : <condition_value_list>,
//      <condition_key_string> : <condition_value_list>,
//      ...
//  }, ...
//}
//<condition_value_list> = [<condition_value>, <condition_value>, ...]
//<condition_value> = ("String" | "Number" | "Boolean")

type Policy struct {
	Id        int         `json:"-"`
	Version   Version     `json:"Version"`
	Statement []Statement `json:"Statements"`
}

// MatchContext 判断该policy是否和Context 所对应
// 对应指的是
//   policy的action  包含context中的 action
//   AND
//   policy的resource包含context中的 resource
func (p *Policy) MatchContext(c *auth.Context) {
	if p.matchContextAction(c) && p.matchContextResource(c) {

	}
}

// matchContextAction 判断policy的action 是否包含context中的 action
func (p *Policy) matchContextAction(c *auth.Context) bool {
	return false
}

// matchContextResource 判断policy的resource 是否包含context中的 resource
func (p *Policy) matchContextResource(c *auth.Context) bool {
	return false
}

func (p *Policy) Evaluate() {}
