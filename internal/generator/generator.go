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

	groupedFields := g.groupFieldsByStruct()

	groupNames := make([]string, 0, len(groupedFields))
	for name := range groupedFields {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		fields := groupedFields[groupName]

		if len(fields) > 0 {
			g.generateGroup(&builder, groupName, fields)
		}
	}

	return builder.String()
}

func (g *Generator) groupFieldsByStruct() map[string][]types.ConfigField {
	grouped := make(map[string][]types.ConfigField)

	for _, configStruct := range g.structs {
		for _, field := range configStruct.Fields {
			grouped[field.Group] = append(grouped[field.Group], field)
		}
	}

	return grouped
}

func (g *Generator) generateGroup(builder *strings.Builder, groupName string, fields []types.ConfigField) {
	groupTitle := g.formatGroupTitle(groupName)
	builder.WriteString(fmt.Sprintf("\n## %s\n", groupTitle))

	for _, field := range fields {
		g.generateField(builder, field)
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
	builder.WriteString("\n")
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
