package main

import (
	"fmt"
	db "github.com/BrandonRomano/drudge/database"
	"github.com/BrandonRomano/drudge/models"
)

func main() {
	db.Open()
	defer db.Close()

	rigby := &models.Animal{
		Name: "Rigby",
		Age:  3,
	}
	rigby.DbWorker = models.DbWorker{
		BaseModel: rigby,
	}

	fmt.Println(rigby)
	rigby.DbWorker.Insert()
	fmt.Println(rigby)

	fmt.Println("--------")

	rigbyTwo := &models.Animal{
		Id: rigby.Id,
	}
	rigbyTwo.DbWorker = models.DbWorker{
		BaseModel: rigbyTwo,
	}

	fmt.Println(rigbyTwo)
	rigbyTwo.DbWorker.Load()
	fmt.Println(rigbyTwo)
}
