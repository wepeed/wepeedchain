package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/gob"
	"bytes"
	"strconv"
	"time"
	"strings"

)


// TxType Symbol Message

type TxData struct {
	Timestamp	int
	From		string
	To			string
	Amount		float64
	Fee			float64
	Hash		string
	Sign		string
	Nonce		int
}

func (td *TxData) Set() {
	// Input
	// Verify
}



func (td *TxData) Input(TxStr string) {
	if len(TxStr)<20 {
		return
	}
	result:=strings.Split(TxStr, "|")
	td.Timestamp=int(time.Now().Unix())
	td.From=result[1]
	td.To=result[2]
	td.Amount,_=strconv.ParseFloat(result[3],64)
	td.Fee,_=strconv.ParseFloat(result[4],64)
	td.Sign=result[5]
	td.Nonce=GetTxCount(td.From)
	td.SetHash()
}

func (td *TxData) HashString() string {
	toStr:=strconv.Itoa(td.Timestamp)+TxDelim+td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+strconv.Itoa(td.Nonce)
	return toStr
}

func (td *TxData) String() string {
	toStr:=strconv.Itoa(td.Timestamp)+TxDelim+td.From+TxDelim+td.To+TxDelim+strconv.FormatFloat(td.Amount,'f',-1,64)+TxDelim+strconv.FormatFloat(td.Fee,'f',-1,64)+TxDelim+td.Hash+TxDelim+td.Sign+TxDelim+strconv.Itoa(td.Nonce)
	return toStr
}


func (td *TxData) SetHash() {
	td.Hash=td.GetHash()
}

func (td *TxData) GetHash() string {
	hashstr:=td.HashString()
	h:=sha256.New()
	h.Write([]byte(hashstr))
	return hex.EncodeToString(h.Sum(nil))
}


func (td *TxData) Verify() bool {
	return td.Hash==td.GetHash()
}


func (td *TxData) ToByte() []byte {
	var buff bytes.Buffer
	enc:=gob.NewEncoder(&buff)
	err:=enc.Encode(td)
	if err!=nil {
		echo (err)
	}
	return buff.Bytes()
}

func (td *TxData) GetCount() int {
	return GetTxCount(td.From)
}

func (td *TxData) Print() {
	echo ("Timestamp =",td.Timestamp)
	echo ("From =",td.From)
	echo ("To =",td.To)
	echo ("Amount =",td.Amount)
	echo ("Fee =",td.Fee)
	echo ("Hash =",td.Hash)
	echo ("Sign =",td.Sign)
	echo ("Nonce =",td.Nonce)
}


func (td *TxData) Transaction() string {
	result:=NodeSend("tx",td.String())
	return result
}



