package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type SqlConfig struct {
	Address  string `init:"ip"`
	Port     int    `init:"port"`
	UserName string `init:"username"`
	Password string `init:"password"`
	Open     bool   `other:"open"`
	Minstep  float32
}

func (s *SqlConfig) compileConfig(str string) bool {
	//文件中的换行要使用\r\n分割，不同平台可能不同
	lines := strings.Split(str, "\r\n")
	//分别拿到键和值的信息
	key := reflect.TypeOf(s)
	value := reflect.ValueOf(s)
	//处理文件，把需要的字段提取出来
	m := make(map[string]string, len(lines))
	for _, line := range lines {
		if strings.Contains(line, "=") {
			kv := strings.Split(line, "=")
			m[kv[0]] = kv[1]
		}
	}
	//v是结构体指针，必须先使用Elem()获得结构体，再调用NumField()获得结构体的变量个数
	// fmt.Println(v.Elem().NumField())
	// fmt.Println(m)

	//结构体中一共有多少个字段，需要用Elem才能通过指针取到结构体.
	//这里使用value.Elem().NumField()或key.Elem().NumField()可以拿到相同的结果
	fd := value.Elem().NumField()
	for i := 0; i < fd; i++ {
		//获得键的tag,如果没有设置tag的话，直接使用键的名字
		k := key.Elem().Field(i).Tag.Get("init")
		if k == "" {
			k = key.Elem().Field(i).Name
		}
		//获得某个键下值得具体类型(在reflec中得类型)
		switch key.Elem().Field(i).Type.Kind() {
		case reflect.String:
			//使用SetXXX方法对 值 进行赋值。使用key会报错
			value.Elem().Field(i).SetString(m[k])
		case reflect.Int:
			//使用strconv包将字符串转为其他类型
			interger, err := strconv.ParseInt(m[k], 10, 32)
			if err != nil {
				fmt.Println("一个整数不符合规范", err)
			}
			value.Elem().Field(i).SetInt(interger)
		case reflect.Float32:
			fmt.Printf(m[k])
			float, err := strconv.ParseFloat(m[k], 32)
			if err != nil {
				fmt.Println("一个浮点数不符合规范", err)
			}
			value.Elem().Field(i).SetFloat(float)
		case reflect.Float64:
			fmt.Printf(m[k])
			float, err := strconv.ParseFloat(m[k], 64)
			if err != nil {
				fmt.Println("一个浮点数不符合规范", err)
			}
			value.Elem().Field(i).SetFloat(float)
		case reflect.Bool:
			b, err := strconv.ParseBool(m[k])
			if err != nil {
				fmt.Println("一个布尔值不符合规范", err)
			}
			value.Elem().Field(i).SetBool(b)
		}
	}
	return true
}
func main() {
	//读取文件
	content, err := ioutil.ReadFile("init.cof")
	if err != nil {
		fmt.Printf("open file failed, err = %v", err)
	}
	var sc SqlConfig
	ok := sc.compileConfig(string(content))
	if ok {
		fmt.Print(sc)
	} else {
		fmt.Println("解析失败")
	}
}
