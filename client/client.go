package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"unicode"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const serverIP = "3.225.60.216:15000"

func main() {
	// Configura o logger para escrever em client.log
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo de log: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	if err := ui.Init(); err != nil {
		log.Fatalf("Erro ao inicializar termui: %v", err)
	}
	defer ui.Close()

	menu, input, result := setupUI()
	ui.Render(menu, input, result)

	uiEvents := ui.PollEvents()

	operations, err := getAvailableOperations()
	if err != nil {
		log.Fatalf("Erro ao obter operações disponíveis: %v", err)
	}

	menu.Rows = operations
	ui.Render(menu)

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
			handleEnter(menu, input, result, uiEvents)
		}
		ui.Render(menu)
	}
}

func setupUI() (*widgets.List, *widgets.Paragraph, *widgets.Paragraph) {
	menu := widgets.NewList()
	menu.Title = "Escolha uma operação"
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

	return menu, input, result
}

func handleEnter(menu *widgets.List, input *widgets.Paragraph, result *widgets.Paragraph, uiEvents <-chan ui.Event) {
	choice := menu.Rows[menu.SelectedRow]
	if choice == "5. Sair" {
		fmt.Println("Saindo...")
		return
	}

	operation := getOperation(choice)
	input.Text = ""
	result.Text = ""
	menu.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	ui.Render(input, menu)

	numbers := getUserInput(input, uiEvents)

	res, err := sendRequest(operation, numbers)
	if err != nil {
		result.TextStyle = ui.NewStyle(ui.ColorRed)
		result.Text = fmt.Sprintf("Erro ao enviar requisição: %v", err)
	} else {
		result.TextStyle = ui.NewStyle(ui.ColorGreen)
		result.Text = res
	}
	ui.Render(result)

	// Aguarda confirmação do usuário antes de reiniciar o fluxo
	confirmRestart(menu, input, result, uiEvents)
}

func getOperation(choice string) string {
	// Extrai o nome da operação a partir da escolha do usuário
	parts := strings.Split(choice, " ")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func getUserInput(input *widgets.Paragraph, uiEvents <-chan ui.Event) string {
	for {
		e := <-uiEvents
		if e.ID == "<Enter>" {
			input.TextStyle = ui.NewStyle(ui.ColorWhite)
			ui.Render(input)
			break
		}
		if e.Type == ui.KeyboardEvent {
			handleKeyboardEvent(input, e)
			ui.Render(input)
		}
	}
	return strings.TrimSpace(input.Text)
}

func handleKeyboardEvent(input *widgets.Paragraph, e ui.Event) {
	if e.ID == "<Space>" && !strings.HasSuffix(input.Text, " ") {
		input.Text += " "
	} else if e.ID == "<Backspace>" {
		if len(input.Text) > 0 {
			input.Text = input.Text[:len(input.Text)-1]
		}
	} else if len(e.ID) == 1 && (unicode.IsDigit(rune(e.ID[0])) || e.ID[0] == '.' || e.ID[0] == '-') {
		input.Text += e.ID
	}
}

func sendRequest(operation, numbers string) (string, error) {
	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		log.Printf("Erro ao conectar ao servidor: %v", err)
		return "", err
	}
	defer conn.Close()

	// Envia a mensagem inicial indicando que deseja realizar uma operação
	_, err = conn.Write([]byte("operation\n"))
	if err != nil {
		log.Printf("Erro ao enviar mensagem inicial: %v", err)
		return "", err
	}

	// Envia a operação e os números para o servidor
	request := fmt.Sprintf("%s %s\n", operation, numbers)
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Printf("Erro ao enviar requisição: %v", err)
		return "", err
	}
	log.Printf("Requisição enviada: %s", request)    // Log para depuração
	log.Printf("Aguardando resposta do servidor...") // Log para depuração

	// Lê a resposta do servidor
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Erro ao ler resposta: %v", err)
		return "", err
	}

	// Remove o prefixo "Resultado: " da resposta
	log.Printf("Resposta do servidor: %s", response) // Log para depuração
	response = strings.TrimPrefix(response, "Resultado: ")
	response = strings.TrimSpace(response)

	return response, nil
}

func getAvailableOperations() ([]string, error) {
	conn, err := net.Dial("tcp", serverIP)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Envia a mensagem inicial indicando que deseja obter as operações disponíveis
	_, err = conn.Write([]byte("list\n"))
	if err != nil {
		log.Printf("Erro ao enviar mensagem inicial: %v", err)
		return nil, err
	}

	// Lê as operações disponíveis do servidor
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return nil, err
	}

	operations := strings.Split(strings.TrimSpace(response), ",")
	for i, op := range operations {
		operations[i] = fmt.Sprintf("%d. %s", i+1, op)
	}
	operations = append(operations, "5. Sair")

	return operations, nil
}

func confirmRestart(menu *widgets.List, input *widgets.Paragraph, result *widgets.Paragraph, uiEvents <-chan ui.Event) {
	confirm := widgets.NewParagraph()
	confirm.Title = "Pressione Enter para continuar"
	confirm.SetRect(0, 13, 50, 16)
	confirm.TextStyle = ui.NewStyle(ui.ColorYellow)
	ui.Render(confirm)

	for {
		e := <-uiEvents
		if e.ID == "<Enter>" {
			confirm.Text = ""
			ui.Render(confirm)
			break
		}
	}

	// Reinicia o fluxo
	input.Text = ""
	result.Text = ""
	menu.SelectedRow = 0
	ui.Render(menu, input, result)
}
