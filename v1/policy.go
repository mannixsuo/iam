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
	Id         int          `json:"-"`
	Version    Version     `json:"Version"`
	Statements []*Statement `json:"Statements"`
}

//
func (p *Policy) Allow(c *Context) (allow bool, match bool, err error) {
	var m bool
	for _, s := range p.Statements {
		m, err = s.match(c)
		if err != nil {
			break
		}
		if m {
			match = true
			if s.Effect == Deny {
				allow = false
				break
			}
			allow = true
		}
	}
	return allow, match, err
}
