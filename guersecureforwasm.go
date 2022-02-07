package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	_ "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "io"
	"log"
	"strings"
	"syscall/js"
	_ "syscall/js"

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

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func encryptData(this js.Value, inputs []js.Value) interface{} {
	var err error
	keystorePass := inputs[0].String()
	keystoreBodystring := inputs[1].String()
	keystoreBody := []byte(keystoreBodystring)
	file := inputs[2].String()
	userFile, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		fmt.Println(err)
	}

	aesiv := make([]byte, 32)
	_, err = rand.Read(aesiv)

	ec, err := aes.NewCipher(aesiv)

	gcm, err := cipher.NewGCM(ec)

	nonce := make([]byte, gcm.NonceSize())

	key, err := keystore.DecryptKey(keystoreBody, keystorePass)

	pubkey := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)

	encryptedFile := gcm.Seal(nonce, nonce, userFile, nil)

	dataHash := crypto.Keccak256Hash(encryptedFile)

	signature, err := crypto.Sign(dataHash.Bytes(), key.PrivateKey)

	companion := Companion{
		PubKey:    []byte(pubkey[:]),
		Signature: []byte(signature[:]),
		DataHash:  []byte(dataHash.Bytes()),
		AES:       []byte(aesiv[:]),
	}

	jsonCompanion, err := json.Marshal(companion)
	eciesPubkey := ecies.ImportECDSAPublic(&key.PrivateKey.PublicKey)

	companionFile, err := ecies.Encrypt(rand.Reader, eciesPubkey, []byte(jsonCompanion), nil, nil)

	b64c := base64.StdEncoding.EncodeToString(companionFile)

	b64f := base64.StdEncoding.EncodeToString(encryptedFile)

	js.Global().Set("companionFile", string(b64c))
	js.Global().Set("encryptedFile", string(b64f))

	return nil
}

func decryptData(this js.Value, inputs []js.Value) interface{} {
	var err error
	if err != nil {
		fmt.Println(err)
	}
	keystorePass := inputs[0].String()
	keystoreBody := inputs[1].String()
	encryptedFile := inputs[2].String()
	companionBody := inputs[3].String()

	log.Print("companionBase64")
	log.Print(companionBody)

	log.Print("encryptedBase64")
	log.Print(encryptedFile)

	cd, err := base64.StdEncoding.DecodeString(companionBody)
	fd, err := base64.StdEncoding.DecodeString(encryptedFile)
	log.Print("companionbase64decrypt")
	log.Print(cd)
	dkey, err := keystore.DecryptKey([]byte(keystoreBody), keystorePass)

	eciesKey := ecies.ImportECDSA(dkey.PrivateKey)

	djsonCompanion, err := eciesKey.Decrypt(cd, nil, nil)

	var companion2 Companion

	json.Unmarshal(djsonCompanion, &companion2)

	aeskey := []byte(companion2.AES)

	dc, err := aes.NewCipher(aeskey)

	dgcm, err := cipher.NewGCM(dc)

	nonceSize := dgcm.NonceSize()

	nonce, ciphertext := fd[:nonceSize], fd[nonceSize:]

	plaintext, err := dgcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)

	log.Print(plaintext)

	return base64.StdEncoding.EncodeToString(plaintext)
}

func eAndD(this js.Value, inputs []js.Value) interface{} {
	var err error
	keystorePass := inputs[0].String()
	keystoreBodystring := inputs[1].String()
	keystoreBody := []byte(keystoreBodystring)
	file := inputs[2].String()
	userFile, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		fmt.Println(err)
	}

	aesiv := make([]byte, 32)
	_, err = rand.Read(aesiv)

	ec, err := aes.NewCipher(aesiv)

	gcm, err := cipher.NewGCM(ec)

	nonce := make([]byte, gcm.NonceSize())

	key, err := keystore.DecryptKey(keystoreBody, keystorePass)

	pubkey := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)

	encryptedFile := gcm.Seal(nonce, nonce, userFile, nil)

	dataHash := crypto.Keccak256Hash(encryptedFile)

	signature, err := crypto.Sign(dataHash.Bytes(), key.PrivateKey)

	companion := Companion{
		PubKey:    []byte(pubkey[:]),
		Signature: []byte(signature[:]),
		DataHash:  []byte(dataHash.Bytes()),
		AES:       []byte(aesiv[:]),
	}

	jsonCompanion, err := json.Marshal(companion)
	eciesPubkey := ecies.ImportECDSAPublic(&key.PrivateKey.PublicKey)

	companionFile, err := ecies.Encrypt(rand.Reader, eciesPubkey, []byte(jsonCompanion), nil, nil)
	log.Print("companion-doall-encrypted .")
	log.Print(companionFile)

	b64c := base64.StdEncoding.EncodeToString(companionFile)
	log.Print("companion-doall-base64 .")
	log.Print(b64c)

	log.Print("companion-doall-base64-cast-to-string .")
	log.Print(string(b64c))

	b64f := base64.StdEncoding.EncodeToString(encryptedFile)

	//this is where decrypt starts

	cd, err := base64.StdEncoding.DecodeString(b64c)
	log.Print(cd)

	fd, err := base64.StdEncoding.DecodeString(b64f)

	dkey, err := keystore.DecryptKey([]byte(keystoreBody), keystorePass)

	eciesKey := ecies.ImportECDSA(dkey.PrivateKey)

	log.Print(cd)
	djsonCompanion, err := eciesKey.Decrypt(cd, nil, nil)

	var companion2 Companion

	json.Unmarshal(djsonCompanion, &companion2)

	aeskey := []byte(companion2.AES)

	dc, err := aes.NewCipher(aeskey)

	dgcm, err := cipher.NewGCM(dc)

	nonceSize := dgcm.NonceSize()

	nonce, ciphertext := fd[:nonceSize], fd[nonceSize:]

	plaintext, err := dgcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)

	log.Print(plaintext)

	return base64.StdEncoding.EncodeToString(plaintext)
}

func main() {
	c := make(chan int)
	js.Global().Set("encryptData", js.FuncOf(encryptData))
	js.Global().Set("decryptData", js.FuncOf(decryptData))
	js.Global().Set("eAndD", js.FuncOf(eAndD))

	<-c
}
