package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/zdzh/goctl-env-template/internal/types"
)

type Parser struct {
	filePath string
	fset     *token.FileSet
	file     *ast.File
}

func NewParser(filePath string) (*Parser, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", filePath)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &Parser{
		filePath: filePath,
		fset:     fset,
		file:     file,
	}, nil
}

func (p *Parser) Parse() ([]types.ConfigStruct, error) {
	var structs []types.ConfigStruct

	for _, decl := range p.file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			configStruct := types.ConfigStruct{
				Name:   typeSpec.Name.Name,
				Fields: p.extractFields(structType, typeSpec.Name.Name),
			}

			if len(configStruct.Fields) > 0 {
				structs = append(structs, configStruct)
			}
		}
	}

	return structs, nil
}

func (p *Parser) extractFields(structType *ast.StructType, groupName string) []types.ConfigField {
	var fields []types.ConfigField

	if structType.Fields == nil {
		return fields
	}

	for _, field := range structType.Fields.List {
		if field.Names == nil {
			continue
		}

		for _, name := range field.Names {
			configField := types.ConfigField{
				Name:  name.Name,
				Type:  p.getTypeString(field.Type),
				Group: groupName,
			}

			if field.Tag != nil {
				configField = p.parseStructTag(configField, field.Tag.Value)
			}

			if field.Comment != nil {
				configField.Comment = field.Comment.Text()
			} else if field.Doc != nil {
				configField.Comment = field.Doc.Text()
			}

			if configField.EnvVar != "" {
				fields = append(fields, configField)
			}
		}
	}

	return fields
}

func (p *Parser) parseStructTag(field types.ConfigField, tag string) types.ConfigField {
	tag = strings.Trim(tag, "`")

	for _, part := range strings.Split(tag, " ") {
		if strings.Contains(part, "json") && strings.Contains(part, "env=") {
			envVar := p.extractEnvVar(part)
			if envVar != "" {
				field.EnvVar = envVar
			}

			if strings.Contains(part, "optional") {
				field.IsOptional = true
			}

			defaultValue := p.extractDefault(part)
			if defaultValue != "" {
				field.DefaultValue = defaultValue
			}
		}
	}

	return field
}

func (p *Parser) extractEnvVar(tag string) string {
	parts := strings.Split(tag, "env=")
	if len(parts) < 2 {
		return ""
	}

	envPart := parts[1]
	commaIdx := strings.Index(envPart, ",")
	if commaIdx != -1 {
		envPart = envPart[:commaIdx]
	}

	envPart = strings.TrimSpace(envPart)
	envPart = strings.TrimSuffix(envPart, `"`)
	envPart = strings.TrimSuffix(envPart, "'")

	return envPart
}

func (p *Parser) extractDefault(tag string) string {
	parts := strings.Split(tag, "default=")
	if len(parts) < 2 {
		return ""
	}

	defaultPart := parts[1]
	commaIdx := strings.Index(defaultPart, ",")
	if commaIdx != -1 {
		defaultPart = defaultPart[:commaIdx]
	}

	defaultPart = strings.TrimSpace(defaultPart)
	defaultPart = strings.TrimSuffix(defaultPart, `"`)
	defaultPart = strings.TrimSuffix(defaultPart, "'")

	return defaultPart
}

func (p *Parser) getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X, t.Sel)
	case *ast.StarExpr:
		return "*" + p.getTypeString(t.X)
	case *ast.ArrayType:
		return "[]" + p.getTypeString(t.Elt)
	default:
		return ""
	}
}
