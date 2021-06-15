package main

import (
	"fmt"
	"log"

	"github.com/hackathon-21-spring-02/back-end/model"
	"github.com/hackathon-21-spring-02/back-end/router"
	sess "github.com/hackathon-21-spring-02/back-end/session"
)

func main() {
	log.Println("Server Started.")

	db, err := model.InitDB()
	if err != nil {
		panic(fmt.Errorf("DB Error: %w", err)) //TODO
	}

	sess, err := sess.NewSession(db.DB)
	if err != nil {
		panic(err) //TODO
	}

	router.SetRouting(sess)
}
