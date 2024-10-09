# Distributed Calculator

Este projeto consiste em um sistema distribuído de calculadora, onde o servidor processa operações matemáticas básicas (soma, subtração, multiplicação e divisão) e o cliente se conecta ao servidor para realizar essas operações. O cliente obtém as operações disponíveis diretamente do servidor, garantindo que as informações sejam sempre atualizadas.

## Estrutura do Projeto

O projeto é dividido em dois componentes principais:

1. **Servidor**: Responsável por receber as requisições do cliente, processar as operações matemáticas e retornar os resultados.
2. **Cliente**: Interface de linha de comando (CLI) que se conecta ao servidor para realizar as operações matemáticas.

## Pré-requisitos

Antes de começar, certifique-se de que você tem o seguinte instalado:

- [Go](https://golang.org/dl/) (versão 1.16 ou superior)
- Conexão à internet (para o cliente se conectar ao servidor)

## Configuração do Servidor

### 1. Clone o Repositório

Primeiro, clone o repositório para o seu ambiente local:

```bash
git clone https://github.com/seu-usuario/distributed-calculator.git
cd distributed-calculator
```

### 2. Navegue até o Diretório do Servidor

```bash
cd server
```

### 3. Execute o Servidor

O servidor pode ser executado diretamente usando o comando `go run`:

```bash
go run main.go
```

O servidor estará agora em execução e ouvindo na porta `15000`.

## Configuração do Cliente

### 1. Navegue até o Diretório do Cliente

```bash
cd ../client
```

### 2. Instale o Cliente

Para instalar o cliente, use o comando `go install`:

```bash
go install
```

Isso compilará o cliente e o instalará no seu `$GOPATH/bin`.

### 3. Execute o Cliente

Agora, você pode executar o cliente usando o comando `go run`:

```bash
go run main.go
```

O cliente se conectará ao servidor e exibirá um menu com as operações disponíveis.

## Como Usar o Cliente

1. **Selecione uma Operação**: Use as setas para cima e para baixo para navegar pelo menu e pressione `Enter` para selecionar uma operação.
2. **Digite os Números**: Após selecionar uma operação, você será solicitado a digitar os números separados por espaço.
3. **Obtenha o Resultado**: O cliente enviará a operação e os números para o servidor, que processará a operação e retornará o resultado.
4. **Sair**: Para sair do cliente, selecione a opção "5. Sair" no menu.

## Configuração do IP do Servidor

O cliente está configurado para se conectar ao servidor no endereço `3.225.60.216:15000`. Se você estiver executando o servidor localmente ou em um endereço diferente, você pode modificar o IP do servidor no arquivo `client/main.go`.

### Modificando o IP do Servidor

Abra o arquivo `client/main.go` e localize a constante `serverIP`:

```go
const serverIP = "3.225.60.216:15000"
```

Altere o valor da constante `serverIP` para o endereço do servidor que você está usando. Por exemplo, se o servidor estiver em execução localmente, você pode alterar para:

```go
const serverIP = "localhost:15000"
```

## Exemplo de Uso

1. **Selecione a Operação**:
   - Use as setas para cima e para baixo para selecionar "1. somar".
   - Pressione `Enter`.

2. **Digite os Números**:
   - Digite `1 2 3` e pressione `Enter`.

3. **Resultado**:
   - O cliente exibirá o resultado: `Resultado: 6.00`.

## Considerações Finais

- **Execução do Servidor**: Certifique-se de que o servidor esteja em execução antes de iniciar o cliente.
- **Configuração do IP**: O cliente se conecta ao servidor no endereço definido na constante `serverIP`. Se o servidor estiver em execução localmente, você pode modificar o endereço no código do cliente para `localhost:15000`.
- **Limitações**: O servidor suporta até 20 números por operação.

## Contribuição

Sinta-se à vontade para contribuir com melhorias, correções de bugs ou novas funcionalidades. Abra uma issue ou envie um pull request.

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).