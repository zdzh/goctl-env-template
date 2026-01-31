package main

import (
	"os"
	"strings"
	"testing"

	"github.com/zdzh/goctl-env-template/internal/generator"
	"github.com/zdzh/goctl-env-template/internal/parser"
	"github.com/zdzh/goctl-env-template/internal/types"
)

func TestCommentFormat(t *testing.T) {
	// 准备测试数据
	testStructs := []types.ConfigStruct{
		{
			Name: "AuthConfig",
			Fields: []types.ConfigField{
				{
					Name:       "Secret",
					EnvVar:     "AUTH_SECRET",
					Type:       "string",
					IsOptional: false,
					Comment:    "JWT secret key",
					Group:      "AuthConfig",
				},
				{
					Name:         "Timeout",
					EnvVar:       "AUTH_TIMEOUT",
					Type:         "int",
					IsOptional:   false,
					DefaultValue: "3600",
					Comment:      "Auth timeout in seconds",
					Group:        "AuthConfig",
				},
				{
					Name:       "RefreshInterval",
					EnvVar:     "AUTH_REFRESH_INTERVAL",
					Type:       "int",
					IsOptional: true,
					Comment:    "Token refresh interval (optional)",
					Group:      "AuthConfig",
				},
			},
		},
	}

	// 生成模板
	gen := generator.NewGenerator(testStructs)
	result := gen.Generate()

	// 测试1: 验证分组注释使用 ##
	if !strings.Contains(result, "## Auth Configuration") {
		t.Error("分组注释应该使用 ##")
	}

	// 测试2: 验证配置项注释在配置项上方
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "##") {
			// 找到注释行，下一行应该是配置项
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine != "" && !strings.HasPrefix(nextLine, "#") {
					// 验证这是配置项
					if strings.Contains(nextLine, "=") {
						// 确认配置项下方没有相同的注释
						if i+2 < len(lines) {
							afterLine := strings.TrimSpace(lines[i+2])
							if strings.HasPrefix(afterLine, "# ") && afterLine == line {
								t.Errorf("配置项说明应该在上方，不应该在下方: %s", line)
							}
						}
					}
				}
			}
		}
	}

	// 测试3: 验证必需变量未被注释
	if strings.Contains(result, "# AUTH_SECRET=") {
		t.Error("必需变量 AUTH_SECRET 不应该被注释")
	}
	if !strings.Contains(result, "\nAUTH_SECRET=") {
		t.Error("必需变量 AUTH_SECRET 应该存在且未被注释")
	}

	// 测试4: 验证可选变量被注释
	if !strings.Contains(result, "# AUTH_REFRESH_INTERVAL=") {
		t.Error("可选变量 AUTH_REFRESH_INTERVAL 应该被注释")
	}

	// 测试5: 验证注释在配置项上方
	secretCommentIndex := strings.Index(result, "# JWT secret key")
	if secretCommentIndex == -1 {
		t.Error("找不到 JWT secret key 注释")
	} else {
		secretIndex := strings.Index(result, "AUTH_SECRET=")
		if secretIndex < secretCommentIndex {
			t.Error("配置项注释应该在配置项上方")
		}
	}

	t.Logf("生成的模板:\n%s", result)
}

func TestRealConfigFile(t *testing.T) {
	// 测试真实的配置文件
	configPath := "config/config.go"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skip("配置文件不存在，跳过测试")
	}

	p, err := parser.NewParser(configPath)
	if err != nil {
		t.Fatalf("创建解析器失败: %v", err)
	}

	structs, err := p.Parse()
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	if len(structs) == 0 {
		t.Skip("没有找到配置结构体")
	}

	gen := generator.NewGenerator(structs)
	result := gen.Generate()

	// 验证格式
	lines := strings.Split(result, "\n")
	commentAboveCount := 0
	commentBelowCount := 0

	for i := 0; i < len(lines)-1; i++ {
		line := strings.TrimSpace(lines[i])
		nextLine := strings.TrimSpace(lines[i+1])

		// 检查注释在配置项上方
		if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "##") {
			if strings.Contains(nextLine, "=") && !strings.HasPrefix(nextLine, "#") {
				commentAboveCount++
			}
		}

		// 检查注释在配置项下方（不应该存在）
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			if strings.HasPrefix(nextLine, "# ") && !strings.HasPrefix(nextLine, "##") {
				commentBelowCount++
			}
		}
	}

	t.Logf("注释在上方: %d", commentAboveCount)
	t.Logf("注释在下方: %d", commentBelowCount)

	if commentBelowCount > 0 {
		t.Error("不应该有注释在配置项下方")
	}

	// 验证分组注释使用 ##
	groupHeaderCount := strings.Count(result, "## ")
	if groupHeaderCount == 0 {
		t.Error("应该有使用 ## 的分组注释")
	}

	t.Logf("生成的模板:\n%s", result)
}
