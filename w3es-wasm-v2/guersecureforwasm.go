package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	shell "github.com/ipfs/go-ipfs-api"
)

type Hashes struct {
	Fid string `json:"fid"`
	Cid string `json:"cid"`
}

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

func encryptData(keystorepass string, keystorehash string, userfile *os.File) []byte {

	keystorePass := keystorepass
	keystoreHash := keystorehash
	userFile := userfile

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(userFile)

	finalurl := join("https://cloudflare-ipfs.com/ipfs/", keystoreHash)

	resp, _ := http.Get(finalurl)

	keystorebody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	aesiv := make([]byte, 32)
	_, err = rand.Read(aesiv)
	if err != nil {
		log.Print(err)
	}
	ec, err := aes.NewCipher(aesiv)
	if err != nil {
		log.Print(err)
	}

	gcm, err := cipher.NewGCM(ec)
	if err != nil {
		log.Print(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Print(err)
	}

	encryptedfile := gcm.Seal(nonce, nonce, buf.Bytes(), nil)

	sh := shell.NewShell("https://ipfs.infura.io:5001")
	fid, err := sh.Add(bytes.NewReader(encryptedfile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	key, err := keystore.DecryptKey(keystorebody, keystorePass)
	if err != nil {
		log.Print(err)
	}

	pubkey := crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	if err != nil {
		log.Print(err)
	}

	dataHash := crypto.Keccak256Hash(encryptedfile)
	if err != nil {
		log.Print("error: write data hash", err)
	}

	signature, err := crypto.Sign(dataHash.Bytes(), key.PrivateKey)
	if err != nil {
		log.Print(err)
	}

	companion := Companion{
		PubKey:    []byte(pubkey[:]),
		Signature: []byte(signature[:]),
		DataHash:  []byte(dataHash.Bytes()),
		AES:       []byte(aesiv[:])}

	jsonCompanion, err := json.Marshal(companion)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf(string(jsonCompanion))

	eciesPubkey := ecies.ImportECDSAPublic(&key.PrivateKey.PublicKey)

	eciesData, err := ecies.Encrypt(rand.Reader, eciesPubkey, []byte(jsonCompanion), nil, nil)
	if err != nil {
		log.Print(err)
	}

	csh := shell.NewShell("https://ipfs.infura.io:5001")
	cid, err := csh.Add(bytes.NewReader(eciesData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	hashes := Hashes{Fid: fid, Cid: cid}

	result, err := json.Marshal(hashes)

	if err != nil {
		fmt.Println(err)
	}

	return result
}

func decryptData(keystorepass string, keystorehash string, companionhash string, filehash string) []byte {

	keystorePass := keystorepass
	keystoreHash := keystorehash
	companionHash := companionhash
	fileHash := filehash

	//plaintextWriteName := "./tmp/unencryptedfile.png"
	keystoreurl := join("https://cloudflare-ipfs.com/ipfs/", keystoreHash)
	ciphertexturl := join("https://cloudflare-ipfs.com/ipfs/", fileHash)
	companionurl := join("https://cloudflare-ipfs.com/ipfs/", companionHash)

	keystoreresp, _ := http.Get(keystoreurl)

	keystorebody, err := ioutil.ReadAll(keystoreresp.Body)
	if err != nil {
		log.Println("1", err)
	}

	ciphertextresp, _ := http.Get(ciphertexturl)

	ciphertextbody, err := ioutil.ReadAll(ciphertextresp.Body)
	if err != nil {
		log.Println("2", err)
	}

	companionresp, _ := http.Get(companionurl)

	companionbody, err := ioutil.ReadAll(companionresp.Body)
	if err != nil {
		log.Println("3", err)
	}

	key, err := keystore.DecryptKey(keystorebody, keystorePass)
	if err != nil {
		log.Print("4", err)
	}

	eciesKey := ecies.ImportECDSA(key.PrivateKey)

	jsonCompanion, err := eciesKey.Decrypt(companionbody, nil, nil)
	if err != nil {
		log.Print("5", err)
	}

	var companion Companion

	json.Unmarshal(jsonCompanion, &companion)

	//Commenting this bit out for now --- sig verification
	//dataHashDecrypt := crypto.Keccak256Hash(ciphertextbody)
	//if err != nil {
	//	log.Print("error: write data hash", err)
	//}
	//
	//sigNoRecoverID := companion.Signature[:len(companion.Signature)-1]

	aeskey := []byte(companion.AES)

	dc, err := aes.NewCipher(aeskey)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(dc)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertextbody) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertextbody[:nonceSize], ciphertextbody[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}

	return plaintext
}

func main() {

}
