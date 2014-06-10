package mylib

import (
	md "../model"
	"database/sql"
	"fmt"
	"github.com/martini-contrib/render"
	"net/http"
)

type My_Json map[string]interface{}

const (
	SUCCESS        int = 1
	FAIL           int = 0
	INTERNAL_ERROR int = -1
)

func CheckNil(r render.Render, x interface{}) bool {
	if x == nil {
		Answer(r, FAIL)
		return true
	}
	return false
}

func CheckNotNil(r render.Render, x interface{}) bool {
	if x != nil {
		Answer(r, FAIL)
		return false
	}
	return true
}
func CheckLength(r render.Render, x string) bool {
	if len(x) == 0 {
		Answer(r, FAIL)
		return true
	}
	return false
}

// func fileserver(m martini.Martini, root string) {
// 	wf := func(path string, f os.FileInfo, err error) error {
// 		if f == nil {
// 			return err
// 		}
// 		if f.IsDir() {
// 			return nil
// 		}
// 		tmp := ReadWholeFile(path)
// 		fmt.Println(tmp)
// 		m.Get("/dssf", Handler_Factory_Content(tmp))
// 		return nil
// 	}
// 	err := filepath.Walk(root, wf)
// 	if err != nil {
// 		fmt.Println("file server error")
// 	}
// }

////////////////////////////////////////for handler///////////////////////////////////////////
func Handler_Factory_Content(content string) func(req *http.Request, w http.ResponseWriter) {
	return func(req *http.Request, w http.ResponseWriter) {
		fmt.Fprintf(w, content)
	}
}

func BasicSearch(req *http.Request, r render.Render, str string, DB *sql.DB) {
	id := "%" + req.Form["bookid"][0] + "%"
	name := "%" + req.Form["bookname"][0] + "%"

	if len(req.Form["bookname"][0]) == 0 {
		if len(req.Form["bookid"][0]) == 0 {
			r.HTML(200, "test", "非法输入")
		} else {
			rows, _ := DB.Query("SELECT * FROM books WHERE bookid like $1", id)
			md.Search_Result(rows, r, str)
			return
		}
	} else {
		if len(req.Form["bookid"][0]) == 0 {
			rows, _ := DB.Query("SELECT * FROM books WHERE bookname like $1", name)
			md.Search_Result(rows, r, str)
			return
		} else {
			rows, _ := DB.Query("SELECT * FROM books WHERE bookid like $1 and bookname like $2", id, name)
			md.Search_Result(rows, r, str)
			return
		}
	}
}

func ListSearch(req *http.Request, r render.Render, str string, DB *sql.DB) {
	rows, _ := DB.Query("SELECT * FROM books ")
	md.Search_Result(rows, r, str)
	return
}

func Json(a interface{}) My_Json {
	return My_Json{"result": a, "txt": ""}
}
func Json2(a interface{}, st string) My_Json {
	return My_Json{"result": a, "txt": st}
}

func Answer(r render.Render, a interface{}) {
	r.JSON(200, Json(a))
}

func Answer_2(r render.Render, a interface{}, st string) {
	r.JSON(200, Json2(a, st))
}
