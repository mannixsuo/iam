package model

type User struct {
	Id   int64
	Name string
}

func CreateUser(u *User) int64 {
	result, e := mysql.Exec("insert into `user`(`name`) values (?)", u.Name)
	checkErr(e)
	i, e := result.LastInsertId()
	checkErr(e)
	return i
}

func QueryUser(id int) *User {
	rows, e := mysql.Query("select `name` from `user` where id = ?", id)
	checkErr(e)
	u := User{}
	if rows.Next() {
		rows.Scan(&u.Name)
	}
	return &u
}

func QueryPolicyIdsByUser(user int) *[]int {
	rows, err := mysql.Query("select p_id from u_p_ref where u_id = ?", user)
	checkErr(err)
	idList := make([]int, 10)
	var id int
	for rows.Next() {
		rows.Scan(&id)
		idList = append(idList, id)
	}
	return &idList
}