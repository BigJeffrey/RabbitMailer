package main

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Nasłuchuję...")
	from := "toppizzapl@77.55.216.72"
	//password := ""

	connRabbit, err := amqp.Dial("amqp://130.61.54.93:5672/")
	if err != nil {
		log.Fatal(err)
	}

	defer connRabbit.Close()

	chRabbit, err := connRabbit.Channel()
	if err != nil {
		log.Fatal(err)
	}

	defer chRabbit.Close()

	q, err := chRabbit.QueueDeclare(
		"login", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := chRabbit.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Otrzymałem wiadomość: %s", d.Body)
			address := string(d.Body)
			to := []string{
				address,
			}

			smtpHost := "77.55.216.72"
			smtpPort := "1025"

			message := []byte("Witaj w TopPizza nowy użytkowniku!!!")

			//auth := smtp.PlainAuth("", from, password, smtpHost)

			err := smtp.SendMail(smtpHost+":"+smtpPort, nil, from, to, message)

			if err != nil {
				//log.Fatal(err)
				log.Println(err)
				fmt.Println("Nie udało się wysłać maila")
			} else {
				fmt.Println("Email wysłany")
			}

		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	<-forever
}
