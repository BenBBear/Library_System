// main.go
package main

import (
	hd "./handlers"
	ml "./utility"
	// md "./model"
	"database/sql"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	// "github.com/martini-contrib/auth"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	// "net/http"
	"fmt"
	"os"
	"path/filepath"
)

func fileserver(root string, pre string) {
	l := len(pre)
	wf := func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		tmp := ml.ReadWholeFile(path)
		m.Get("/"+path[l-1:], ml.Handler_Factory_Content(tmp))
		return nil
	}
	err := filepath.Walk(root, wf)
	if err != nil {
		fmt.Println("file server error")
	}
}

var (
	Channal_Size    int = 100
	Global_User_Map     = ml.User_Map{make(map[string]ml.User)}
	Global_Channel      = make(ml.User_Channel, Channal_Size)
	m                   = martini.Classic()
)

//暂时现不实现登录，使用middleware实现起来不难，而且这是过早优化
func main() {
	DB, _ := sql.Open("postgres", "postgres://beviszhang:12345678@localhost/testcourse")
	m.Use(render.Renderer())
	m.Use(hd.Create_DB_Service(DB))
	My_Session := sessions.NewCookieStore([]byte("authentication"))
	m.Use(sessions.Sessions("Session", My_Session))
	m.Use(hd.Create_Login_Service(Global_Channel, Global_User_Map))
	//////////////////////////////////////////////////////////////////////////////////////////////////
	fileserver("./templates/Pack", "./templates")
	/////////////////////////////////////////////////////////////////////////////////
	m.Get("/", hd.Static_Page(200, "index", nil))
	///////////////////////////////////////////////////////////////////////////////////
	m.Get("/signup", hd.Static_Page(200, "signup", "gosignup/")) //注册功能 tested
	m.Post("/signup/gosignup", hd.SignUpHandler())
	///////////////////////////////////////////////////////////////////////////////////
	m.Get("/user_login", hd.Static_Page(200, "login", "login/")) //登录功能 tested
	m.Post("/user_login/login", hd.LoginHandler())

	///////////////////////////////////////////////////////////////////////////////////
	m.Get("/logout", hd.LogOutHandler()) // 注销功能 tested
	///////////////////////////////////////////////////////////////////////////////////
	m.Get("/search", hd.Static_Page(200, "search", "/search/result/")) //匿名搜索功能，tested
	m.Post("/search/result", hd.SearchHandler("search_result"))
	///////////////////////////////////////////////////////////////////////////////////
	m.Get("/home", hd.HomeHandler())             //主页功能，tested
	m.Post("/home/goreturn", hd.ReturnHandler()) //还书功能，tested
	m.Get("/home/search", hd.Static_Page(200, "search", "/home/search/result/"))
	m.Post("/home/search/result", hd.SearchHandler("user_borrow"))
	m.Post("/home/search/result/goborrow", hd.BorrowHandler()) //借书功能 tested

	///////////////////////////////////////////////////////////////////////////////////
	m.Get("/admin", hd.Static_Page(200, "admin", nil)) //管理员功能 tested
	m.Get("/admin/add", hd.Static_Page(200, "admin_add", nil))
	m.Post("/admin/add/goadd", hd.Admin_AddHandler())

	m.Get("/admin/login", hd.Static_Page(200, "admin_login", "judge/"))
	m.Post("/admin/login/judge", hd.Admin_LoginHandler())
	m.Get("/admin/search", hd.Static_Page(200, "search", "result/"))
	m.Post("/admin/search/result", hd.Admin_Edit_SearchHandler("admin_edit"))
	m.Post("/admin/search/result/goedit", hd.Admin_ModifyHandler())
	m.Get("/admin/list", hd.Admin_ListHandler("admin_edit"))
	m.Post("/admin/list/goedit", hd.Admin_ModifyHandler())
	///////////////////////////////////////////////////////////////////////////////////
	go Global_User_Map.Log_inOut(Global_Channel)

	m.Run()

	defer func() {
		DB.Close()
	}()
}
