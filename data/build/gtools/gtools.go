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
	FileName        string
	ProjectName     string
	EntityName      string
	EntityNote      string
	EntityTableName string
	EntityContent   string
	EntityToContent string
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
	t, err := template.ParseFiles("templates/entity_model.txt") // 找到其中需要替换的模板变量
	checkErr(err)

	f, err := excelize.OpenFile(config.GetInFilePath())
	if err != nil {
		fmt.Println(err)
	}

	rows := f.GetRows("Sheet1")

	var genStruct *GenStruct
	var genElements []GenElement
	// genElements := make([]GenElement, 0)
	for _, row := range rows {
		var line string
		// for _, colCell := range row {
		// fmt.Print(colCell, "\t")
		line = fmt.Sprintf("%s,%s,%s,%s,", row[0], row[1], row[2], row[3])

		isTable := line == newLine

		if ok, _ := regexp.MatchString(`^[^,]*([\(|（]+)[^,]*([a-zA-Z][a-zA-Z]+)_?([a-zA-Z]+)([\)|）]+)`, line); ok {

			if len(genElements) > 0 {
				doGenStruct(t, genStruct, genElements)
			}

			// fmt.Println(line)
			genElements = make([]GenElement, 0)

			tableName := row[1]
			tableNote := row[0]
			modelName := row[2]

			tableName = strings.ToUpper(model.GetCamelName(tableName))
			fileName := strings.ToLower(model.GetCamelName(modelName))

			genStruct = &GenStruct{FileName: fileName, ProjectName: "wms_platform", EntityName: modelName, EntityTableName: tableName, EntityNote: tableNote, EntityContent: "夕阳西下", EntityToContent: "夕阳西下"}

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

		// // 是新行
		// isNewModel := line == newLine

		// if isNewModel {
		// 	tableName := rows[index-1][1]
		// 	tableNote := rows[index-1][0]
		// 	modelName := rows[index-1][2]

		// 	tableName = strings.ToUpper(model.GetCamelName(tableName))

		// 	fileName := strings.ToLower(model.GetCamelName(modelName))

		// 	fmt.Println(tableName, tableNote, modelName)

		// 	for _, v := range genElements {
		// 		fmt.Println(v.Name, v.NameLower, v.Type, v.Notes)
		// 	}

		// 	data := GenStruct{ProjectName: "wms_platform", EntityName: modelName, EntityTableName: tableName, EntityNote: tableNote, EntityContent: "夕阳西下", EntityToContent: "夕阳西下"}

		// 	// 输出到buf
		// 	buf := new(bytes.Buffer)
		// 	t.Execute(buf, data) // 执行模板的替换
		// 	writeFile("e_", fileName, buf.String())
		// 	genElements = make([]GenElement, 0)
		// } else {
		// 	if line == ",,,," {
		// 		continue
		// 	}

		// 	eName := row[1]
		// 	etype := row[2]
		// 	eNote := row[0]

		// 	nameLower := ""

		// 	if eName != "" {
		// 		nameLower = strings.ToLower(model.GetCamelName(eName))
		// 	}

		// 	element := GenElement{Name: eName, Type: etype, Notes: eNote, NameLower: nameLower}
		// 	genElements = append(genElements, element)
		// }
		// }

	}

	if len(genElements) > 0 {
		doGenStruct(t, genStruct, genElements)
	}

	// for _, v := range genElements {
	// 	fmt.Println(v.Name, v.NameLower, v.Type, v.Notes)
	// }
}

func doGenStruct(t *template.Template, genStruct *GenStruct, genElements []GenElement) {
	fmt.Println(genStruct.EntityTableName, genStruct.EntityNote, genStruct.EntityName)
	for _, v := range genElements {
		fmt.Println(v.Name, v.NameLower, v.Type, v.Notes)
	}
	// 输出到buf
	buf := new(bytes.Buffer)
	t.Execute(buf, genStruct) // 执行模板的替换
	writeFile("e_", genStruct.FileName, buf.String())
}

// 保存文件
func writeFile(prefix, fname, content string) bool {
	path := fmt.Sprintf("%s/%s%s.go", config.GetOutDir(), prefix, fname)
	return tools.WriteFile(path, []string{content}, true)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

const eModel string = `
package project

import (
	"context"

	schema "{{ project_name }}/internal/app/project/schema/project"
	"{{ project_name }}/pkg/gormplus"
	"{{ project_name }}/internal/app/project/model/gorm/entity"
)

// {{ entity_name }} 系统设置
type {{ entity_name }} struct {
	entity.Model
	{{ entity_content }}
}

// TableName 表名
func (a {{ entity_name }}) TableName() string {
	return a.Model.TableName("{{ entity_table_name }}")
}

// Get{{ entity_name }}DB 获取{{ entity_name }}存储
func Get{{ entity_name }}DB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	return entity.GetDBWithModel(ctx, defDB, {{ entity_name }}{})
}

// Schema{{ entity_name }} {{ entity_name }}对象
type Schema{{ entity_name }} schema.{{ entity_name }}

func (a {{ entity_name }}) String() string {
	return entity.ToString(a)
}

// {{ entity_name }}s {{ entity_name }}列表
type {{ entity_name }}s []*{{ entity_name }}

// ToSchema{{ entity_name }}s 转换为{{ entity_name }}对象列表
func (a {{ entity_name }}s) ToSchema{{ entity_name }}s() []*schema.{{ entity_name }} {
	list := make([]*schema.{{ entity_name }}, len(a))
	for i, item := range a {
		list[i] = item.ToSchema{{ entity_name }}()
	}
	return list
}

// To{{ entity_name }} 转换为{{ entity_name }}实体
func (a Schema{{ entity_name }}) To{{ entity_name }}() *{{ entity_name }} {
	item := &{{ entity_name }}{
		RecordID:     a.RecordID,
		{{ entity_to_content }}
	}
	return item
}

// ToSchema{{ entity_name }} 转换为{{ entity_name }}对象
func (a {{ entity_name }}) ToSchema{{ entity_name }}() *schema.{{ entity_name }} {
	item := &schema.{{ entity_name }}{
		RecordID:     a.RecordID,
		{{ entity_to_content }}
	}
	return item
}

`
