package db

import (
	v1 "auth/v1"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var sqlConn *sql.DB

func Init() {
	var err error
	sqlConn, err = sql.Open("mysql", "root:969929899@localhost:3306/test")
	if err != nil {
		panic(err)
	}
}

func CreatePolicy(p *v1.Policy) {
	_, err := sqlConn.Exec("insert into policy(`version`,`statements`) values (?,?)", p.Version, p.Statements)
	if err != nil {
		panic(err)
	}
}
func CreateRole(r *v1.Role) {
	_, err := sqlConn.Exec("insert into role(`name`) values (?)", r.Name)
	if err != nil {
		panic(err)
	}
}
func QueryUserPolicy(uid int64) []*v1.Policy {
	query, err := sqlConn.Query("select p.id,p.version,p.statements from policy p left join policy_ref pr on p.id = pr.policy_id where pr.ref_id = ? and pr.ref_type = 'user'", uid)
	if err != nil {
		panic(err)
	}
	pl := make([]*v1.Policy, 0)
	for query.Next() {
		p := v1.Policy{}
		err := query.Scan(&p.Id, &p.Version, &p.Statements)
		if err != nil {
			panic(err)
		}
		pl = append(pl, &p)
	}
	return pl
}

func QueryRolePolicy(rid int64) []v1.Policy {

}

func QueryGroupPolicy(gid int64) []v1.Policy {

}
