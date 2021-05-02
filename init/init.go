package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type SqlConfig struct {
	Address  string `ini:"ip"`
	Port     int    `ini:"port"`
	UserName string `ini:"username"`
	Password string `ini:"password"`
	Open     bool   `other:"open"`
	Minstep  float32
}
type Redis struct {
	Name string  `ini:"name"`
	Time float32 `ini:"time"`
}
type Cfg struct {
	SqlConfig *SqlConfig `ini:"mysql"`
	Redis     *Redis     `ini:"redis"`
}

func ParseIni(cfgStr string, a interface{}) error {
	//判断传入数据是否为结构体指针类型，不是结构体指针会报错
	key := reflect.TypeOf(a)
	value := reflect.ValueOf(a)
	fmt.Println(key, value)
	if key.Kind() != reflect.Ptr || key.Elem().Kind() != reflect.Struct {
		return errors.New("非结构体指针")
	}
	lines := strings.Split(cfgStr, "\r\n")
	var structName string
	var subValue reflect.Value
	var subKey reflect.Type
	for index, line := range lines {
		line = strings.TrimSpace(line)
		//如果是注释则跳过
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || line == "" {
			continue
		}
		//如果以[]包裹,则记录为结构体名
		if line[0] == '[' {
			reg := regexp.MustCompile(`^\[[a-zA-Z0-9_]{2,6}\]$`)
			if reg.MatchString(line) {
				structName = string(line[1 : len(line)-1])
				for i := 0; i < value.Elem().NumField(); i++ {
					field := key.Elem().Field(i)
					if field.Tag.Get("ini") == structName || field.Name == structName {
						structName = field.Name
						break
					}
				}
				//写一年
				subValue = value.Elem().FieldByName(structName)
				subKey = value.Elem().FieldByName(structName).Type()

				fmt.Println(subKey, subValue)
			} else {
				return fmt.Errorf("%d行名称格式错误", index+1)
			}
			continue
		}
		if strings.Contains(line, "=") {
			reg := regexp.MustCompile(`([a-zA-Z0-9]*)=(.*)`)
			ret := reg.FindAllSubmatch([]byte(line), -1)
			k := string(ret[0][1])
			v := string(ret[0][2])

			for i := 0; i < subKey.Elem().NumField(); i++ {
				field := subKey.Elem().Field(i)
				if field.Tag.Get("ini") == k || field.Name == k {
					switch field.Type.Kind() {
					case reflect.String:
						subValue.Elem().Field(i).SetString(v)
					case reflect.Int:
						ret, err := strconv.ParseInt(v, 10, 32)
						if err != nil {
							return fmt.Errorf("%d行变量转换失败", index+1)
						}
						subValue.Elem().Field(i).SetInt(ret)
					case reflect.Float32:
						ret, err := strconv.ParseFloat(v, 32)
						if err != nil {
							return fmt.Errorf("%d行变量转换失败", index+1)
						}
						subValue.Elem().Field(i).SetFloat(ret)
					case reflect.Float64:
						ret, err := strconv.ParseFloat(v, 64)
						if err != nil {
							return fmt.Errorf("%d行变量转换失败", index+1)
						}
						subValue.Elem().Field(i).SetFloat(ret)
					case reflect.Bool:
						ret, err := strconv.ParseBool(v)
						if err != nil {
							return fmt.Errorf("%d行变量转换失败", index+1)
						}
						subValue.Elem().Field(i).SetBool(ret)
					default:
						return fmt.Errorf("%d行没有找到合适的转换", index+1)
					}

				}
			}
		} else {
			return fmt.Errorf("%d行名称格式错误", index+1)
		}

	}
	return nil
}
func main() {
	cfg := Cfg{
		SqlConfig: &SqlConfig{},
		Redis:     &Redis{},
	}
	fullText, err := ioutil.ReadFile("ini.cfg")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ParseIni(string(fullText), &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.SqlConfig)
	fmt.Println(cfg.Redis)
}
