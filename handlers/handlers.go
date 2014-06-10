package handlers

import (
	// md "../model"
	ml "../utility"
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/lib/pq"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
	"strconv"
)

type handler func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB)
type simple_handler func(req *http.Request, w http.ResponseWriter)
type Channal_handler func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, user_chan ml.User_Channel)
type Map_handler func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map)
type Service func(c martini.Context)
type Static_Page_handler func(r render.Render)

const (
	// SUCCESS        string = "1"
	// FAIL           string = "0"
	// INTERNAL_ERROR string = "-1"
	SUCCESS        int = 1
	FAIL           int = 0
	INTERNAL_ERROR int = -1
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Handler_Factory_Content(content string) simple_handler {
	return func(req *http.Request, w http.ResponseWriter) {
		fmt.Fprintf(w, content)
	}
}
func Create_DB_Service(DB *sql.DB) Service {
	return func(c martini.Context) {
		c.Map(DB)
		c.Next()
	}
}

func Create_Login_Service(user_chan ml.User_Channel, Global_User_Map ml.User_Map) Service {
	return func(c martini.Context) {
		c.Map(user_chan)
		c.Map(Global_User_Map)
		c.Next()
	}
}

// func Create_Hash_Service(EnCoder hash.Hash) Service {
// 	return func(c martini.Context) {
// 		c.Map(EnCoder)
// 		c.Next()
// 	}
// }

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func Static_Page(rn int, str string, obj interface{}) Static_Page_handler {
	return func(r render.Render) {
		r.HTML(rn, str, obj)
	}
} //http return num, template to use,struct for the template
func SearchHandler(str string) handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB) {
		req.ParseForm()
		ml.BasicSearch(req, r, str, DB)
		return
	}
}

func Admin_Edit_SearchHandler(str string) Map_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		req.ParseForm()
		x := ses.Get("userid")
		if x == nil {
			r.HTML(200, "test", "管理员同学你还没登录")
		} else {
			if !Global_User_Map.Admin(x.(string)) {
				r.HTML(200, "test", "管理员同学是你吗？")
				return
			}
		}
		ml.BasicSearch(req, r, str, DB)
		return
	}
}
func Admin_ListHandler(str string) Map_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		req.ParseForm()
		x := ses.Get("userid")
		if x == nil {
			r.HTML(200, "test", "管理员同学你还没登录")
		} else {
			if !Global_User_Map.Admin(x.(string)) {
				r.HTML(200, "test", "管理员同学是你吗？")
				return
			}
		}
		ml.ListSearch(req, r, str, DB)
		return
	}
}
func Admin_ModifyHandler() Map_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		req.ParseForm()
		x := ses.Get("userid")
		if x == nil {
			ml.Answer_2(r, FAIL, "管理员你还没登录")
		} else {
			if !Global_User_Map.Admin(x.(string)) {
				ml.Answer_2(r, FAIL, "管理员是你吗")
				return
			}
		}
		if ml.CheckNil(r, req.Form["bookid"][0]) ||
			ml.CheckNil(r, req.Form["bookname"][0]) ||
			ml.CheckLength(r, req.Form["bookid"][0]) ||
			ml.CheckLength(r, req.Form["bookname"][0]) {
			ml.Answer(r, FAIL)
			return
		}
		if len(req.Form["delete"]) != 0 {
			_, err := DB.Exec("DELETE FROM  books   where bookid = $1 and bookname = $2", req.Form["bookid"][0], req.Form["bookname"][0])
			if err != nil {
				ml.Answer_2(r, INTERNAL_ERROR, "内部错误")
			}
		} else {
			_, err := DB.Exec("UPDATE books set numberleft =$3  where bookid = $1 and bookname = $2", req.Form["bookid"][0], req.Form["bookname"][0], req.Form["left_num"][0])
			if err != nil {
				ml.Answer_2(r, INTERNAL_ERROR, "内部错误")
			}

		}
		ml.Answer_2(r, SUCCESS, "修改成功")
		return
	}
}

