package v1

//policy  = {
//     <version_block>,
//     <statement_block>
//}
//<version_block> = "Version" : ("1")
//<statement_block> = "Statements" : [ <statement>, <statement>, ... ]
//<statement> = {
//    <effect_block>,
//    <action_block>,
//    <resource_block>,
//    <condition_block?>
//}
//<effect_block> = "Effect" : ("Evaluate" | "Deny")
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
	Version    Version      `json:"Version"`
	Statements []*Statement `json:"Statements"`
}

// evaluate by context
// allow: represent this policy allow actions in context
// match: represent action and resource in policy statements match these in context
func (p *Policy) Evaluate(c *Context) (allow bool, match bool, err error) {
	for _, statement := range p.Statements {
		statementMatch, err := statement.match(c)
		if err != nil {
			break
		}
		if statementMatch {
			match = true
			if statement.Effect == Deny {
				allow = false
				break
			}
			if statement.Effect == Allow {
				allow = true
			}
		}
	}
	return allow, match, err
}
