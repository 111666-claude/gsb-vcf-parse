package vcard

import "fmt"

// Scenario 是一个固定场景。
type Scenario struct {
	Name string
	Run  func() map[string]any
}

func card(rows ...string) string {
	out := "BEGIN:VCARD\n"
	for _, r := range rows {
		out += r + "\n"
	}
	return out + "END:VCARD\n"
}

// Samples 返回全部场景。
func Samples() []Scenario {
	return []Scenario{
		{Name: "fold", Run: func() map[string]any {
			book := Parse(card("FN:John ", " Smith"), &Counter{})
			return map[string]any{"fn": book.Cards[0].Get("FN")}
		}},
		{Name: "esc", Run: func() map[string]any {
			book := Parse(card(`NOTE:a\,b`), &Counter{})
			return map[string]any{"len": len(book.Cards[0].Get("NOTE"))}
		}},
		{Name: "param", Run: func() map[string]any {
			book := Parse(card("TEL;TYPE=WORK:+123"), &Counter{})
			return map[string]any{"name": book.Cards[0].Props[0].Name}
		}},
		{Name: "dup", Run: func() map[string]any {
			book := Parse(card("FN:John")+card("FN:John"), &Counter{})
			return map[string]any{"cards": len(book.Cards)}
		}},
	}
}

// Work 跑规模线场景。
func Work(n int) map[string]any {
	text := ""
	for i := 0; i < n; i++ {
		text += card(fmt.Sprintf("FN:p%04d", i))
	}
	book := Parse(text, &Counter{})
	c := &Counter{}
	for i := 0; i < n; i++ {
		book.Find(fmt.Sprintf("p%04d", i), c)
	}
	return map[string]any{"cards": n, "scanned": c.Scanned}
}