func Admin_AddHandler() Map_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		req.ParseForm()
		x := ses.Get("userid")
		if x == nil {
			ml.Answer_2(r, FAIL, "管理员你还没登录")
		} else {
			if !Global_User_Map.Admin(x.(string)) {
				ml.Answer_2(r, FAIL, "管理员是你吗")
				return
			}
		}
		if ml.CheckNil(r, req.Form["bookid"][0]) ||
			ml.CheckNil(r, req.Form["bookname"][0]) ||
			ml.CheckLength(r, req.Form["bookid"][0]) ||
			ml.CheckLength(r, req.Form["bookname"][0]) {
			ml.Answer_2(r, FAIL, "输入非法")
			return
		}
		if _, err := strconv.Atoi(req.Form["left_num"][0]); err != nil {
			ml.Answer_2(r, FAIL, "输入非法")
			return
		}
		_, err := DB.Exec("INSERT INTO books VALUES($1,$2,$3)", req.Form["bookname"][0], req.Form["bookid"][0], req.Form["left_num"][0])
		if err != nil {
			ml.Answer_2(r, INTERNAL_ERROR, "内部错误")
		}
		ml.Answer_2(r, SUCCESS, "修改成功")
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func BorrowHandler() Map_handler { //提交的表单完整性监察很重要，所以必需要用那个binding form 的库,这次是很好的教训
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		x := ses.Get("userid")
		if x == nil {
			ml.Answer_2(r, FAIL, "你还没登录")
		} else {
			if !Global_User_Map.Contain(x.(string)) {
				ml.Answer_2(r, FAIL, "你还没登录")
				return
			}
		}
		req.ParseForm()
		if ml.CheckNil(r, req.Form["bookid"][0]) ||
			ml.CheckNil(r, req.Form["bookname"][0]) ||
			ml.CheckLength(r, req.Form["bookid"][0]) ||
			ml.CheckLength(r, req.Form["bookname"][0]) {
			ml.Answer_2(r, FAIL, "输入非法")
			return
		}
		var Left_N int
		err := DB.QueryRow("SELECT numberleft FROM books WHERE bookid = $1 or bookname = $2", req.Form["bookid"][0], req.Form["bookname"][0]).Scan(&Left_N)
		switch {
		case err == sql.ErrNoRows:
			ml.Answer_2(r, FAIL, "同学在干什么？")
		case err != nil:
			ml.Answer_2(r, INTERNAL_ERROR, "内部错误")
		case Left_N <= 0:
			ml.Answer_2(r, FAIL, "书不够啦")
		default:
			Global_User_Map.BorrowBook(x.(string), r, DB, req.Form["bookid"][0], req.Form["bookname"][0])
		}
	}
}
func ReturnHandler() Map_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		x := ses.Get("userid")
		if x == nil {
			ml.Answer_2(r, FAIL, "你还没登录")
		} else {
			if !Global_User_Map.Contain(x.(string)) {
				ml.Answer_2(r, FAIL, "你还没登录")
				return
			}
		}
		req.ParseForm()
		if ml.CheckNil(r, req.Form["bookid"][0]) ||
			ml.CheckNil(r, req.Form["bookname"][0]) ||
			ml.CheckLength(r, req.Form["bookid"][0]) ||
			ml.CheckLength(r, req.Form["bookname"][0]) {
			ml.Answer_2(r, FAIL, "输入非法")
			return
		}
		var Left_N int
		err := DB.QueryRow("SELECT numberleft FROM books WHERE bookid = $1", req.Form["bookid"][0]).Scan(&Left_N)
		switch {
		case err == sql.ErrNoRows:
			ml.Answer_2(r, FAIL, "输入非法")
		case err != nil:
			ml.Answer_2(r, INTERNAL_ERROR, "内部错误")
		default:
			Global_User_Map.ReturnBook(x.(string), r, DB, req.Form["bookid"][0], req.Form["bookname"][0])
		}
	}
}

