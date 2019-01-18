package core

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

type MicroBlock struct {
	Index			int
	Timestamp		int
	Data			TxData
	Hash			string
	PrevHash		string
	Nonce			int
}

func (mblock *MicroBlock) Input(td TxData) {
	bstr:=td.Transaction()
	if len(bstr)<20 {
		return
	}
	result:=strings.Split(bstr,"|")
	if result[1]=="" {
		return
	}
	mblock.Index,_=strconv.Atoi(result[0])
	mblock.Timestamp,_=strconv.Atoi(result[1])
	mblock.Nonce,_=strconv.Atoi(result[2])
	mblock.Data=td
	mblock.PrevHash=result[3]
	mblock.SetHash()

	res:=mblock.ToStamp()
	if res=="Success." {
		mblock.Save()
	}
}

func (mblock *MicroBlock) InputString(Str string) {
	var td TxData
	td.Input(Str)
	mblock.Input(td)
}

func (mblock *MicroBlock) ToStamp() string {
	result:=NodeSend("mblock",mblock.String())
	return result
}

func (mblock *MicroBlock) Save() error {
	filename:=mblock.FileName()
	datapath:=mblock.FilePath()
	err:=FileWrite(datapath+string(filepathSeparator)+filename,mblock)
	return err
}

func (mblock *MicroBlock) Read() error {
	filename:=mblock.FileName()
	datapath:=mblock.FilePath()
	err:=FileRead(datapath+string(filepathSeparator)+filename,mblock)
	return err
}

func (mblock *MicroBlock) HashString() string {
	toStr:=strconv.Itoa(mblock.Index)+MblockDelim+strconv.Itoa(mblock.Timestamp)+MblockDelim+mblock.Data.String()+MblockDelim+mblock.PrevHash+MblockDelim+strconv.Itoa(mblock.Nonce)
	return toStr
}

func (mblock *MicroBlock) String() string {
	toStr:=strconv.Itoa(mblock.Index)+MblockDelim+strconv.Itoa(mblock.Timestamp)+MblockDelim+mblock.Data.String()+MblockDelim+mblock.Hash+MblockDelim+mblock.PrevHash+MblockDelim+strconv.Itoa(mblock.Nonce)
	return toStr
}

func (mblock *MicroBlock) SetHash() {
	mblock.Hash=mblock.GetHash()
}

func (mblock *MicroBlock) GetHash() string {
	hashstr:=mblock.HashString()
	h:=sha256.New()
	h.Write([]byte(hashstr))
	return hex.EncodeToString(h.Sum(nil))
}

func (mblock *MicroBlock) Verify() bool {
	pidx:=mblock.Index-1
	if pidx>0 {
		pblock:=&MicroBlock{}
		pblock.Read()
		if pblock.Hash!=mblock.PrevHash {
			return false
		} else if mblock.Hash!=mblock.GetHash() {
			return false
		}
	}
	return true
}

func (mblock *MicroBlock) FileName() string {
	filename:=""
	if mblock.Hash=="" {
		filename=fileSearch(mblock.FilePath(),strconv.Itoa(mblock.Index)+"_")
	} else {
		filename=strconv.Itoa(mblock.Index) + "_" + mblock.Hash + ".mbk"
	}
	return filename	
}

func (mblock *MicroBlock) FilePath() string {
	filepath:=FilePath(mblock.Index)
	return filepath	
}

func (mblock *MicroBlock) FileSize() int64 {
	filepath:=mblock.FilePath()
	filesize:=FileSize(filepath)
	return filesize	
}

func (mblock *MicroBlock) Balance(addr string) float64 {
	result:=0.0
	if mblock.Data.From==addr {
		result+=(mblock.Data.Amount+mblock.Data.Fee)*(-1)
	}
	if mblock.Data.To==addr {
		result+=mblock.Data.Amount
	}
	return result
}

func (mblock *MicroBlock) TxCount(addr string) (int,int) {
	send:=0
	receive:=0
	if mblock.Data.From==addr {
		send++
	}
	if mblock.Data.To==addr {
		receive++
	}
	return send,receive
}

func (mblock *MicroBlock) TxList(addr string) (TxData,bool) {
	var txd TxData
	result:=false
	if mblock.Data.From==addr || mblock.Data.To==addr {
		result=true
		txd=mblock.Data
	}
	return txd,result
}

func (mblock *MicroBlock) TxDetail(hash string) (TxData,bool) {
	var txd TxData
	result:=false
	if mblock.Data.Hash==hash {
		result=true
		txd=mblock.Data
	}
	return txd,result
}

func (mblock *MicroBlock) Address(addr string) string {
	return mblock.Data.From
}

func (mblock *MicroBlock) Print() {
	echo ("==============Block Head==============")
	echo ("Index =",mblock.Index)
	echo ("Timestamp =",mblock.Timestamp)
	echo ("Hash =",mblock.Hash)
	echo ("PrevHash =",mblock.PrevHash)
	echo ("Nonce =",mblock.Nonce)
	echo ("==============Block Body==============")
	mblock.Data.Print()
}

func (mblock *MicroBlock) PrintInfo() {
	echo ("==============File Info==============")
	echo ("FileName=",mblock.FileName())
	echo ("FilePath =",mblock.FilePath())
	echo ("FileSize =",mblock.FileSize())
}
