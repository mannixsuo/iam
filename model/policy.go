package model

import (
	"fmt"
	"strconv"
	"strings"
)

func AttachPolicyToUser(userId int64, policyId int64) (int64, error) {
	result, e := mysql.Exec("insert into u_p_ref(u_id,p_id) values (?,?)", userId, policyId)
	checkErr(e)
	return result.LastInsertId()
}

func AttachPolicyToRole(roleId int64, policyId int64) (int64, error) {
	result, e := mysql.Exec("insert into r_p_ref(r_id,p_id) values (?,?)", roleId, policyId)
	checkErr(e)
	return result.LastInsertId()
}

func AttachPolicyToGroup(groupId int64, policyId int64) (int64, error) {
	result, e := mysql.Exec("insert into g_p_ref(g_id,p_id) values (?,?)", groupId, policyId)
	checkErr(e)
	return result.LastInsertId()
}

// QueryPolicyByIdsAndActions query policy list by policy id list and action list
func QueryPolicyByIdsAndActions(idList []int, actionParams []string) *[]*Policy {
	ids := BuildNumberList(idList)
	actionParameter := BuildStringList(actionParams)
	queryString := "SELECT p.id as policy_id, p.version, s.id as statement_id, s.effect, a.id as action_id,a.action,r.id as resource_id, r.resource FROM policy p 	LEFT JOIN statement s ON p.id = s.policy LEFT JOIN action a ON a.statement = s.id LEFT JOIN resource r ON r.statement = s.id WHERE p.id in %s and (a.action in \"%s\" or a.action LIKE '%\\*%')"
	queryString = fmt.Sprintf(queryString, ids, actionParameter)
	fmt.Printf(queryString)
	rows, e := mysql.Query(queryString)
	checkErr(e)
	var result []*Policy
	policys := make(map[int]*Policy)
	statements := make(map[int]*Statement)
	actions := make(map[int]*string)
	resources := make(map[int]*string)
	var policy *Policy
	var statement *Statement
	var present bool
	var policyId, statementId, actionId, resourceId int
	var version, effect, action, resource string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&policyId, &version, &statementId, &effect, &actionId, &action, &resourceId, &resource)
		if policy, present = policys[policyId]; !present {
			policy = &Policy{
				Id:        policyId,
				Version:   version,
				Statement: make([]*Statement, 0, 100),
			}
			policys[policyId] = policy
			result = append(result, policy)
		}
		if statement, present = statements[statementId]; !present {
			statement = &Statement{
				Id:       statementId,
				Effect:   effect,
				Action:   make([]string, 0, 100),
				Resource: make([]string, 0, 100),
			}
			policy.Statement = append(policy.Statement, statement)
			statements[statementId] = statement
		}
		if _, present := actions[actionId]; !present {
			statement.Action = append(statement.Action, action)
			actions[actionId] = &action
		}
		if _, present := resources[resourceId]; !present {
			statement.Resource = append(statement.Resource, resource)
			resources[resourceId] = &resource
		}
	}
	return &result
}

// QueryPolicy simple query policy info but this is not efficient
func QueryPolicy(policyId int) *Policy {
	rows, e := mysql.Query("select `version` from `policy` where id = ?", policyId)
	checkErr(e)
	p := Policy{}
	for rows.Next() {
		e := rows.Scan(&p.Version)
		checkErr(e)
		statement, e := mysql.Query("select `id` , `effect` from statement where policy = ? ", policyId)
		checkErr(e)
		for statement.Next() {
			s := Statement{}
			e := statement.Scan(&s.Id, &s.Effect)
			checkErr(e)
			s.Action = GetActions(s.Id)
			s.Resource = GetResources(s.Id)
			p.Statement = append(p.Statement, &s)
		}
	}
	return &p
}

// GetActions query action by statementId
func GetActions(statementId int) []string {
	actionRows, e := mysql.Query("select action from action where statement = ?", statementId)
	checkErr(e)
	var actions []string
	for actionRows.Next() {
		var action string
		e := actionRows.Scan(&action)
		checkErr(e)
		actions = append(actions, action)
	}
	return actions
}

// GetResources query resource by statementId
func GetResources(statementId int) []string {
	resourceRows, e := mysql.Query("select resource from resource where statement = ?", statementId)
	checkErr(e)
	var resources []string
	for resourceRows.Next() {
		var resource string
		e := resourceRows.Scan(&resource)
		checkErr(e)
		resources = append(resources, resource)
	}
	return resources
}

