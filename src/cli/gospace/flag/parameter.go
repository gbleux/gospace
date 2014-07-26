package flag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Parameter struct {
	Short       byte
	Long        string
	ValueName   string
	Description string
}

func (p *Parameter) IsFlag() bool {
	return 0 == len(p.ValueName)
}

func (p *Parameter) Matches(value string) bool {
	short := "-" + string(p.Short)
	long := "--" + p.Long

	return strings.HasPrefix(value, short) ||
		strings.HasPrefix(value, long)
}

func (p *Parameter) ParseValue(value string) string {
	stop := false
	trim := func(char rune) bool {
		if stop {
			return false
		} else if '=' == char {
			stop = true
		}

		return true
	}

	return strings.TrimLeftFunc(value, trim)
}

func (p *Parameter) ParseValueOr(value string, env string, fallback string) string {
	result := p.ParseValue(value)

	switch {
	case 0 < len(result):
		return result
	case 0 < len(env):
		return os.Getenv(env)
	case 0 < len(fallback):
		return fallback
	default:
		return ""
	}
}

func (p *Parameter) ParseIntValueOr(value string, fallback int) int {
	result := p.ParseValue(value)

	if 0 < len(result) {
		if num, err := strconv.Atoi(result); nil == err {
			return num
		}
	}

	return fallback
}

func (p *Parameter) Usage() string {
	long := p.Long

	if false == p.IsFlag() {
		long = p.Long + "=" + p.ValueName
	}

	return fmt.Sprintf("\t-%c, --%-15s %s\n",
		p.Short,
		long,
		p.Description)
}

func NewArgParameter(short byte, long string, argument string, description string) *Parameter {
	return &Parameter{short, long, argument, description}
}

func NewFlagParameter(short byte, long string, description string) *Parameter {
	return &Parameter{short, long, "", description}
}
