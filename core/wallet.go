package core

import (
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"bytes"
	"fmt"
	"time"
	"strconv"
	//"log"
	"os"

	"../wallet"
)


const walletFile="wallet_%s.dat"

type WepWallet struct {
	Timestamp	int
	Address		string
	PublicKey	string
	PrivateKey	string
	PrivateHash	string
}

func (ww *WepWallet) Create(pass string) {
	w:=wallet.CreateWallet()
	ww.Timestamp=int(time.Now().Unix())
	ww.PublicKey=ByteToStr(w.PublicKey)
	ww.PrivateKey=fmt.Sprintf("%d",w.PrivateKey.D)
	ww.Address=ByteToStr(w.GetAddress())
	if(string(ww.Address[0:1])!="W") {
		ww.Create(pass)
	}
	ww.PrivateHash=ww.GetHash(pass)
}

func (ww *WepWallet) Save() {
	var content bytes.Buffer
	wFile:=fmt.Sprintf(walletFile,ww.Address)
	gob.Register(elliptic.P256())
	encoder:=gob.NewEncoder(&content)
	err:=encoder.Encode(ww)
	Err(err,0)
	err=ioutil.WriteFile(wFile,content.Bytes(),0644)
	Err(err,0)
}

func (ww *WepWallet) Load() error {
	wFile:=fmt.Sprintf(walletFile,ww.Address)
	if _, err := os.Stat(wFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(wFile)
	Err(err,0)
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err=decoder.Decode(ww)
	Err(err,0)
	return nil
}

func (ww *WepWallet) GetHash(pass string) string {
	return setHash(strconv.Itoa(ww.Timestamp)+"|"+setHash(pass))
}


func (ww *WepWallet) Verify(pass string) bool {
	return ww.PrivateHash==ww.GetHash(pass)
}
