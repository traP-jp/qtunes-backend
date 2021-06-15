package main

import (
	"fmt"
	"log"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/hackathon-21-spring-02/back-end/router"
	"github.com/hackathon-21-spring-02/back-end/session"
)

func main() {
	log.Println("Server Started.")

	db, err := model.InitDB()
	if err != nil {
		panic(fmt.Errorf("DB Error: %w", err))
	}

	sess, err := session.NewSession(db.DB)
	if err != nil {
		panic(fmt.Errorf("Session Error: %w", err))
	}

	router.SetRouting(sess)
}
