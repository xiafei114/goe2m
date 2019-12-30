package gtools

import (
	"bytes"
	"fmt"
	"goe2m/data/build/model"
	"goe2m/data/config"
	"log"
	"strings"
	"text/template"

	"github.com/xxjwxc/public/tools"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const newLine string = "说明,字段名,字段类型,备注,"

// Data 数据
type Data struct {
	ProjectName     string
	EntityName      string
	EntityNote      string
	EntityTableName string
	EntityContent   string
	EntityToContent string
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
	for index, row := range rows {
		var line string
		for _, colCell := range row {
			// fmt.Print(colCell, "\t")
			line = line + colCell + ","

			// 是新行
			isNewModel := line == newLine

			if isNewModel {
				tableName := rows[index-1][1]
				tableNote := rows[index-1][0]
				modelName := rows[index-1][2]

				tableName = strings.ToUpper(model.GetCamelName(tableName))

				fileName := strings.ToLower(model.GetCamelName(modelName))

				fmt.Println(tableName, tableNote, modelName)
				data := Data{ProjectName: "wms_platform", EntityName: modelName, EntityTableName: tableName, EntityNote: tableNote, EntityContent: "夕阳西下", EntityToContent: "夕阳西下"}

				// 输出到buf
				buf := new(bytes.Buffer)
				t.Execute(buf, data) // 执行模板的替换
				writeFile("e_", fileName, buf.String())
			}
		}

	}
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
