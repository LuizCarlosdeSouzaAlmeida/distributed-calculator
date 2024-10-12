package br.ufg.inf.es.rcsd.lab1;

import java.io.*;
import java.net.*;
import java.util.StringTokenizer;

class Server {

    public static void main(String argv[]) throws Exception {
        String clientSentence;
        String result;
        int porta = 15000; // Porta definida para 15000
        ServerSocket welcomeSocket = new ServerSocket(porta);

        System.out.println("Esperando conexões na porta " + porta);

        int numConn = 1;

        while (true) {

            Socket connectionSocket = welcomeSocket.accept();
            System.out.println("Conexão [" + numConn++ + "] aceita...");

            BufferedReader inFromClient = new BufferedReader(
                    new InputStreamReader(connectionSocket.getInputStream()));

            DataOutputStream outToClient = new DataOutputStream(
                    connectionSocket.getOutputStream());

            clientSentence = inFromClient.readLine();

            // A operação é recebida e processada
            result = processOperation(clientSentence);

            outToClient.writeBytes(result + '\n');

            connectionSocket.close();
        }
    }

    // Função para processar a operação aritmética
    private static String processOperation(String input) {
        StringTokenizer st = new StringTokenizer(input);
        String operation = st.nextToken();
        double result = 0;
        boolean first = true;

        switch (operation) {
            case "SOMAR":
                while (st.hasMoreTokens()) {
                    double num = Double.parseDouble(st.nextToken());
                    result += num;
                }
                break;
            case "SUBTRAIR":
                if (st.hasMoreTokens()) {
                    result = Double.parseDouble(st.nextToken());
                }
                while (st.hasMoreTokens()) {
                    result -= Double.parseDouble(st.nextToken());
                }
                break;
            case "MULTIPLICAR":
                result = 1;
                while (st.hasMoreTokens()) {
                    result *= Double.parseDouble(st.nextToken());
                }
                break;
            case "DIVIDIR":
                if (st.hasMoreTokens()) {
                    result = Double.parseDouble(st.nextToken());
                }
                while (st.hasMoreTokens()) {
                    double num = Double.parseDouble(st.nextToken());
                    if (num != 0) {
                        result /= num;
                    } else {
                        return "Erro: divisao por zero.";
                    }
                }
                break;
            default:
                return "Operação inválida";
        }
        return "Resultado: " + result;
    }
}
