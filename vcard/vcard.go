// Package vcard 解析 vCard 文本。
package vcard

import (
	"strings"
)

// Counter 记录比较次数。
type Counter struct {
	Scanned int
}

// Property 是一条属性。
type Property struct {
	Name   string
	Params map[string]string
	Value  string
}

// Card 是一张卡片。
type Card struct {
	Props []Property
}

// Get 取第一个同名属性的值。
func (c *Card) Get(name string) string {
	for i := range c.Props {
		if c.Props[i].Name == name {
			return c.Props[i].Value
		}
	}
	return ""
}

// Book 是通讯录。
type Book struct {
	Cards []Card
	index map[string]int
}

// Parse 解析 vCard 文本。
func Parse(text string, c *Counter) *Book {
	book := &Book{index: map[string]int{}}
	inCard := false
	var cur []Property
	for _, line := range logicalLines(text) {
		c.Scanned++
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "BEGIN:VCARD") {
			cur = nil
			inCard = true
			continue
		}
		if strings.EqualFold(line, "END:VCARD") {
			if inCard {
				book.add(Card{Props: cur})
			}
			inCard = false
			continue
		}
		if !inCard {
			continue
		}
		if p, ok := parseProperty(line); ok {
			cur = append(cur, p)
		}
	}
	return book
}

// logicalLines 还原折行：以空格或制表符开头的行去掉行首空白后接到上一条。
func logicalLines(text string) []string {
	var out []string
	for _, l := range strings.Split(text, "\n") {
		l = strings.TrimRight(l, "\r")
		if l == "" {
			continue
		}
		if (l[0] == ' ' || l[0] == '\t') && len(out) > 0 {
			out[len(out)-1] += strings.TrimLeft(l, " \t")
			continue
		}
		out = append(out, l)
	}
	return out
}

// parseProperty 拆出属性名、参数和还原转义后的值。
func parseProperty(line string) (Property, bool) {
	head, value, ok := strings.Cut(line, ":")
	if !ok {
		return Property{}, false
	}
	name, paramsPart, _ := strings.Cut(head, ";")
	p := Property{Name: name, Value: unescape(value)}
	if paramsPart != "" {
		p.Params = map[string]string{}
		for _, kv := range strings.Split(paramsPart, ";") {
			k, v, _ := strings.Cut(kv, "=")
			p.Params[k] = v
		}
	}
	return p, true
}

// unescape 把 \, \; \n \\ 还原成对应单字符。
func unescape(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n', 'N':
				b.WriteByte('\n')
			case ',':
				b.WriteByte(',')
			case ';':
				b.WriteByte(';')
			case '\\':
				b.WriteByte('\\')
			default:
				b.WriteByte(s[i])
				b.WriteByte(s[i+1])
			}
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// add 按 FN 去重：同一个 FN 只保留第一次出现的卡片。
func (b *Book) add(card Card) {
	fn := card.Get("FN")
	if fn == "" {
		b.Cards = append(b.Cards, card)
		return
	}
	if _, dup := b.index[fn]; dup {
		return
	}
	b.index[fn] = len(b.Cards)
	b.Cards = append(b.Cards, card)
}

// Find 按 FN 找卡片。
func (b *Book) Find(fn string, c *Counter) *Card {
	c.Scanned++
	if i, ok := b.index[fn]; ok {
		return &b.Cards[i]
	}
	return nil
}
