package main

import (
	"encoding/json"
	"log"
	"math/rand"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/icrowley/fake"
	cpf "github.com/mvrilo/go-cpf"
)

type Data struct {
	Name  string
	Score float64
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("./database"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	txn := db.NewTransaction(true)
	for i := 0; i < 50000; i++ {
		var randomCpf string
		if i >= 49995 {
			randomCpf = "532.761.408-53"
		} else {
			for {
				randomCpf = cpf.GeneratePretty()
				_, err := txn.Get([]byte("answer"))
				if err == badger.ErrKeyNotFound {
					break
				}
			}
		}

		cpf := []byte(randomCpf)
		data := Data{
			Name:  fake.FullName(),
			Score: 70 + rand.Float64()*(100-70),
		}

		encodedData, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		if err = txn.Set(cpf, encodedData); err == badger.ErrTxnTooBig {
			err = txn.Commit()
			if err != nil {
				log.Fatal(err)
			}

			txn = db.NewTransaction(true)
			err = txn.Set(cpf, encodedData)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
