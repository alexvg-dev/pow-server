package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"net"
	"strings"
)

const (
	port       = "4444"
	difficulty = 2 // Количество нулей в начале хэша
)

func main() {
	conn, err := net.Dial("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Отправляем HelloRequest
	err = writeWithLength(conn, []byte("Hello"))
	if err != nil {
		fmt.Println("Error sending Hello request:", err)
		return
	}

	// Читаем challenge от сервера
	buf, err := readWithLength(conn)
	if err != nil {
		fmt.Println("Error reading challenge:", err)
		return
	}

	challenge := hex.EncodeToString(buf)
	fmt.Println("Received challenge:", challenge)

	// Генерируем ответ на challenge
	response := solveChallenge(buf)
	fmt.Println("generated response:", hex.EncodeToString(response))

	// Отправляем ответ на challenge
	err = writeWithLength(conn, response)
	if err != nil {
		fmt.Println("Error sending response:", err)
		return
	}

	// Читаем ответ от сервера
	buf, err = readWithLength(conn)
	if err != nil {
		fmt.Println("Error reading server response:", err)
		return
	}

	fmt.Println("server response:", string(buf))
}

func solveChallenge(challenge []byte) []byte {

	nonce := 0 // Начинаем с нуля
	for {
		// Преобразуем nonce в байты
		nonceBytes := []byte(fmt.Sprintf("%d", nonce))

		// Вычисляем хэш с использовоНогом
		hash, err := scrypt.Key(nonceBytes, challenge, 16384, 8, 1, 32)
		if err != nil {
			fmt.Println("Error computing scrypt:", err)
			return nil
		}

		hashStr := hex.EncodeToString(hash)

		// Проверяем, соответствует ли хэш условиям
		if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
			fmt.Printf("Solution found: %s (nonce=%s)\n", hashStr, nonceBytes)
			return nonceBytes
		}

		nonce++ // Увеличиваем nonce для следующей итерации
	}
}

func writeWithLength(conn net.Conn, data []byte) error {
	// Отправляем размер сообщения
	err := binary.Write(conn, binary.BigEndian, uint64(len(data)))
	if err != nil {
		fmt.Println("Error writing message size:", err)
		return err
	}

	// Отправляем само сообщение
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func readWithLength(conn net.Conn) ([]byte, error) {
	var length uint64
	err := binary.Read(conn, binary.BigEndian, &length)
	if err != nil {
		return nil, nil
	}

	fmt.Printf("Got message. Len = %d\n", length)

	// Читаем само сообщение
	data := make([]byte, length)
	_, err = conn.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return data, nil
}
