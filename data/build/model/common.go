package model

import "github.com/xxjwxc/public/mybigcamel"

// GetCamelName Big Hump or Capital Letter.大驼峰或者首字母大写
func GetCamelName(name string) string {
	return mybigcamel.UnMarshal(name)
}

// titleCase title case.首字母大写
func titleCase(name string) string {
	vv := []rune(name)
	if len(vv) > 0 {
		if bool(vv[0] >= 'a' && vv[0] <= 'z') { // title case.首字母大写
			vv[0] -= 32
		}
	}

	return string(vv)
}

// // getTypeName Type acquisition filtering.类型获取过滤
// func getTypeName(name string) string {
// 	// Precise matching first.先精确匹配
// 	if v, ok := cnf.TypeMysqlDicMp[name]; ok {
// 		return v
// 	}

// 	// Fuzzy Regular Matching.模糊正则匹配
// 	for k, v := range cnf.TypeMysqlMatchMp {
// 		if ok, _ := regexp.MatchString(k, name); ok {
// 			return v
// 		}
// 	}

// 	panic(fmt.Sprintf("type (%v) not match in any way.maybe need to add on (https://gormt/blob/master/data/view/cnf/def.go)", name))
// }
