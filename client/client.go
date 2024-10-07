package main

import (
	"bufio"
	"fmt"
	"net"

	"strings"

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
	menu.TextStyle = ui.NewStyle(ui.ColorYellow)
	menu.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorYellow)

	input := widgets.NewParagraph()
	input.Title = "Digite os números separados por espaço:"
	input.SetRect(0, 7, 50, 10)

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
			ui.Render(input)

			numbers := ""
			for {
				uiEvents := ui.PollEvents()
				e := <-uiEvents
				if e.ID == "<Enter>" {
					break
				}
				if e.Type == ui.KeyboardEvent {
					input.Text += e.ID
					ui.Render(input)
				}
			}

			numbers = strings.TrimSpace(input.Text)
			res, err := sendRequest(operation, numbers)
			if err != nil {
				result.Text = fmt.Sprintf("Erro ao enviar requisição: %v", err)
			} else {
				result.Text = res
			}
			ui.Render(result)
		}

		ui.Render(menu, input, result)
	}
}

func sendRequest(operation, numbers string) (string, error) {
	conn, err := net.Dial("tcp", "3.80.253.101:15000")
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
