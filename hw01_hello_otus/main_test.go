package main

import "testing"

func TestHello(t *testing.T) {
	input := "Проверка, юнит-hello!"
	expected := "!olleh-тиню ,акреворП"
	actual := reverseString(input)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
