package core

import (
	"strconv"
	"time"
)

type POT struct {
	WalletAddress	string
	StartTimestamp	int
	MidTimestamp	int
	EndTimestamp	int
	StartHash		string
	EndHash			string
}

func (pt *POT) SetAddress(address string) {
	pt.WalletAddress=address
}

func (pt *POT) String() string {
	str:=pt.WalletAddress+PotDelim+strconv.Itoa(pt.StartTimestamp)+PotDelim+strconv.Itoa(pt.MidTimestamp)+PotDelim+strconv.Itoa(pt.EndTimestamp)+PotDelim+pt.StartHash+PotDelim+pt.EndHash
	return str
}

func (pt *POT) Start() {
	pt.StartTimestamp=int(time.Now().Unix())
	result:=pt.SendPot()
	pt.StartHash=result
	echo (result)
}

func (pt *POT) MidTime() {
	pt.MidTimestamp=int(time.Now().Unix())
	result:=pt.SendPot()
	echo (result)
}

func (pt *POT) End() {
	pt.EndTimestamp=int(time.Now().Unix())
	pt.SetHash()
	result:=pt.SendPot()
	echo (result)
}

func (pt *POT) SetHash() {
	pt.EndHash=pt.GetHash()
}

func (pt *POT) GetHash() string {
	return setHash(pt.String())
}

func (pt *POT) SendPot() string {
	result:=NodeSend("pot",pt.String())
	return result
}


/*

POT
1. start time, address -> TSN (starthash값 부여)
2. 1시간 단위로 middle time, address, starthash
3. end time, address, key값 및 hash값(starthash+endtime) -> TSN


TSN
1. start time 일때 StartHash값 부여
2. middle time 일때 StartHash값 검증후 등록
3. end time StartHash값 및 EndHash값 endtime 검증

*/