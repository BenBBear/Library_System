package main

import (
	"errors"
	"fmt"
	"strings"
)

type StringSlice []string

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

func main() {
	// x := ""
	// si := strings.Split(x, ",")
	// for i := 0; i < len(si); i++ {
	// 	fmt.Println("first here we have: " + si[i])
	// }
	// var s []string
	// s = append(s, "trouble")
	// s = append(s, "trouble")
	// s = append(s, "trouble")
	// s = append(s, "trouble")
	// xx := strings.Join(s, ",")

	// ss := strings.Split(xx, ",")
	// for i := 0; i < len(ss); i++ {
	// 	fmt.Println("here we have: " + ss[i])
	// // }
	// var x []string
	// xx := strings.Join(x, ",")
	// if xx == "" {
	// 	fmt.Println("s")
	// }
	// fmt.Println(xx)
	var idx int
	var p, q []string
	p = append(p, "123")
	p = append(p, "12pp3")
	p = append(p, "1233")
	q = append(q, "1")
	q = append(q, "2")
	q = append(q, "3")
	p, idx, _ = DeleteFromSlice(p, "123")
	fmt.Println(idx)
	q, _ = DeleteFromSlice_index(q, idx)
	qq := strings.Join(q, ",")
	if qq == "" {
		fmt.Println("b")
	}
	fmt.Println(qq)
}
