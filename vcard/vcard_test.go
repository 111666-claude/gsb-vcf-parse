package vcard

import "testing"

func TestSingleCardCount(t *testing.T) {
	book := Parse(card("FN:John"), &Counter{})
	if len(book.Cards) != 1 {
		t.Fatalf("应有一张卡片：%d", len(book.Cards))
	}
}

func TestFnValue(t *testing.T) {
	book := Parse(card("FN:John"), &Counter{})
	if got := book.Cards[0].Get("FN"); got != "John" {
		t.Fatalf("FN 应是 John：%q", got)
	}
}

func TestPropValue(t *testing.T) {
	book := Parse(card("TEL:+123"), &Counter{})
	if got := book.Cards[0].Get("TEL"); got != "+123" {
		t.Fatalf("TEL 应是 +123：%q", got)
	}
}

func TestTwoCards(t *testing.T) {
	book := Parse(card("FN:A")+card("FN:B"), &Counter{})
	if len(book.Cards) != 2 {
		t.Fatalf("应有两张卡片：%d", len(book.Cards))
	}
}

func TestFindMissing(t *testing.T) {
	book := Parse(card("FN:John"), &Counter{})
	if book.Find("nope", &Counter{}) != nil {
		t.Fatal("缺失姓名应返回 nil")
	}
}

func TestCounterCounts(t *testing.T) {
	c := &Counter{}
	Parse(card("FN:John"), c)
	if c.Scanned == 0 {
		t.Fatal("应统计比较次数")
	}
}