// SavePolicy it will save all the policy info into database
func SavePolicy(p *Policy) {
	tx, _ := mysql.Begin()
	result, e := tx.Exec("insert into `policy`(`version`) values (?)", p.Version)
	checkErr(e)
	policyId, e := result.LastInsertId()
	checkErr(e)
	for _, s := range p.Statement {
		result, e = tx.Exec("insert into `statement`(`effect`,`policy`) values (?,?)", s.Effect, policyId)
		checkErr(e)
		statementId, e := result.LastInsertId()
		checkErr(e)
		for _, a := range s.Action {
			fmt.Println(a)
			_, e := tx.Exec("insert into `action`(action,statement) values (?,?)", a, statementId)
			checkErr(e)
		}
		for _, r := range s.Resource {
			fmt.Println(r)
			_, e = tx.Exec("insert into `resource`(resource,statement) values (?,?)", r, statementId)
			checkErr(e)
		}
	}
	e = tx.Commit()
	checkErr(e)
}

type PolicyCheck struct {
	Policys  []*Policy
	Context  map[string]interface{}
	Resource QueryResource
}

type UserInfo struct {
	name         string
	role         []string
	resourceName string
}

type QueryResource struct {
	Resource string
}

// CheckAllow check access by policy and resource and statement's effect
func (p *PolicyCheck) CheckAllow() bool {
	var decision bool
	for _, policy := range p.Policys {
		for _, statement := range policy.Statement {
			for _, resource := range statement.Resource {
				fullResourceName := p.buildResource(resource)
				fmt.Println(fullResourceName)
				if resourceEquals(fullResourceName, p.Resource.Resource) {
					if statement.Effect == "allow" {
						decision = true
					} else {
						return false
					}
				}
			}
		}
	}
	fmt.Println(decision)
	return decision
}

// resourceEquals compare whether two resource are equals
// a:b:c:d:relation-id
// first compare the prefix
// split a:b:c:d by `:`(a,b,c,d)
// then compare one by one
// * equals anything
// if policy is only * it means any resource
// if two string are equals continue compare next one else the two resource are not equal
// if two prefix are equals then we compare relation id
func resourceEquals(fullResourceName string, policyResource string) bool {
	if policyResource == "*" {
		return true
	}
	frList := strings.Split(fullResourceName, ":")
	prList := strings.Split(policyResource, ":")
	prLen := len(prList)
	if len(frList) != prLen {
		return false
	}
	for i := 0; i < len(prList)-1; i++ {
		if prList[i] == "*" || prList[i] == frList[i] {
			continue
		} else {
			fmt.Printf("resource %s and %s are not equals \n", prList[i], frList[i])
			return false
		}
	}
	if !compareRelationId(frList[prLen-1], prList[prLen-1]) {
		return false
	}
	fmt.Printf("resource %s and %s are equals \n", fullResourceName, policyResource)
	return true
}

func ActionEquals(queryAction string, policyAction string) bool {
	if policyAction == "*" {
		return true
	}
	frList := strings.Split(queryAction, ":")
	prList := strings.Split(policyAction, ":")
	prLen := len(prList)
	if len(frList) != prLen {
		return false
	}
	for i := 0; i < len(prList)-1; i++ {
		if prList[i] == "*" || prList[i] == frList[i] {
			continue
		} else {
			fmt.Printf("action %s and %s are not equals \n", prList[i], frList[i])
			return false
		}
	}
	if !compareRelationId(frList[prLen-1], prList[prLen-1]) {
		return false
	}
	fmt.Printf("action %s and %s are equals \n", queryAction, policyAction)
	return true
}

// compareRelationId relationId look like car/tyres which means tyres of a car
// car/* means anything of the car
// we first split the relation by `/`
// then compare then one by one like resource compare

// a/b/c/*
// a/*
// deny

// a
// a/b/c
// allow

func compareRelationId(fullRelation string, policyRelation string) bool {
	if policyRelation == "*" {
		return true
	}

	frList := strings.Split(fullRelation, "/")
	prList := strings.Split(policyRelation, "/")
	if len(prList) > len(frList) {
		return false
	}
	for i := 0; i < len(prList); i++ {
		if prList[i] == "*" || prList[i] == frList[i] {
			continue
		} else {
			fmt.Printf("relation %s and %s are not equals \n", prList[i], frList[i])
			return false
		}
	}
	fmt.Printf("relation %s and %s are equals \n", fullRelation, policyRelation)
	return true
}

func (p *PolicyCheck) buildResource(resource string) string {
	splitedResource := strings.Split(resource, ":")

	for index, syntax := range splitedResource {
		if strings.HasPrefix(syntax, "$") {
			value := p.getValue(strings.TrimPrefix(syntax, "$"))
			splitedResource[index] = interface2String(value)
		}
	}
	return strings.Join(splitedResource, ":")
}

func (p *PolicyCheck) getValue(key string) interface{} {
	if value, ok := p.Context[key]; ok {
		return value
	}
	return nil
}

func interface2String(inter interface{}) string {

	switch inter.(type) {
	case string:
		return inter.(string)
	case int:
		return strconv.Itoa(inter.(int))
	}
	return ""
}
