package main

import (
	"context"
	"fmt"

	"github.com/kaito2/getting-started/bob/models"
	"github.com/stephenafamo/bob"
)

func main() {
	db, err := bob.Open("mysql", "root:pass@tcp(localhost:3306)/example")
	if err != nil {
		panic(err)
	}

	user, err := models.FindUser(context.Background(), db, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	users, err := models.Users.Query(
		models.SelectWhere.Users.ID.GT(1),
	).All(context.Background(), db)
	if err != nil {
		panic(err)
	}
	for i, user := range users {
		fmt.Printf("user[%d]:%v\n", i, user)
	}
}
