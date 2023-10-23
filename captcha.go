package main

import (
	"fmt"
	"math/rand"
)

var font = [][]string{
	{"  ", "_ ", "_ ", "   ", " _ ", " _ ", "__", " _ ", " _ "},
	{"/|", " )", "_)", "|_|", "|_ ", "|_ ", " /", "(_)", "(_|"},
	{" |", "/_", "_)", "  |", " _)", "|_)", "/ ", "(_)", " _|"},
}

type captcha struct {
	value []int
}

func NewCaptcha(len int) captcha {
	c := captcha{}
	c.value = rand.Perm(len)
	return c
}

func (c captcha) IsValid(v string) bool {
	s := ""
	for _, v := range c.value {
		s += fmt.Sprintf("%d", v+1)
	}
	return v == s
}

func (c captcha) View() string {
	captcha := ""
	for y := 0; y < len(font); y++ {
		captcha += fmt.Sprintf("%s %s %s %s\n", font[y][c.value[0]], font[y][c.value[1]], font[y][c.value[2]], font[y][c.value[3]])
	}
	return captcha
}
