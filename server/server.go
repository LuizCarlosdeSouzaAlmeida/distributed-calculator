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
	listener, err := net.Listen("tcp", "0.0.0.0:15000")
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
	defer listener.Close()

	log.Println("Servidor iniciado e ouvindo na porta 15000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Erro ao aceitar conexão: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	operations := []string{"somar", "subtrair", "multiplicar", "dividir"}
	conn.Write([]byte(strings.Join(operations, ",") + "\n"))

	reader := bufio.NewReader(conn)
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
			conn.Write([]byte("Erro: Requisição inválida\n"))
			log.Println("Requisição inválida recebida")
			continue
		}

		operation := parts[0]
		numbers := parts[1:]

		if len(numbers) > 20 {
			conn.Write([]byte("Erro: Limite de 20 números excedido\n"))
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
			conn.Write([]byte("Erro: Operação inválida\n"))
			log.Printf("Operação inválida recebida: %s", operation)
			continue
		}

		if !valid {
			conn.Write([]byte("Erro: Entrada inválida\n"))
			log.Println("Entrada inválida recebida")
			continue
		}

		conn.Write([]byte(fmt.Sprintf("Resultado: %.2f\n", result)))
		log.Printf("Operação %s realizada com sucesso. Resultado: %.2f", operation, result)
	}
}

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
