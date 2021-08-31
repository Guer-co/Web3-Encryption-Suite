package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"io/ioutil"
	"log"
	"os"
)

func createKs() {
	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "motdepasse"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

}

func importKs() {
	file := "./tmp/UTC--2021-06-11T10-03-17.844395300Z--71ab8dfbd056b92da467d61bf7171f77558980cc"
	ks := keystore.NewKeyStore("./tmp", keystore.LightScryptN, keystore.LightScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	password := "motdepasse"
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex())

	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}

}

func main() {
	createKs()
}
