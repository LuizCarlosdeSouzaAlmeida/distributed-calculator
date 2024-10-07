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
		switch operation {
		case "somar":
			result = sum(numbers)
		case "subtrair":
			result = subtract(numbers)
		case "multiplicar":
			result = multiply(numbers)
		case "dividir":
			result = divide(numbers)
		default:
			conn.Write([]byte("Erro: Operação inválida\n"))
			continue
		}

		conn.Write([]byte(fmt.Sprintf("Resultado: %.2f\n", result)))
	}
}

func sum(numbers []string) float64 {
	var result float64
	for _, num := range numbers {
		n, _ := strconv.ParseFloat(num, 64)
		result += n
	}
	return result
}

func subtract(numbers []string) float64 {
	var result float64
	for i, num := range numbers {
		n, _ := strconv.ParseFloat(num, 64)
		if i == 0 {
			result = n
		} else {
			result -= n
		}
	}
	return result
}

func multiply(numbers []string) float64 {
	result := 1.0
	for _, num := range numbers {
		n, _ := strconv.ParseFloat(num, 64)
		result *= n
	}
	return result
}

func divide(numbers []string) float64 {
	var result float64
	for i, num := range numbers {
		n, _ := strconv.ParseFloat(num, 64)
		if i == 0 {
			result = n
		} else {
			result /= n
		}
	}
	return result
}
