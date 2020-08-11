package main

import (
	"auth/policy"
	"auth/policy/v1"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

func main() {
	policy.Init()
	p := v1.Policy{}
	// policys := "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":[\"s3:ListBucket\"],\"Effect\":\"Allow\",\"Resource\":[\"arn:aws:s3:::personal-files\"],\"Condition\":{\"StringLike\":{\"s3:prefix\":[\"tyrchen/*\"]}}},{\"Action\":[\"s3:GetObject\",\"s3:PutObject\"],\"Effect\":\"Allow\",\"Resource\":[\"arn:aws:s3:::personal-files/tyrchen/*\"]}]}"
	policys := "{\"Version\":\"1\",\"Statement\":[{\"Effect\":\"allow\",\"Action\":[\"tcd:ListS2Vm\"],\"Resource\":[\"tcd:server:*:$tenantId:vm/s2/*\"]},{\"Effect\":\"allow\",\"Action\":[\"tcd:ListS3Vm\"],\"Resource\":[\"tcd:server:*:$tenantId:vm/s2/*\"]},{\"Effect\":\"allow\",\"Action\":[\"tcd:createS2Vm\"],\"Resource\":[\"tcd:server:*:$tenantId:vm/s2\"]},{\"Effect\":\"allow\",\"Action\":[\"tcd:createS3Vm\"],\"Resource\":[\"tcd:server:*:$tenantId:vm/s3\"]}]}"
	policys = "{\"Version\":\"1\",\"Statements\":[{\"Effect\":\"allow\",\"Action\":[\"tcd:describeS2Vm\",\"tcd:ListS2Vm\",\"tcd:startS2Vm\",\"tcd:stopS2Vm\",\"tcd:rebootS2Vm\"],\"Resource\":[\"tcd:server:*:$tenantId:vm/s2/*\"]}]}"
	policys = "{\"Version\":\"1\",\"Statements\":[{\"Effect\":\"allow\",\"Action\":[\"tcd:*\"],\"Resource\":[\"tcd:*:*:*:*\"]}]}"
	policys = "{\"Version\":\"1\",\"Statements\":[{\"Effect\":\"allow\",\"Action\":[\"tcd:*\"],\"Resource\":[\"tcd:server:*:$tenantId:vm/s2/id/123\"]}]}"
	policys = "{\"Version\":\"1\",\"Statements\":[{\"Effect\":\"deny\",\"Action\":[\"tcd:*\"],\"Resource\":[\"tcd:server:*:$tenantId:vm/s2/id/123\"]}]}"
	err := json.NewDecoder(strings.NewReader(policys)).Decode(&p)
	if err != nil {
		panic(err)
	}
	// for i := 0; i < 100; i++ {
	// 	model.SavePolicy(&p)
	// }
	now := time.Now()
	// model.SavePolicy(&p)
	// policy := mysqlTableTest()
	ids := []int{37453, 37454, 6}
	policy := policy.QueryPolicyByIdsAndActions(ids, []string{"tcd:ListS3Vm"})
	check := policy.PolicyCheck{
		Policys: *policy,
		Context: map[string]interface{}{
			"tenantId": 11,
			"user":     "mannix",
		},
		Resource: policy.QueryResource{
			Resource: "tcd:server:*:10:vm/s3",
		},
	}
	fmt.Println(check.CheckAllow())
	fmt.Println(time.Since(now).String())
	fmt.Println(policy)
}

func mysqlTableTest() *v1.Policy {
	return policy.QueryPolicy(4)
}
