package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"envmon/internal/parser"
)

type Line struct {
	Original   string
	Rule       *parser.FlagRule
	IsFlagLine bool
	Content    string
}

func ParseLine(line string) Line {
	l := Line{Original: line}

	if match := parser.FlagLineRegex.FindStringSubmatch(line); match != nil {
		rule := parser.ParseFlagRule(match[1])
		l.Rule = &rule
		l.IsFlagLine = true
		return l
	}

	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "#") {
		l.Content = strings.TrimPrefix(trimmed, "#")
	} else {
		l.Content = line
	}

	return l
}

func ReadLines() ([]Line, error) {
	file, err := os.Open(".env")
	if err != nil {
		return nil, fmt.Errorf("could not open .env file: %w", err)
	}
	defer file.Close()

	var lines []Line
	var currentRule *parser.FlagRule
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parsed := ParseLine(line)

		if parsed.IsFlagLine {
			currentRule = parsed.Rule
			lines = append(lines, parsed)
		} else {
			if strings.TrimSpace(line) == "" {
				currentRule = nil
				lines = append(lines, parsed)
			} else {
				parsed.Rule = currentRule
				lines = append(lines, parsed)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .env file: %w", err)
	}

	return lines, nil
}

func Process(deploymentFlag string) (string, error) {
	lines, err := ReadLines()
	if err != nil {
		return "", err
	}

	var output strings.Builder
	for _, l := range lines {
		if l.IsFlagLine {
			output.WriteString(l.Original)
			output.WriteString("\n")
			continue
		}

		trimmed := strings.TrimSpace(l.Original)
		if trimmed == "" {
			output.WriteString(l.Original)
			output.WriteString("\n")
			continue
		}

		if l.Rule == nil || l.Rule.Invalid {
			output.WriteString(l.Original)
			output.WriteString("\n")
			continue
		}

		isActive := l.Rule.IsActive(deploymentFlag)
		isCommented := strings.HasPrefix(trimmed, "#")

		if isActive && isCommented {
			// Uncomment: remove leading #
			uncommented := strings.TrimPrefix(trimmed, "#")
			output.WriteString(uncommented)
			output.WriteString("\n")
		} else if !isActive && !isCommented {
			// Comment: add leading #
			output.WriteString("#")
			output.WriteString(l.Original)
			output.WriteString("\n")
		} else {
			output.WriteString(l.Original)
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}

func Write(content string) error {
	return os.WriteFile(".env", []byte(content), 0644)
}

func GetConfigs() ([]string, error) {
	file, err := os.Open(".env")
	if err != nil {
		return nil, fmt.Errorf("could not open .env file: %w", err)
	}
	defer file.Close()

	configs := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if match := parser.FlagLineRegex.FindStringSubmatch(line); match != nil {
			rule := parser.ParseFlagRule(match[1])
			if !rule.Invalid {
				for _, f := range rule.Flags {
					configs[f] = true
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .env file: %w", err)
	}

	result := make([]string, 0, len(configs))
	for cfg := range configs {
		result = append(result, cfg)
	}
	return result, nil
}

func GetCurrentDeployment() (string, error) {
	file, err := os.Open(".env")
	if err != nil {
		return "", fmt.Errorf("could not open .env file: %w", err)
	}
	defer file.Close()

	flagActivity := make(map[string]int)
	allFlags := make(map[string]bool)

	var currentRule *parser.FlagRule
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if match := parser.FlagLineRegex.FindStringSubmatch(line); match != nil {
			rule := parser.ParseFlagRule(match[1])
			currentRule = &rule
			if !rule.Invalid {
				for _, f := range rule.Flags {
					allFlags[f] = true
				}
			}
			continue
		}

		if trimmed == "" {
			currentRule = nil
			continue
		}

		if currentRule != nil && !currentRule.Invalid {
			isCommented := strings.HasPrefix(trimmed, "#")
			if !isCommented && !currentRule.Inverse {
				for _, f := range currentRule.Flags {
					flagActivity[f]++
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading .env file: %w", err)
	}

	if len(allFlags) == 0 {
		return "", nil
	}

	var bestFlag string
	bestCount := 0
	for flag, count := range flagActivity {
		if count > bestCount {
			bestCount = count
			bestFlag = flag
		}
	}

	return bestFlag, nil
}
