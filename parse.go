package main

import (
	"strings"
	"unicode"
	"bytes"
)

const (
	Bold		= 1 << iota
	Italic		= 1 << iota
	Underline	= 1 << iota
	Citation	= 1 << iota
	Escaped		= 1 << iota
)

type rule struct {
	rid	int
	s,e	string
}

var rules = map[rune]*rule {
	'*':	&rule{Bold, "<b>", "</b>"},
	'/':	&rule{Italic, "<i>", "</i>"},
	'_':	&rule{Underline, "<u>", "</u>"},
}
	
func parse(s string) string {
	p := s + "\n"	/* Add a newline to the end to make everything work better. */
	var r bytes.Buffer
	cite := []rune{}
	format := 0
	pchar := rune(0)
	
	for _, c := range p {
		if format & Citation == Citation {
			if unicode.IsDigit(c) {
				cite = append(cite, c)
			} else {
				format &= ^Citation
				if len(cite) > 0 {
					r.WriteString(`<a href="#` +
						string(cite) +
						`">@` + string(cite) + `</a>` +
						string(c))
					cite = []rune{}
				} else {
					r.WriteString("@" + string(c))
				}
			}
		} else {
			switch c {
			case '*', '/', '_':
				rule,_ := rules[c] 
				if pchar == c {
					if format & Escaped == Escaped {
						r.WriteString(string(c) + string(c))
						format &= ^Escaped
					} else if format & rule.rid == rule.rid {
						r.WriteString(rule.e)
						format &= ^rule.rid
					} else {
						r.WriteString(rule.s)
						format |= rule.rid
					}
					pchar = rune(0)
				} else {
					if pchar != rune(0) {
						r.WriteString(string(pchar))
					}
					pchar = c
				}
			case '@':
				if format & Escaped == Escaped {
					r.WriteString("@")
					format &= ^Escaped
				}
				pchar = rune(0)
				format |= Citation
			case '\\':
				if format & Escaped == Escaped {
					r.WriteString("\\")
					format &= ^Escaped
				} else {
					format |= Escaped
				}
				pchar = rune(0)
			default:
				if format & Escaped == Escaped {
					r.WriteString("\\")
					format &= ^Escaped
				}
				if pchar != rune(0) {
					r.WriteString(string(pchar))
				}
				pchar = rune(0)
				r.WriteString(string(c))
			}
		}
	}
	return strings.Replace(strings.TrimSpace(r.String()), "\n", "<br />", -1)
}

/*
func parseTrip(s string) string {
	ind := Index(s, "#")
	if ind == -1 {
		return s
	}
}
 */