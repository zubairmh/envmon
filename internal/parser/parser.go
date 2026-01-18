package parser

import (
	"regexp"
	"strings"
)

var (
	FlagLineRegex = regexp.MustCompile(`^#\s*\[([^\]]+)\]`)

	inverseMultiRegex = regexp.MustCompile(`^~\(([^)]+)\)$`)

	inverseSimpleRegex = regexp.MustCompile(`^~([^/]+)$`)
)

type FlagRule struct {
	Flags   []string
	Inverse bool
	Invalid bool
}

func ParseFlagRule(raw string) FlagRule {
	raw = strings.TrimSpace(raw)
	
	if match := inverseMultiRegex.FindStringSubmatch(raw); match != nil {
		flags := strings.Split(match[1], "/")
		for i := range flags {
			flags[i] = strings.TrimSpace(flags[i])
		}
		return FlagRule{
			Flags:   flags,
			Inverse: true,
		}
	}

	if match := inverseSimpleRegex.FindStringSubmatch(raw); match != nil {
		return FlagRule{
			Flags:   []string{strings.TrimSpace(match[1])},
			Inverse: true,
		}
	}
	if strings.HasPrefix(raw, "~") && strings.Contains(raw, "/") {
		return FlagRule{Invalid: true}
	}
	flags := strings.Split(raw, "/")
	for i := range flags {
		flags[i] = strings.TrimSpace(flags[i])
	}
	return FlagRule{
		Flags:   flags,
		Inverse: false,
	}
}


func (r FlagRule) IsActive(deploymentFlag string) bool {
	if r.Invalid {
		return false
	}

	matched := false
	for _, f := range r.Flags {
		if f == deploymentFlag {
			matched = true
			break
		}
	}

	if r.Inverse {
		return !matched
	}
	return matched
}