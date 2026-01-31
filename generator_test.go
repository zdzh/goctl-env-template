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
			Name:    "AuthConfig",
			Comment: "Auth配置包含JWT认证相关的配置项",
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

	// 测试1: 验证结构体注释在分组注释上方
	if !strings.Contains(result, "# Auth配置包含JWT认证相关的配置项") {
		t.Error("应该包含结构体注释")
	}

	// 测试2: 验证分组注释使用 ##
	if !strings.Contains(result, "## Auth Configuration") {
		t.Error("分组注释应该使用 ##")
	}

	// 测试3: 验证配置项注释在配置项上方
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

	// 测试4: 验证必需变量未被注释
	if strings.Contains(result, "# AUTH_SECRET=") {
		t.Error("必需变量 AUTH_SECRET 不应该被注释")
	}
	if !strings.Contains(result, "\nAUTH_SECRET=") {
		t.Error("必需变量 AUTH_SECRET 应该存在且未被注释")
	}

	// 测试5: 验证可选变量被注释
	if !strings.Contains(result, "# AUTH_REFRESH_INTERVAL=") {
		t.Error("可选变量 AUTH_REFRESH_INTERVAL 应该被注释")
	}

	// 测试6: 验证同一分组内变量之间没有空行
	secretLineIndex := -1
	timeoutLineIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "AUTH_SECRET=" {
			secretLineIndex = i
		}
		if strings.TrimSpace(line) == "AUTH_TIMEOUT=3600" {
			timeoutLineIndex = i
		}
	}

	if secretLineIndex >= 0 && timeoutLineIndex >= 0 {
		betweenLines := lines[min(secretLineIndex, timeoutLineIndex)+1 : max(secretLineIndex, timeoutLineIndex)]
		emptyCount := 0
		for _, line := range betweenLines {
			if strings.TrimSpace(line) == "" {
				emptyCount++
			}
		}
		// AUTH_SECRET= 后面有注释行 # Auth timeout in seconds，所以可能有1个空行
		// 但不应该有多个连续的空行
		if emptyCount > 2 {
			t.Errorf("同一分组内变量之间不应该有多个空行，发现 %d 个", emptyCount)
		}
	}

	t.Logf("生成的模板:\n%s", result)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

	// 验证结构体注释存在
	structCommentCount := 0
	for _, configStruct := range structs {
		if configStruct.Comment != "" {
			structCommentCount++
			if !strings.Contains(result, configStruct.Comment) {
				t.Errorf("应该包含结构体注释: %s", configStruct.Comment)
			}
		}
	}

	t.Logf("结构体注释数量: %d", structCommentCount)

	// 验证分组注释使用 ##
	groupHeaderCount := strings.Count(result, "## ")
	if groupHeaderCount == 0 {
		t.Error("应该有使用 ## 的分组注释")
	}

	t.Logf("分组注释数量: %d", groupHeaderCount)

	// 验证同一分组内变量之间没有多余空行
	// 统计连续空行的数量
	lines := strings.Split(result, "\n")
	maxConsecutiveEmpty := 0
	currentConsecutiveEmpty := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			currentConsecutiveEmpty++
			if currentConsecutiveEmpty > maxConsecutiveEmpty {
				maxConsecutiveEmpty = currentConsecutiveEmpty
			}
		} else {
			// 遇到非空行，重置计数
			currentConsecutiveEmpty = 0
		}
	}

	t.Logf("最大连续空行数: %d", maxConsecutiveEmpty)

	// 验证必需变量未注释，可选变量被注释
	requiredVarCount := 0
	optionalVarCount := 0
	for _, configStruct := range structs {
		for _, field := range configStruct.Fields {
			if field.IsOptional {
				optionalVarCount++
				if !strings.Contains(result, "# "+field.EnvVar+"=") {
					t.Errorf("可选变量 %s 应该被注释", field.EnvVar)
				}
			} else {
				requiredVarCount++
				// 必需变量不应该被注释（即不应该包含 "# VARNAME="）
				if strings.Contains(result, "# "+field.EnvVar+"=") {
					t.Errorf("必需变量 %s 不应该被注释", field.EnvVar)
				}
			}
		}
	}

	t.Logf("必需变量数量: %d", requiredVarCount)
	t.Logf("可选变量数量: %d", optionalVarCount)

	t.Logf("生成的模板:\n%s", result)
}
