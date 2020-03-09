package model

func CreateGroup(g *Group) (int64, error) {
	result, err := mysql.Exec("insert into `group`(`name`) values (?)", g.Name)
	checkErr(err)
	return result.LastInsertId()
}

func QueryGroup(groupId int64) *Group {
	rows, e := mysql.Query("select id, name from `group` where id = ?", groupId)
	checkErr(e)
	g := Group{}
	if rows.Next() {
		rows.Scan(&g.Name)
	}
	return &g
}

func UserJoinGroup(userId int64, groupId int64) (int64, error) {
	result, e := mysql.Exec("insert into u_g_ref(`user`,`group`) values (?,?)", userId, groupId)
	checkErr(e)
	return result.LastInsertId()
}

func QueryUserGroup(userId int64) *[]Group {
	result, e := mysql.Query("select g.id,g.name from group g left join u_g_ref ugr where ugr.user=?", userId)
	checkErr(e)
	groups := make([]Group, 10)
	if result.Next() {
		g := Group{}
		result.Scan(&g.Id, &g.Name)
		groups = append(groups, g)
	}
	return &groups
}

func QueryPolicyIdsByGroup(group int) *[]int {
	rows, err := mysql.Query("select p_id from g_p_ref where g_id = ?", group)
	checkErr(err)
	idList := make([]int, 10)
	var id int
	for rows.Next() {
		rows.Scan(&id)
		idList = append(idList, id)
	}
	return &idList
}