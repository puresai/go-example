package main

import (
	"demo5-OOP/oop"
)

func main() {
	gf := oop.NewOne()
	gf.SetName("Lily")
	gf.SetHeight(170)
	gf.SetWeight(50)
	gf.SetGreeting("hello,").SetContent("thanks for your great creavity!").SetAge(24)
	gf.Show()
}
