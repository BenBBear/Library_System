package mylib

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
)

type Littlebox struct {
	Bookid   string
	Bookname string
	Left_num int
}

func (l Littlebox) String() string {
	return l.Bookid + " " + l.Bookname
}

func Search_Result(rows *sql.Rows, r render.Render, str string) {
	tmp_container := []Littlebox{}
	for rows.Next() {
		var bookid string
		var bookname string
		var left_num int
		err := rows.Scan(&bookname, &bookid, &left_num)
		if err != nil {
			fmt.Errorf("DATABASE ERROR")
		}
		tmp := Littlebox{bookid, bookname, left_num}
		fmt.Println(tmp)
		tmp_container = append(tmp_container, tmp)
	}
	r.HTML(200, str, tmp_container)
}
func Left_Number(rows *sql.Rows) int {
	var left_num int
	for rows.Next() {
		err := rows.Scan(&left_num)
		if err != nil {
			fmt.Errorf("DATABASE ERROR")
		}
	}
	return left_num
}
