package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"unicode"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		fmt.Printf("Erro ao inicializar termui: %v", err)
		return
	}
	defer ui.Close()

	menu := widgets.NewList()
	menu.Title = "Escolha uma operação"
	menu.Rows = []string{
		"1. Somar",
		"2. Subtrair",
		"3. Multiplicar",
		"4. Dividir",
		"5. Sair",
	}
	menu.SetRect(0, 0, 25, 7)
	menu.TextStyle = ui.NewStyle(ui.ColorWhite)
	menu.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGreen)

	input := widgets.NewParagraph()
	input.Title = "Digite os números separados por espaço:"
	input.SetRect(0, 7, 50, 10)
	input.TextStyle = ui.NewStyle(ui.ColorGreen)

	result := widgets.NewParagraph()
	result.Title = "Resultado"
	result.SetRect(0, 10, 50, 13)

	ui.Render(menu, input, result)

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			menu.ScrollDown()
		case "k", "<Up>":
			menu.ScrollUp()
		case "<Enter>":
			choice := menu.Rows[menu.SelectedRow]
			if choice == "5. Sair" {
				fmt.Println("Saindo...")
				return
			}

			operation := ""
			switch choice {
			case "1. Somar":
				operation = "somar"
			case "2. Subtrair":
				operation = "subtrair"
			case "3. Multiplicar":
				operation = "multiplicar"
			case "4. Dividir":
				operation = "dividir"
			}

			input.Text = ""
			result.Text = ""
			menu.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
			ui.Render(input, menu)

			numbers := ""
			for {
				e := <-uiEvents
				if e.ID == "<Enter>" {
					input.TextStyle = ui.NewStyle(ui.ColorWhite)
					ui.Render(input)
					break
				}
				if e.Type == ui.KeyboardEvent {
					if e.ID == "<Space>" && !strings.HasSuffix(input.Text, " ") {
						input.Text += " "
					} else if e.ID == "<Backspace>" {
						if len(input.Text) > 0 {
							input.Text = input.Text[:len(input.Text)-1]
						}
					} else if len(e.ID) == 1 && (unicode.IsDigit(rune(e.ID[0])) || e.ID[0] == ' ') {
						input.Text += e.ID
					}
					ui.Render(input)
				}
			}

			input.TextStyle = ui.NewStyle(ui.ColorWhite)
			numbers = strings.TrimSpace(input.Text)
			res, err := sendRequest(operation, numbers)
			if err != nil {
				result.TextStyle = ui.NewStyle(ui.ColorRed)
				result.Text = fmt.Sprintf("Erro ao enviar requisição: %v", err)
			} else {
				result.TextStyle = ui.NewStyle(ui.ColorGreen)
				result.Text = res
			}
			ui.Render(result)
		}
		menu.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorGreen) // Verde enquanto selecionado
		result.TextStyle = ui.NewStyle(ui.ColorWhite)
		ui.Render(menu)
	}
}

func sendRequest(operation, numbers string) (string, error) {
	conn, err := net.Dial("tcp", "18.208.231.110:15000")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	request := fmt.Sprintf("%s %s\n", operation, numbers)
	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", err
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}
