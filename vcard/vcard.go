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
}

// Parse 解析 vCard 文本。
func Parse(text string, c *Counter) *Book {
	book := &Book{}
	inCard := false
	var cur []Property
	for _, l := range strings.Split(text, "\n") {
		c.Scanned++
		line := strings.TrimSpace(strings.TrimRight(l, "\r"))
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
				book.Cards = append(book.Cards, Card{Props: cur})
			}
			inCard = false
			continue
		}
		if !inCard {
			continue
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		cur = append(cur, Property{Name: name, Value: value})
	}
	return book
}

// Find 按 FN 找卡片。
func (b *Book) Find(fn string, c *Counter) *Card {
	for i := range b.Cards {
		c.Scanned++
		if b.Cards[i].Get("FN") == fn {
			return &b.Cards[i]
		}
	}
	return nil
}
