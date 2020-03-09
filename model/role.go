package model

func CreateRokle(r *Role) int64 {
	result, e := mysql.Exec("insert into `role`(`name`) values (?)", r.Name)
	checkErr(e)
	i, e := result.LastInsertId()
	checkErr(e)
	return i
}

func QueryRole(id int64) *Role {
	rows, e := mysql.Query("select id,name from `role` where id = ?", id)
	checkErr(e)
	r := Role{}
	if rows.Next() {
		rows.Scan(&r.Id, &r.Name)
	}
	return &r
}

func QueryPolicyIdsByRole(role int) *[]int {
	rows, err := mysql.Query("select p_id from r_p_ref where r_id = ?", role)
	checkErr(err)
	idList := make([]int, 10)
	var id int
	for rows.Next() {
		rows.Scan(&id)
		idList = append(idList, id)
	}
	return &idList
}
