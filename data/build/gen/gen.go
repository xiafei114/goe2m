package gen

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
	EntityNameLower     string
	EntityNote          string
	EntityTableName     string
	EntityContent       string
	EntityToContent     string
	EntityToContentMap  string
	EntitySchemaContent string
	HTMLEntityContent   string
	HTMLEntitys         string
	HTMLElementContent  string
}

// GenElement 数据
type GenElement struct {
	Name      string              // Name.元素名
	NameLower string              // 小写
	Type      string              // Type.类型标记
	Notes     string              // Notes.注释
	Remarks   string              // Remarks.注释
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

	tInterface, err := template.ParseFiles("templates/interface.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tHTMList, err := template.ParseFiles("templates/html/list.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tHTMLForm, err := template.ParseFiles("templates/html/form.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tHTMLElementText, err := template.ParseFiles("templates/html/element/text.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	tHTMLJs, err := template.ParseFiles("templates/html/js.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	f, err := excelize.OpenFile(config.GetInFilePath())
	if err != nil {
		fmt.Println(err)
	}

	sheet := f.GetSheetMap()

	schemaContent := ""
	interfaceContent := ""

	for _, v := range sheet {
		rows := f.GetRows(v)

		var genStruct *GenStruct
		var genElements []GenElement
		for _, row := range rows {
			line := fmt.Sprintf("%s,%s,%s,%s,", row[0], row[1], row[2], row[3])

			isTable := line == newLine

			if ok, _ := regexp.MatchString(`^[^,]*([\(|（]+)[^,]*([a-zA-Z][a-zA-Z]+)_?([a-zA-Z]+)([\)|）]+)`, line); ok {

				if len(genElements) > 0 {
					pSchemaContent, pInterfaceContent := doGen(tEntity, tModel, tBll, tCtl, tSchema, tInterface, tHTMList, tHTMLForm, tHTMLElementText, tHTMLJs, genStruct, genElements)
					schemaContent += pSchemaContent + "\n"
					interfaceContent += pInterfaceContent + "\n"
				}

				// fmt.Println(line)
				genElements = make([]GenElement, 0)

				tableName := row[1]
				tableNote := row[0]
				modelName := row[2]
				entityNameLower := strings.ToLower(modelName[0:1]) + modelName[1:]

				tableName = strings.ToUpper(model.GetCamelName(tableName))
				fileName := strings.ToLower(model.GetCamelName(modelName))

				genStruct = &GenStruct{
					FileName:        fileName,
					ProjectName:     config.GetProjectName(),
					EntityName:      modelName,
					EntityTableName: tableName,
					EntityNote:      tableNote,
					EntityNameLower: entityNameLower,
				}

			} else if !isTable {
				if line == ",,,," {
					continue
				}
				eName := row[1]
				etype := row[2]
				eNote := row[0]
				eRemark := row[3]

				nameLower := ""

				if eName != "" {
					nameLower = strings.ToLower(model.GetCamelName(eName))
				}

				element := GenElement{Name: eName, Type: etype, Notes: eNote, NameLower: nameLower, Remarks: eRemark}
				genElements = append(genElements, element)
			}

		}

		if len(genElements) > 0 {
			pSchemaContent, pInterfaceContent := doGen(tEntity, tModel, tBll, tCtl, tSchema, tInterface, tHTMList, tHTMLForm, tHTMLElementText, tHTMLJs, genStruct, genElements)
			schemaContent += pSchemaContent + "\n"
			interfaceContent += pInterfaceContent + "\n"
		}
	}

	writeFile("schema", "s_project", "", "go", schemaContent)
	writeFile("interface", "m_project", "", "go", interfaceContent)
}

// 生成 文件
func doGen(tEntity, tModel, tBll, tCtl, tSchema, tInterface, tHTMList, tHTMLForm, tHTMLElementText, tHTMLJs *template.Template,
	genStruct *GenStruct, genElements []GenElement) (pSchemaContent string, pInterfaceContent string) {
	fmt.Println(genStruct.EntityTableName, genStruct.EntityNote, genStruct.EntityName)

	content := ""
	toContent := ""
	toContentMap := ""
	schemaContent := ""

	htmlElementContent := ""
	htmlEntityContent := ""
	htmlEntity := ""
	for _, v := range genElements {
		pGorm, pType, pjson := genGorm(&v)
		content += fmt.Sprintf("%s %s `%s` // %s  %s\n", v.Name, pType, pGorm, v.Notes, v.Remarks)
		toContent += fmt.Sprintf("%s: a.%s,\n", v.Name, v.Name)
		toContentMap += fmt.Sprintf(`item["%s"] = a.%s`, v.NameLower, v.Name) + "\n"
		schemaContent += fmt.Sprintf("%s %s `%s` // %s %s \n", v.Name, pType, pjson, v.Notes, v.Remarks)

		buf := new(bytes.Buffer)
		tHTMLElementText.Execute(buf, v) // 执行模板的替换
		htmlElementContent += buf.String()

		switch pType {
		case "string":
			htmlEntityContent += fmt.Sprintf(`        '%s': '',`+"\n", v.NameLower)
		case "int", "int64":
			htmlEntityContent += fmt.Sprintf(`        '%s': 0,`+"\n", v.NameLower)
		case "*time.Time":
			htmlEntityContent += fmt.Sprintf(`        '%s': '',`+"\n", v.NameLower)
		case "bool":
			htmlEntityContent += fmt.Sprintf(`        '%s': false,`+"\n", v.NameLower)
		case "float64", "float32":
			htmlEntityContent += fmt.Sprintf(`        '%s': 0,`+"\n", v.NameLower)
		default:
			htmlEntityContent += fmt.Sprintf(`        '%s': '',`+"\n", v.NameLower)
		}
		htmlEntity += fmt.Sprintf(`'%s', `, v.NameLower)
	}
	genStruct.EntityContent = content
	genStruct.EntityToContent = toContent
	genStruct.EntityToContentMap = toContentMap
	genStruct.EntitySchemaContent = schemaContent

	if htmlEntityContent[len(htmlEntityContent)-2:] == ",\n" {
		htmlEntityContent = htmlEntityContent[:len(htmlEntityContent)-2]
	}

	if htmlEntity[len(htmlEntity)-2:] == ", " {
		htmlEntity = htmlEntity[:len(htmlEntity)-2]
	}
	genStruct.HTMLEntityContent = htmlEntityContent
	genStruct.HTMLEntitys = htmlEntity
	genStruct.HTMLElementContent = htmlElementContent

	// 输出到buf
	buf := new(bytes.Buffer)
	tEntity.Execute(buf, genStruct) // 执行模板的替换
	writeFile("entity", "e_", genStruct.FileName, "go", buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tModel.Execute(buf, genStruct) // 执行模板的替换
	writeFile("model", "m_", genStruct.FileName, "go", buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tBll.Execute(buf, genStruct) // 执行模板的替换
	writeFile("bll", "b_", genStruct.FileName, "go", buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tCtl.Execute(buf, genStruct) // 执行模板的替换
	writeFile("ctl", "c_", genStruct.FileName, "go", buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tHTMList.Execute(buf, genStruct) // 执行模板的替换
	writeFile("vue", "", genStruct.EntityName+"List", "vue", buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tHTMLForm.Execute(buf, genStruct) // 执行模板的替换
	writeFile("vue", "", genStruct.EntityName+"Form", "vue", buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tHTMLJs.Execute(buf, genStruct) // 执行模板的替换
	writeFile("js", "", genStruct.EntityNameLower, "js", buf.String())

	// 输出到buf
	buf = new(bytes.Buffer)
	tSchema.Execute(buf, genStruct) // 执行模板的替换

	pSchemaContent = buf.String()

	// 输出到buf
	buf = new(bytes.Buffer)
	tInterface.Execute(buf, genStruct) // 执行模板的替换

	pInterfaceContent = buf.String()

	return
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
func writeFile(stype, prefix, fname, suffix, content string) (string, bool) {
	path := fmt.Sprintf("%s/%s/%s%s.%s", config.GetOutDir(), stype, prefix, fname, suffix)
	return path, tools.WriteFile(path, []string{content}, true)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
