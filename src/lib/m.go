package main

import (
	"fmt"
	"os"
)

type People interface {
	Speak(string) string
}

type Stduent struct{}

func (stu *Stduent) Speak(think string) (talk string) {
	if think == "love" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return

}

func main() {
	//var peo People = Stduent{}
	//think := "love"
	//fmt.Println(peo.Speak(think))
	//s := Stduent{}
	//_s := &s
	//s.Speak("1111")
	//_s.Speak("111111")
	s := make([]int, 5, 10)
	s1 := append(s, []int{1, 2}...)
	fmt.Printf("S:  %v  len %v cap %v", s, len(s), cap(s))
	fmt.Printf("S1: %v  len %v cap %v", s1, len(s1), cap(s1))
	os, os.NewFile()

}
