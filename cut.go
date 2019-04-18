package gocommon

import "strings"

// cut 'line', output 'pos' position when splitting by character 'split'
// works like awk -F'split' '{print pos}', i.e. if there are multiple consecutive split characters, they are treated as one
func cut(line string, pos int, split string) string {
	p := 0
	for _, v := range strings.Split(line, split) {
		if v != "" {
			p = p + 1
		}
		if p == pos {
			return v
		}
	}
	return ""
}
