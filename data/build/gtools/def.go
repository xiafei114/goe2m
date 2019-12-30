package gtools

import (
	"fmt"
	"regexp"
	"strconv"
)

// EImportsHead imports head options. import包含选项
var EImportsHead = map[string]string{
	"stirng":     `"string"`,
	"time.Time":  `"time"`,
	"gorm.Model": `"github.com/jinzhu/gorm"`,
}

// TypeMysqlDicMp Accurate matching type.精确匹配类型
var TypeMysqlDicMp = map[string]string{
	"int":                 "int",
	"bigint":              "int64",
	"varchar":             "string",
	"char":                "string",
	"date":                "*time.Time",
	"datetime":            "*time.Time",
	"bit(1)":              "bool",
	"tinyint(1)":          "bool",
	"tinyint(1) unsigned": "bool",
	"json":                "string",
	"text":                "string",
	"timestamp":           "*time.Time",
	"double":              "float64",
	"mediumtext":          "string",
	"longtext":            "string",
	"float":               "float32",
	"tinytext":            "string",
	"longblob":            "string",
}

// TypeMysqlMatchMp Fuzzy Matching Types.模糊匹配类型
var TypeMysqlMatchMp = map[string]string{
	`^(tinyint)[(]\d+[)]`:     "int8",
	`^(smallint)[(]\d+[)]`:    "int8",
	`^(int)[(]\d+[)]`:         "int",
	`^(bigint)[(]\d+[)]`:      "int64",
	`^(char)[(](\d+)[)]`:      "string",
	`^(varchar)[(](\d+)[)]`:   "string",
	`^(varbinary)[(]\d+[)]`:   "[]byte",
	`^(decimal)[(]\d+,\d+[)]`: "float64",
	`^(mediumint)[(](\d+)[)]`: "string",
	`^(double)[(]\d+,\d+[)]`:  "float64",
	`^(float)[(]\d+,\d+[)]`:   "float64",
	`^(numeric)[(]\d+,\d+[)]`: "float64",
}

// getTypeName Type acquisition filtering.类型获取过滤
func getTypeName(name string) (string, int) {
	// Precise matching first.先精确匹配
	if v, ok := TypeMysqlDicMp[name]; ok {
		return getNameLenght(v, "", name)
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for k, v := range TypeMysqlMatchMp {
		if ok, _ := regexp.MatchString(k, name); ok {
			return getNameLenght(v, k, name)
		}
	}

	panic(fmt.Sprintf("type (%v) not match in any way.maybe need to add", name))
}

func getNameLenght(value, k, name string) (string, int) {
	if value != "string" || k == "" {
		return value, 0
	}

	reg := regexp.MustCompile(k)
	values := reg.FindStringSubmatch(name)
	// fmt.Println(values)
	if len(values) > 2 {
		len, _ := strconv.Atoi(values[2])
		return value, len
	} else {
		return value, 0
	}
}