////////////////////// r.JSON(200, map[string]interface{}{"hello": "world"})///////////////////////////////////////////////login///////////////////////////////////////////////////////////////////////////////////////////////////
func LoginHandler() Channal_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, user_chan ml.User_Channel) {
		req.ParseForm()
		if ml.CheckNil(r, req.Form["password"][0]) ||
			ml.CheckNil(r, req.Form["username"][0]) ||
			ml.CheckLength(r, req.Form["password"][0]) ||
			ml.CheckLength(r, req.Form["username"][0]) {
			ml.Answer(r, FAIL)
			return
		}
		var username string
		err := DB.QueryRow("SELECT username FROM users WHERE username = $1 and passwd = $2", req.Form["username"][0], req.Form["password"][0]).Scan(&username)
		switch {
		case err == sql.ErrNoRows:
			ml.Answer(r, FAIL)
			return
		case err != nil:
			ml.Answer(r, INTERNAL_ERROR)
			return
		default:
			tmp_u := ml.NewUser(req.Form["username"][0], req.Form["password"][0])
			user_chan <- *tmp_u
			ses.Set("userid", tmp_u.Sha256)
			// fmt.Fprintf(w, "sfdffsx")
			ml.Answer(r, SUCCESS)
		}
		return
	}
}

func Admin_LoginHandler() Channal_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, user_chan ml.User_Channel) {
		req.ParseForm()
		if ml.CheckNil(r, req.Form["password"][0]) ||
			ml.CheckNil(r, req.Form["username"][0]) ||
			ml.CheckLength(r, req.Form["password"][0]) ||
			ml.CheckLength(r, req.Form["username"][0]) {
			ml.Answer(r, FAIL)
			return
		}
		var username string
		err := DB.QueryRow("SELECT username FROM administrater WHERE username = $1 and passwd = $2", req.Form["username"][0], req.Form["password"][0]).Scan(&username)
		switch {
		case err == sql.ErrNoRows:
			ml.Answer(r, FAIL)
			return
		case err != nil:
			ml.Answer(r, INTERNAL_ERROR)
			return
		default:
			tmp_u := ml.NewUser(req.Form["username"][0], req.Form["password"][0])
			tmp_u.Who = ml.ADMIN
			user_chan <- *tmp_u
			ses.Set("userid", tmp_u.Sha256)
			ml.Answer(r, SUCCESS)
		}
		return
	}
}
func LogOutHandler() Map_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		x := ses.Get("userid")
		if x == nil {
			r.HTML(200, "test", "error")
		} else {
			if Global_User_Map.Contain(x.(string)) {
				Global_User_Map.Del_User(x.(string))
				r.HTML(200, "test", "注销成功")
			} else {
				r.HTML(200, "test", "你还没登录")
			}
		}
	}
}

func SignUpHandler() handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB) {
		req.ParseForm()
		if ml.CheckNil(r, req.Form["password"][0]) ||
			ml.CheckNil(r, req.Form["username"][0]) ||
			ml.CheckLength(r, req.Form["password"][0]) ||
			ml.CheckLength(r, req.Form["username"][0]) {
			ml.Answer(r, FAIL)
			return
		}
		result, err := DB.Exec(fmt.Sprintf("INSERT INTO %s VALUES ($1,$2,0,$3,$4)", pq.QuoteIdentifier("users")), req.Form["username"][0], req.Form["password"][0], "", "")
		if err != nil {
			ml.Answer(r, INTERNAL_ERROR)
			return
		}
		if result == nil {
			ml.Answer(r, FAIL)
			return
		}
		ml.Answer(r, SUCCESS)
		return
	}
}

func HomeHandler() Map_handler {
	return func(req *http.Request, w http.ResponseWriter, r render.Render, param martini.Params, ses sessions.Session, DB *sql.DB, Global_User_Map ml.User_Map) {
		x := ses.Get("userid")
		if x == nil {
			r.HTML(200, "test", "你还没登录OR你还未成为用户")
		} else {
			if Global_User_Map.Contain(x.(string)) {
				// ml.P("ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
				Global_User_Map.Home(x.(string), r, DB)
				return
			} else {
				// ml.P("ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
				r.HTML(200, "test", "你还没登录OR你还未成为用户")
			}
		}

	}
}
