package br.ufg.inf.es.rcsd.lab1;

import java.io.*;
import java.net.*;

class Client {

    public static void main(String argv[]) throws Exception {
        BufferedReader inFromUser = new BufferedReader(new InputStreamReader(System.in));

        int porta = 15000; // Porta definida para 15000
        String servidor = "localhost";

        while (true) {
            System.out.println("\nEscolha a operação:");
            System.out.println("1. Somar");
            System.out.println("2. Subtrair");
            System.out.println("3. Multiplicar");
            System.out.println("4. Dividir");
            System.out.println("5. Sair");

            String choice = inFromUser.readLine();
            String operation = "";
            if (choice.equals("1")) {
                operation = "SOMAR";
            } else if (choice.equals("2")) {
                operation = "SUBTRAIR";
            } else if (choice.equals("3")) {
                operation = "MULTIPLICAR";
            } else if (choice.equals("4")) {
                operation = "DIVIDIR";
            } else if (choice.equals("5")) {
                System.out.println("Encerrando...");
                break;
            } else {
                System.out.println("Escolha inválida.");
                continue;
            }

            System.out.println("Digite os números separados por espaço:");
            String numbers = inFromUser.readLine();

            // Conectar ao servidor
            Socket clientSocket = new Socket(servidor, porta);
            DataOutputStream outToServer = new DataOutputStream(clientSocket.getOutputStream());
            BufferedReader inFromServer = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()));

            // Enviar a operação e os números
            String sentence = operation + " " + numbers;
            outToServer.writeBytes(sentence + '\n');

            // Receber e imprimir o resultado
            String modifiedSentence = inFromServer.readLine();
            System.out.println("Recebido do servidor: " + modifiedSentence);

            clientSocket.close();
        }
    }
}
