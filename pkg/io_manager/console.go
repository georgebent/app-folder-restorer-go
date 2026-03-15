package io_manager

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func Write(description string) {
	fmt.Println(description)
}

func Read(question string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	text, err := reader.ReadString('\n')
	if err != nil && len(text) == 0 {
		return ""
	}
	text = strings.TrimSpace(text)

	fmt.Println("Ви ввели:", text)

	return text
}

func Ask(question string, options map[string]string) string {
	if len(options) == 0 {
		os.Exit(0)
	}

	if !supportsTTYMenu() {
		return askPlain(question, options, Read)
	}

	menu := NewMenu(question)
	for _, key := range optionKeys(options) {
		menu.AddItem(options[key], key)
	}

	choice := menu.Display()
	if strings.ToLower(choice) == "exit" {
		fmt.Println("Вихід з програми...")
		os.Exit(0)
	}

	if _, ok := options[choice]; !ok {
		return Ask("Вибір не вірний. Будь ласка, спробуйте ще раз.", options)
	}

	fmt.Printf("Choice: %s\n", options[choice])

	return choice
}

func askPlain(question string, options map[string]string, read func(string) string) string {
	for {
		fmt.Printf("%s:\n", question)
		for _, key := range optionKeys(options) {
			fmt.Printf("%s. %s\n", key, options[key])
		}

		choice := read("Choice: ")
		if strings.ToLower(choice) == "exit" {
			fmt.Println("Вихід з програми...")
			os.Exit(0)
		}

		if _, ok := options[choice]; ok {
			fmt.Printf("Choice: %s\n", options[choice])
			return choice
		}

		question = "Вибір не вірний. Будь ласка, спробуйте ще раз."
	}
}

func optionKeys(options map[string]string) []string {
	keys := make([]string, 0, len(options))
	for key := range options {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
