package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:15000")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor iniciado e ouvindo na porta 15000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		request, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler requisição:", err)
			return
		}

		request = strings.TrimSpace(request)
		if request == "" {
			continue
		}

		parts := strings.Split(request, " ")
		if len(parts) < 2 {
			conn.Write([]byte("Erro: Requisição inválida\n"))
			continue
		}

		operation := parts[0]
		numbers := parts[1:]

		if len(numbers) > 20 {
			conn.Write([]byte("Erro: Limite de 20 números excedido\n"))
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
			continue
		}

		if !valid {
			conn.Write([]byte("Erro: Entrada inválida\n"))
			continue
		}

		conn.Write([]byte(fmt.Sprintf("Resultado: %.2f\n", result)))
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
