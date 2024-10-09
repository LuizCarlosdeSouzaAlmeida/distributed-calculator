package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	// Inicia o servidor e ouve na porta 15000
	listener, err := startServer("0.0.0.0:15000")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
	defer listener.Close()

	log.Println("Servidor iniciado e ouvindo na porta 15000...")

	// Loop principal para aceitar conexões
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

// startServer inicia o servidor e retorna o listener
func startServer(address string) (net.Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// handleConnection lida com uma conexão individual
func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Lê a mensagem inicial do cliente
	initialMessage, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Erro ao ler mensagem inicial: %v", err)
		return
	}

	initialMessage = strings.TrimSpace(initialMessage)

	switch initialMessage {
	case "list":
		// Envia a lista de operações disponíveis
		sendOperationsList(conn)
		return
	case "operation":
		// Continua para o processamento de operações
	default:
		sendError(conn, "Comando inicial inválido")
		log.Printf("Comando inicial inválido recebido: %s", initialMessage)
		return
	}

	// Loop para processar requisições
	for {
		request, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Conexão fechada pelo cliente")
			} else {
				log.Printf("Erro ao ler requisição: %v", err)
			}
			return
		}

		request = strings.TrimSpace(request)
		if request == "" {
			log.Println("Requisição vazia recebida")
			continue
		}

		parts := strings.Split(request, " ")
		if len(parts) < 2 {
			sendError(conn, "Requisição inválida")
			log.Println("Requisição inválida recebida")
			continue
		}

		operation := parts[0]
		numbers := parts[1:]

		if len(numbers) > 20 {
			sendError(conn, "Limite de 20 números excedido")
			log.Println("Limite de números excedido")
			continue
		}

		var result float64
		var valid bool
		switch operation {
		case "somar":
			result, valid = sum(numbers)
		case "subtrair":
			result, valid = subtract(numbers)
		case "multiplicar":
			result, valid = multiply(numbers)
		case "dividir":
			result, valid = divide(numbers)
		default:
			sendError(conn, "Operação inválida")
			log.Printf("Operação inválida recebida: %s", operation)
			continue
		}

		if !valid {
			sendError(conn, "Entrada inválida")
			log.Println("Entrada inválida recebida")
			continue
		}

		sendResult(conn, result)
		log.Printf("Operação %s realizada com sucesso. Resultado: %.2f", operation, result)
	}
}

// sendOperationsList envia a lista de operações disponíveis para o cliente
func sendOperationsList(conn net.Conn) {
	operations := []string{"somar", "subtrair", "multiplicar", "dividir"}
	conn.Write([]byte(strings.Join(operations, ",") + "\n"))
}

// sendError envia uma mensagem de erro para o cliente
func sendError(conn net.Conn, message string) {
	conn.Write([]byte(fmt.Sprintf("Erro: %s\n", message)))
}

// sendResult envia o resultado da operação para o cliente
func sendResult(conn net.Conn, result float64) {
	conn.Write([]byte(fmt.Sprintf("Resultado: %.2f\n", result)))
}

// sum realiza a soma dos números
func sum(numbers []string) (float64, bool) {
	var result float64
	for _, num := range numbers {
		n, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0, false
		}
		result += n
	}
	return result, true
}

// subtract realiza a subtração dos números
func subtract(numbers []string) (float64, bool) {
	var result float64
	for i, num := range numbers {
		n, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0, false
		}
		if i == 0 {
			result = n
		} else {
			result -= n
		}
	}
	return result, true
}

// multiply realiza a multiplicação dos números
func multiply(numbers []string) (float64, bool) {
	result := 1.0
	for _, num := range numbers {
		n, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0, false
		}
		result *= n
	}
	return result, true
}

// divide realiza a divisão dos números
func divide(numbers []string) (float64, bool) {
	var result float64
	for i, num := range numbers {
		n, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0, false
		}
		if i == 0 {
			result = n
		} else {
			if n == 0 {
				return 0, false
			}
			result /= n
		}
	}
	return result, true
}
