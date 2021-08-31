package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

type Companion struct {
	PubKey    []byte `json:"pubkey"`
	Signature []byte `json:"signature"`
	DataHash  []byte `json:"datahash"`
	AES       []byte `json:"aesiv"`
}

func main() {

	fmt.Println("GUER Decryption A v0.01")

	//Required User Inputs:
	ciphertextReadFilepath := "./tmp/myfile.data"
	keystorePath := "./tmp/UTC--2021-06-11T10-03-17.844395300Z--71ab8dfbd056b92da467d61bf7171f77558980cc"
	keystorePass := "motdepasse"
	companionImportFilepath := "./tmp/companion.json"
	plaintextWriteName := "./tmp/unencryptedDog.txt"

	//Import keystore
	keystoreImport, err := ioutil.ReadFile(keystorePath)
	if err != nil {
		log.Fatal(err)
	}

	//decrypts key from keystore for use in application
	key, err := keystore.DecryptKey(keystoreImport, keystorePass)
	if err != nil {
		log.Fatal(err)
	}

	eciesKey := ecies.ImportECDSA(key.PrivateKey)

	//Import jsoncompanion
	companionImport, err := ioutil.ReadFile(companionImportFilepath)
	if err != nil {
		log.Fatal(err)
	}

	//Decrypt jsoncompanion using ECIES
	jsonCompanion, err := eciesKey.Decrypt(companionImport, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Parse JSON

	var companion Companion

	json.Unmarshal(jsonCompanion, &companion)

	//Import encrypted data
	ciphertextRead, err := ioutil.ReadFile(ciphertextReadFilepath)
	if err != nil {
		log.Fatal("error: read encrypted file", err)
	}

	//generate hash from encrypted file for sig verification
	dataHashDecrypt := crypto.Keccak256Hash(ciphertextRead)
	if err != nil {
		log.Fatal("error: write data hash", err)
	}

	//verify signature of file you're trying to decrypt
	sigNoRecoverID := companion.Signature[:len(companion.Signature)-1]

	verified := crypto.VerifySignature(companion.PubKey, dataHashDecrypt.Bytes(), sigNoRecoverID)
	if verified == false {
		log.Fatal(verified)
		fmt.Println("Invalid Signature. Ending Process")
	} else {
		fmt.Println("Valid Signature: ", verified)
	}

	//Decrypt file
	aeskey := []byte(companion.AES)

	c, err := aes.NewCipher(aeskey)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertextRead) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertextRead[:nonceSize], ciphertextRead[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	//Read file to terminal
	fmt.Println(string(plaintext))

	//Write plaintext file
	ioutil.WriteFile(plaintextWriteName, plaintext, 0664)

	//conclude decrypt
	fmt.Println("Successfully Decrypted Data")
	log.Fatal()
}
