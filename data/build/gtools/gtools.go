package gtools

import (
	"bytes"
	"fmt"
	"goe2m/data/build/model"
	"goe2m/data/config"
	"log"
	"regexp"
	"strings"
	"text/template"

	"github.com/xxjwxc/public/tools"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const newLine string = "说明,字段名,字段类型,备注,"

// GenStruct 数据
type GenStruct struct {
	FileName            string
	ProjectName         string
	EntityName          string
	EntityNote          string
	EntityTableName     string
	EntityContent       string
	EntityToContent     string
	EntitySchemaContent string
}

// GenElement 数据
type GenElement struct {
	Name      string              // Name.元素名
	NameLower string              // 小写
	Type      string              // Type.类型标记
	Notes     string              // Notes.注释
	Tags      map[string][]string // tages.标记
}

// Execute 执行
func Execute() {
	tEntity, err := template.ParseFiles("templates/entity.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tModel, err := template.ParseFiles("templates/model.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tBll, err := template.ParseFiles("templates/bll.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tCtl, err := template.ParseFiles("templates/ctl.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tSchema, err := template.ParseFiles("templates/schema.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	f, err := excelize.OpenFile(config.GetInFilePath())
	if err != nil {
		fmt.Println(err)
	}

	rows := f.GetRows("Sheet1")

	var genStruct *GenStruct
	var genElements []GenElement
	for _, row := range rows {
		line := fmt.Sprintf("%s,%s,%s,%s,", row[0], row[1], row[2], row[3])

		isTable := line == newLine

		if ok, _ := regexp.MatchString(`^[^,]*([\(|（]+)[^,]*([a-zA-Z][a-zA-Z]+)_?([a-zA-Z]+)([\)|）]+)`, line); ok {

			if len(genElements) > 0 {
				doGen(tEntity, tModel, tBll, tCtl, tSchema, genStruct, genElements)
			}

			// fmt.Println(line)
			genElements = make([]GenElement, 0)

			tableName := row[1]
			tableNote := row[0]
			modelName := row[2]

			tableName = strings.ToUpper(model.GetCamelName(tableName))
			fileName := strings.ToLower(model.GetCamelName(modelName))

			genStruct = &GenStruct{FileName: fileName, ProjectName: config.GetProjectName(), EntityName: modelName, EntityTableName: tableName, EntityNote: tableNote}

		} else if !isTable {
			if line == ",,,," {
				continue
			}
			eName := row[1]
			etype := row[2]
			eNote := row[0]

			nameLower := ""

			if eName != "" {
				nameLower = strings.ToLower(model.GetCamelName(eName))
			}

			element := GenElement{Name: eName, Type: etype, Notes: eNote, NameLower: nameLower}
			genElements = append(genElements, element)
		}

	}

	if len(genElements) > 0 {
		doGen(tEntity, tModel, tBll, tCtl, tSchema, genStruct, genElements)
	}
}

// 生成 文件
func doGen(tEntity *template.Template, tModel *template.Template, tBll *template.Template, tCtl *template.Template, tSchema *template.Template, genStruct *GenStruct, genElements []GenElement) {
	fmt.Println(genStruct.EntityTableName, genStruct.EntityNote, genStruct.EntityName)

	content := ""
	toContent := ""
	schemaContent := ""
	for _, v := range genElements {
		pGorm, pType, pjson := genGorm(&v)
		content += fmt.Sprintf("%s %s `%s` // %s \n", v.Name, pType, pGorm, v.Notes)
		toContent += fmt.Sprintf("%s: a.%s,\n", v.Name, v.Name)
		schemaContent += fmt.Sprintf("%s %s `%s` // %s \n", v.Name, pType, pjson, v.Notes)
	}
	genStruct.EntityContent = content
	genStruct.EntityToContent = toContent
	genStruct.EntitySchemaContent = schemaContent

	// 输出到buf
	buf := new(bytes.Buffer)
	tEntity.Execute(buf, genStruct) // 执行模板的替换
	writeFile("entity", "e_", genStruct.FileName, buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tModel.Execute(buf, genStruct) // 执行模板的替换
	writeFile("model", "m_", genStruct.FileName, buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tBll.Execute(buf, genStruct) // 执行模板的替换
	writeFile("bll", "b_", genStruct.FileName, buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tCtl.Execute(buf, genStruct) // 执行模板的替换
	writeFile("ctl", "c_", genStruct.FileName, buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tSchema.Execute(buf, genStruct) // 执行模板的替换
	writeFile("schema", "s_", genStruct.FileName, buf.String())

	// path, _ := writeFile("e_", genStruct.FileName, buf.String())

	// fmt.Println("formatting differs from goimport's:")
	// cmd, _ := exec.Command("goimports", "-l", "-w", path).Output()
	// fmt.Println(string(cmd))

	// fmt.Println("formatting differs from gofmt's:")
	// cmd, _ := exec.Command("gofmt", "-l", "-w", path).Output()
	// fmt.Println(string(cmd))
}

func genGorm(v *GenElement) (string, string, string) {
	fmt.Println(v.Name, v.NameLower, v.Type, v.Notes)
	gorm := fmt.Sprintf("column:%s;", v.NameLower)

	stype, len := getTypeName(strings.ToLower(v.Type))

	if len != 0 {
		gorm += fmt.Sprintf("size:%d;", len)
		if len == 36 {
			gorm += "index;"
		}
	}

	// fmt.Println(v.Name, stype)

	return fmt.Sprintf("gorm:\"%s\"", gorm), stype, fmt.Sprintf("json:\"%s\" swaggo:\"false,%s\"", v.NameLower, v.Notes)
}

// 保存文件
func writeFile(stype, prefix, fname, content string) (string, bool) {
	path := fmt.Sprintf("%s/%s/%s%s.go", config.GetOutDir(), stype, prefix, fname)
	return path, tools.WriteFile(path, []string{content}, true)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
