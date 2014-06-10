package mylib

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/martini-contrib/render"
	"io"
	"os"
	"strings"
	"syscall"
)

type StringSlice []string
type Status int

const (
	ADMIN  = 1
	COMMON = 0
)

type User_Data struct {
	Who      User
	Num      int
	Book     StringSlice
	BookName StringSlice
	Nil      bool
}

// type Bookpair struct {
// 	Book     string
// 	BookName string
// }

type User struct {
	Name, Password string
	Sha256         string
	Who            Status
}

type User_Channel chan User
type User_Map struct {
	Content map[string]User
}

func ReadWholeFile(fn string) (rt string) {
	f, _ := os.Open(fn)
	fi, _ := f.Stat()
	data, _ := syscall.Mmap(int(f.Fd()), 0, int(fi.Size()),
		syscall.PROT_READ, syscall.MAP_PRIVATE)
	rt = string(data)
	f.Close()
	return rt
}

func (Online_People *User_Map) Add_User(usr_detail User) {
	Online_People.Content[usr_detail.Sha256] = usr_detail
	fmt.Println("From Add User:  " + usr_detail.Sha256)
}
func (Online_People *User_Map) Del_User(usr string) {
	delete(Online_People.Content, usr)
}

func (Online_People *User_Map) Log_inOut(uc User_Channel) {
	for {
		select {
		case user_new := <-uc:
			go Online_People.Add_User(user_new)
		}
	}
}
func (Online_People *User_Map) Contain(u string) bool {
	_, ok := Online_People.Content[u]
	// P("dddddddddddddddddddddddddddddddddddddddddddddddddddddddd")
	// fmt.Println(u)
	// P(ok)
	return ok
}
func (Online_People *User_Map) Admin(u string) bool {
	tmp, ok := Online_People.Content[u]
	if !ok {
		return false
	} else if tmp.Who != ADMIN {
		return false
	}
	return true
}
func (O User) String() string {
	return O.Name + " " + O.Password
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func P(a ...interface{}) {
	fmt.Println(a)
}

func SHA256(str string) string {
	EnCoder := sha256.New()
	io.WriteString(EnCoder, str)
	return string(EnCoder.Sum(nil)[:])
}

func DeleteFromSlice(sl StringSlice, str string) (results StringSlice, idx int, err error) {
	for i := 0; i < len(sl); i++ {
		if sl[i] == str {
			if i == 0 {
				if len(sl) > 1 {
					return sl[1:], i, nil
				} else {
					var tmp StringSlice
					return tmp, i, nil
				}
			} else {

				return append(sl[0:i-1], sl[i:]...), i, nil
			}
		}
	}
	results = nil
	idx = -1
	err = errors.New("fail")
	return
}
func DeleteFromSlice_index(sl StringSlice, idx int) (results StringSlice, err error) {
	if idx == 0 {
		if len(sl) == 1 {
			var tmp StringSlice
			return tmp, nil
		} else {
			return sl[1:], nil
		}
	} else {
		return append(sl[0:idx], sl[idx+1:]...), nil
	}
	err = errors.New("fail")
	results = nil
	return
}

func NewUser(username string, passwd string) *User {
	return &User{Name: username, Password: passwd, Sha256: SHA256(passwd + username), Who: COMMON}
}

func CheckErr(r render.Render, err error, str string) bool {
	if err != nil {
		r.HTML(200, "test", str)
		return true
	}
	return false
}

func BookList(str string) StringSlice {
	var tmp StringSlice
	if str == "" {
		return tmp
	} else {
		tmp = strings.Split(str, ",")
	}
	return tmp
}

/////////////////////////////////////////////////////////////////////////////data base interactive///////////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (Online_People *User_Map) BorrowBook(usr string, r render.Render, DB *sql.DB, bookid string, bookname string) {
	tmp_user, ok := Online_People.Content[usr]
	if !ok {
		Answer_2(r, FAIL, "请登录")
		return
	}
	var num int
	var book string
	var name string
	err := DB.QueryRow("SELECT num,book,bookname FROM users WHERE username = $1", tmp_user.Name).Scan(&num, &book, &name)
	CheckErr(r, err, "error1")
	if num >= 5 {
		Answer_2(r, FAIL, "你已经借满书了")
		return
	}

	bookset := BookList(book)
	nameset := BookList(name)
	bookset = append(bookset, bookid)
	nameset = append(nameset, bookname)
	book = strings.Join(bookset, ",")
	name = strings.Join(nameset, ",")
	_, err = DB.Exec("UPDATE users set book=$2,num=num+1,bookname=$3  where username = $1", tmp_user.Name, book, name)
	CheckErr(r, err, "error2")
	_, err = DB.Exec("UPDATE books set numberleft = numberleft-1  where bookid = $1", bookid)
	CheckErr(r, err, "error3")
	// fmt.Println(bookid)
	Answer_2(r, SUCCESS, "成功借书")
	return
}

func (Online_People *User_Map) ReturnBook(usr string, r render.Render, DB *sql.DB, bookid string, bookname string) {
	tmp_user, ok := Online_People.Content[usr]
	if !ok {
		Answer_2(r, FAIL, "请登录")
		return
	}
	var num, idx int
	var book string
	var name string
	err := DB.QueryRow("SELECT num,book,bookname FROM users WHERE username = $1", tmp_user.Name).Scan(&num, &book, &name)
	CheckErr(r, err, "error")
	if num == 0 {
		Answer_2(r, FAIL, "同学你没有借这本书")
		return
	}
	bookset := strings.Split(book, ",")
	nameset := strings.Split(name, ",")
	bookset, idx, err = DeleteFromSlice(bookset, bookid)
	CheckErr(r, err, "error1")
	nameset, err = DeleteFromSlice_index(nameset, idx)
	CheckErr(r, err, "erroridx")
	book = strings.Join(bookset, ",")
	name = strings.Join(nameset, ",")
	_, err = DB.Exec("UPDATE users set book=$2,num=num-1,bookname=$3  where username = $1", tmp_user.Name, book, name)
	if CheckErr(r, err, "error2") {
		return
	}
	_, err = DB.Exec("UPDATE books set numberleft = numberleft+1  where bookid = $1", bookid)
	if CheckErr(r, err, "同学在干什么") {
		return
	}
	Answer_2(r, SUCCESS, "成功还书")
	return
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (Online_People *User_Map) Home(usr string, r render.Render, DB *sql.DB) {
	tmp_user := Online_People.Content[usr]
	var num int
	var book, bookname string
	err := DB.QueryRow("SELECT num,book,bookname FROM users WHERE username = $1", tmp_user.Name).Scan(&num, &book, &bookname)
	CheckErr(r, err, "internal error")
	if len(book) == 0 {
		r.HTML(200, "home", User_Data{tmp_user, num, strings.Split(book, ","), strings.Split(bookname, ","), false})
	} else {
		r.HTML(200, "home", User_Data{tmp_user, num, strings.Split(book, ","), strings.Split(bookname, ","), true})
	}
	return
}
