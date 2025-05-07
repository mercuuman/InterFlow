package main

import "log"

func main() {
	err := sendMail("devatkinkonstantin@yandex.ru", "1312dsadasd123123")
	if err != nil {
		log.Fatal(err) // Выведет ошибку и завершит программу
	}
}
