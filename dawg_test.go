package dawg

import (
	"fmt"
	"testing"
)

func TestSimple(t *testing.T) {
	builder := NewDawgBuilder()
	fmt.Println("Inserting apple:", builder.InsertString("apple"))
	fmt.Println("Inserting appliance:", builder.InsertString("appliance"))
	fmt.Println("Inserting applied:", builder.InsertString("applied"))
	fmt.Println("Inserting apply:", builder.InsertString("apply"))
	fmt.Println("Inserting banana:", builder.InsertString("banana"))
	fmt.Println("Inserting changed:", builder.InsertString("changed"))
	fmt.Println("Inserting cherry:", builder.InsertString("cherry"))
	fmt.Println("Inserting durian:", builder.InsertString("durian"))
	fmt.Println("Inserting mandarin:", builder.InsertString("mandarin"))
	fmt.Println("Inserting murdered:", builder.InsertString("murdered"))
	fmt.Println("Inserting office:", builder.InsertString("office"))

	dawg := NewDawg()
	builder.Finish(dawg)

	dict := dawg.Build()
	fmt.Println("Contains apple:", dict.ContainsString("apple"))
	fmt.Println("Contains cherry:", dict.ContainsString("cherry"))
	fmt.Println("Contains durian:", dict.ContainsString("durian"))
	fmt.Println("Contains green:", dict.ContainsString("green"))
	fmt.Println("Contains mandarin:", dict.ContainsString("mandarin"))
	fmt.Println("Contains applied:", dict.ContainsString("applied"))
	fmt.Println("Contains murdered:", dict.ContainsString("murdered"))
	fmt.Println("Contains changed:", dict.ContainsString("changed"))
	fmt.Println("Contains appliance:", dict.ContainsString("appliance"))
	fmt.Println("Contains change:", dict.ContainsString("change"))
}

func TestZeros(t *testing.T) {
	builder := NewDawgBuilder()
	fmt.Println("Inserting 00:", builder.InsertKeyValue([]ucharType{0, 0}, 2, 0))
	fmt.Println("Inserting 00000:", builder.InsertKeyValue([]ucharType{0, 0, 0, 0, 0}, 2, 0))

	dawg := NewDawg()
	builder.Finish(dawg)

	dict := dawg.Build()
	fmt.Println("Contains 0:", dict.ContainsString("\x00"))
	fmt.Println("Contains 00:", dict.ContainsString("\x00\x00"))
	fmt.Println("Contains 000:", dict.ContainsString("\x00\x00\x00"))
	fmt.Println("Contains 00000:", dict.ContainsString("\x00\x00\x00\x00\x00"))
	fmt.Println("Contains 000000:", dict.ContainsString("\x00\x00\x00\x00\x00\x00"))
}
