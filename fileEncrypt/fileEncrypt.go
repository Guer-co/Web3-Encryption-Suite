package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	ecies "github.com/ethereum/go-ethereum/crypto/ecies"
)

type Companion struct {
	PubKey    []byte `json:"pubkey"`
	Signature []byte `json:"signature"`
	DataHash  []byte `json:"datahash"`
	AES       []byte `json:"aesiv"`
}

func main() {

	fmt.Println("GUER Encryption A v0.01")

	//UI Inputs below:
	filePath := "./tmp/dog.txt"
	keystorePath := "./tmp/UTC--2021-06-11T10-03-17.844395300Z--71ab8dfbd056b92da467d61bf7171f77558980cc"
	keystorePass := "motdepasse"
	encryptedPath := "./tmp/myfile.data"

	//import file
	plaintext, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	//generate aes key from rand
	aesiv := make([]byte, 32)
	_, err = rand.Read(aesiv)
	if err != nil {
		log.Fatal(err)
		return
	}
	c, err := aes.NewCipher(aesiv)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	//encrypt file using aes key and write to file
	err = ioutil.WriteFile("./tmp/myfile.data", gcm.Seal(nonce, nonce, plaintext, nil), 0777)
	if err != nil {
		log.Fatal(err)
	}

	//import keystore
	//**UI**: User should be able to specify keystore file and/or eventually sign with wallet
	importKeystore, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		log.Fatal(err)
	}

	//**UI**: User needs to provide pw to their keystore file; alt: be logged into wallet
	password := keystorePass

	//decrypts key from keystore for use in application
	key, err := keystore.DecryptKey(importKeystore, password)
	if err != nil {
		log.Fatal(err)
	}

	//write to file pubKey from imported keystore
	pubkey := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	//read encrypted file for HMAC hash
	data, err := ioutil.ReadFile(encryptedPath)
	if err != nil {
		log.Fatal("error: read encrypted file", err)
	}

	//hash encrypted file
	dataHash := crypto.Keccak256Hash(data)
	if err != nil {
		log.Fatal("error: write data hash", err)
	}

	//sign hmac
	signature, err := crypto.Sign(dataHash.Bytes(), key.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	companion := Companion{
		PubKey:    []byte(pubkey[:]),
		Signature: []byte(signature[:]),
		DataHash:  []byte(dataHash.Bytes()),
		AES:       []byte(aesiv[:])}

	jsonCompanion, err := json.Marshal(companion)
	if err != nil {
		log.Fatal(err)
	}

	eciesPubkey := ecies.ImportECDSAPublic(&key.PrivateKey.PublicKey)

	eciesData, err := ecies.Encrypt(rand.Reader, eciesPubkey, []byte(jsonCompanion), nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile("./tmp/companion.json", eciesData, 0644)
}
