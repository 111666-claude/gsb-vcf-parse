// Package vcard 解析 vCard 文本。
package vcard

import "strings"

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

// unescape 还原值里的转义序列。
func unescape(s string) string {
	if !strings.ContainsRune(s, '\\') {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n', 'N':
				b.WriteByte('\n')
				i++
				continue
			case ',', ';', '\\':
				b.WriteByte(s[i+1])
				i++
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// parseProp 把一行拆成属性名、参数和值。
func parseProp(line string) (Property, bool) {
	head, value, ok := strings.Cut(line, ":")
	if !ok {
		return Property{}, false
	}
	parts := strings.Split(head, ";")
	p := Property{Name: parts[0], Value: unescape(value)}
	for _, kv := range parts[1:] {
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		if p.Params == nil {
			p.Params = map[string]string{}
		}
		p.Params[k] = v
	}
	return p, true
}

// Parse 解析 vCard 文本。
func Parse(text string, c *Counter) *Book {
	book := &Book{index: map[string]int{}}
	inCard := false
	var cur []Property
	endCard := func() {
		if !inCard {
			return
		}
		inCard = false
		card := Card{Props: cur}
		cur = nil
		fn := card.Get("FN")
		if fn != "" {
			if _, dup := book.index[fn]; dup {
				return
			}
			book.index[fn] = len(book.Cards)
		}
		book.Cards = append(book.Cards, card)
	}
	handle := func(line string) {
		if strings.TrimSpace(line) == "" {
			return
		}
		if strings.EqualFold(line, "BEGIN:VCARD") {
			cur = nil
			inCard = true
			return
		}
		if strings.EqualFold(line, "END:VCARD") {
			endCard()
			return
		}
		if !inCard {
			return
		}
		if p, ok := parseProp(line); ok {
			cur = append(cur, p)
		}
	}
	var logical string
	hasLogical := false
	flush := func() {
		if hasLogical {
			handle(logical)
			logical = ""
			hasLogical = false
		}
	}
	for _, l := range strings.Split(text, "\n") {
		c.Scanned++
		line := strings.TrimRight(l, "\r")
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			logical += strings.TrimLeft(line, " \t")
			continue
		}
		flush()
		logical = line
		hasLogical = true
	}
	flush()
	endCard()
	return book
}

// Find 按 FN 找卡片。
func (b *Book) Find(fn string, c *Counter) *Card {
	c.Scanned++
	if i, ok := b.index[fn]; ok {
		return &b.Cards[i]
	}
	return nil
}
