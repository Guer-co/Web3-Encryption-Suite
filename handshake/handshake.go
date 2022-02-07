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

	"github.com/ethereum/go-ethereum/crypto"
	ecies "github.com/ethereum/go-ethereum/crypto/ecies"
)

type Request struct {
	PartAPubkey      []byte `json:"partApubkey"`
	SignatureReq     []byte `json:"signatureRequest"`
	SignatureReqData []byte `json:"signatureReqData"`
}

type Confirm struct {
	PartBPubkey      []byte `json:"partBpubkey"`
	SignatureCon     []byte `json:"signatureConfirm"`
	SignatureConData []byte `json:"signatureConData"`
}

type Distribution struct {
	PartASymkey []byte `json:"sym key for requester"`
	PartBSymkey []byte `json:"sym key for confirmer"`
}

func main() {

	//UI inputs
	requestMsgPath := "./tmp/request.json"
	confirmMsgPath := "./tmp/confirm.json"

	//import request
	importRequestMsg, err := ioutil.ReadFile(requestMsgPath)
	if err != nil {
		fmt.Println(err)
	}

	//unmarshal json
	var request Request

	json.Unmarshal(importRequestMsg, &request)

	//verify sig
	sigReqNoRecoverID := request.SignatureReq[:len(request.SignatureReq)-1]

	verifiedReq := crypto.VerifySignature(request.PartAPubkey, request.SignatureReqData, sigReqNoRecoverID)
	if verifiedReq == false {
		fmt.Println("Invalid Signature. Ending Process")
		log.Fatal(verifiedReq)
	} else {
		fmt.Println("Valid Signature: ", verifiedReq)
	}

	//import confirmation
	importConfirmMsg, err := ioutil.ReadFile(confirmMsgPath)
	if err != nil {
		fmt.Println(err)
	}

	//unmarshal json
	var confirm Confirm

	json.Unmarshal(importConfirmMsg, &confirm)

	//verify sig
	sigConNoRecoverID := confirm.SignatureCon[:len(confirm.SignatureCon)-1]

	verifiedCon := crypto.VerifySignature(confirm.PartBPubkey, confirm.SignatureConData, sigConNoRecoverID)
	if verifiedReq == false {
		fmt.Println("Invalid Signature. Ending Process")
		log.Fatal(verifiedCon)
	} else {
		fmt.Println("Valid Signature: ", verifiedCon)
	}

	//gen rand aes key
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

	//ecies encrypt part A

	partAHex, err := crypto.UnmarshalPubkey(request.PartAPubkey)
	if err != nil {
		fmt.Println(err)
	}

	partAEciesPubkey := ecies.ImportECDSAPublic(partAHex)

	eciesPartA, err := ecies.Encrypt(rand.Reader, partAEciesPubkey, []byte(aesiv), nil, nil)
	if err != nil {
		fmt.Println("eciesPartA failed")
		log.Fatal(err)
	}

	//ecies encrypt part B
	partBHex, err := crypto.UnmarshalPubkey(confirm.PartBPubkey)
	if err != nil {
		fmt.Println(err)
	}

	partBEciesPubkey := ecies.ImportECDSAPublic(partBHex)

	eciesPartB, err := ecies.Encrypt(rand.Reader, partBEciesPubkey, []byte(aesiv), nil, nil)
	if err != nil {
		fmt.Println("eciesPartB failed")
		log.Fatal(err)
	}

	//marshal json part A key, part B key
	guerKeys := Distribution{
		PartASymkey: []byte(eciesPartA[:]),
		PartBSymkey: []byte(eciesPartB[:]),
	}

	jsonGuerKeys, err := json.Marshal(guerKeys)
	if err != nil {
		log.Fatal(err)
	}
	//write json key distribution
	ioutil.WriteFile("./tmp/guerKeys.json", jsonGuerKeys, 0644)
}
