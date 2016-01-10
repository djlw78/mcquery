package mcquery_test

import (
	"fmt"
	"time"

	"github.com/sean-callahan/mcquery"
)

func ExampleDial() {
	_, err := mcquery.Dial("127.0.0.1:25565", time.Second)
	if err != nil {
		panic(err)
	}
}

func ExampleMcQuery_GetStatus() {
	mcq, err := mcquery.Dial("127.0.0.1:25565", time.Second)
	if err != nil {
		panic(err)
	}

	status, _, err := mcq.GetStatus()
	if err != nil {
		panic(err)
	}

	fmt.Println(status["game_id"])
}
