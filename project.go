package weiqi

import (
    "strings"
)

type Item struct {
    Name string
    Values []string
}

func newItem(name string, values ...string) Item {
    var i = Item{name, make([]string, 0)}
    for _, v := range values {
        i.Values = append(i.Values, v)
    }
    return i
}

func (i Item) LenValue() int {
    return len(i.Values)
}

func (i *Item) AddValue( value string) {
    i.Values = append(i.Values, value)
}

func (i *Item) ClearValue() {
    i.Values = make([]string, 0)
}

func (i Item) StrValue() string {
    return strings.Join(i.Values, ",")
}

type Project struct {
    Name string
    Text string
    Items []Item
}

type Event struct {
    Name string
    Number int
    Items []Item
    Text string
}