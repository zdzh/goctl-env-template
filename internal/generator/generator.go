package generator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zdzh/goctl-env-template/internal/types"
)

type Generator struct {
	structs []types.ConfigStruct
}

func NewGenerator(structs []types.ConfigStruct) *Generator {
	return &Generator{
		structs: structs,
	}
}

func (g *Generator) Generate() string {
	var builder strings.Builder

	// 按结构体名称分组
	structMap := g.groupStructsByStruct()

	groupNames := make([]string, 0, len(structMap))
	for name := range structMap {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	for i, groupName := range groupNames {
		configStruct := structMap[groupName]

		if len(configStruct.Fields) > 0 {
			g.generateGroup(&builder, configStruct)

			// 不是最后一个分组时，添加两个空行
			if i < len(groupNames)-1 {
				builder.WriteString("\n\n")
			}
		}
	}

	return builder.String()
}

func (g *Generator) groupStructsByStruct() map[string]types.ConfigStruct {
	structMap := make(map[string]types.ConfigStruct)

	for _, configStruct := range g.structs {
		// 如果已经存在，合并字段
		if existing, ok := structMap[configStruct.Name]; ok {
			existing.Fields = append(existing.Fields, configStruct.Fields...)
			structMap[configStruct.Name] = existing
		} else {
			structMap[configStruct.Name] = configStruct
		}
	}

	return structMap
}

func (g *Generator) generateGroup(builder *strings.Builder, configStruct types.ConfigStruct) {
	// 添加结构体注释
	if configStruct.Comment != "" {
		cleanedComment := g.cleanComment(configStruct.Comment)
		builder.WriteString(fmt.Sprintf("# %s\n", cleanedComment))
	}

	groupTitle := g.formatGroupTitle(configStruct.Name)
	builder.WriteString(fmt.Sprintf("## %s\n", groupTitle))

	for i, field := range configStruct.Fields {
		g.generateField(builder, field)

		// 不是最后一个字段时，不添加空行
		if i < len(configStruct.Fields)-1 {
			// 字段之间不添加空行
		} else {
			// 最后一个字段后添加一个空行（与下一组分隔）
			builder.WriteString("\n")
		}
	}
}

func (g *Generator) generateField(builder *strings.Builder, field types.ConfigField) {
	comment := g.cleanComment(field.Comment)
	varPrefix := "# "

	if !field.IsOptional {
		varPrefix = ""
	}

	value := g.formatValue(field)

	if comment != "" {
		builder.WriteString(fmt.Sprintf("# %s\n", comment))
	}

	builder.WriteString(fmt.Sprintf("%s%s=%s\n", varPrefix, field.EnvVar, value))
}

func (g *Generator) formatGroupTitle(groupName string) string {
	title := strings.ReplaceAll(groupName, "Config", "")
	title = strings.ReplaceAll(title, "_", " ")
	title = strings.TrimSpace(title)

	if title == "" {
		return "Configuration"
	}

	return title + " Configuration"
}

func (g *Generator) cleanComment(comment string) string {
	comment = strings.TrimSpace(comment)
	comment = strings.TrimPrefix(comment, "//")
	comment = strings.TrimSpace(comment)
	comment = strings.TrimPrefix(comment, "*")
	comment = strings.TrimSpace(comment)

	return comment
}

func (g *Generator) formatValue(field types.ConfigField) string {
	if field.DefaultValue != "" {
		return field.DefaultValue
	}

	switch field.Type {
	case "string", "*string":
		return ""
	case "int", "*int", "int64", "*int64":
		return "0"
	case "bool", "*bool":
		return "false"
	case "float64", "*float64":
		return "0.0"
	default:
		return ""
	}
}
