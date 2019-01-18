package core

import (
)

func GetBalance(addr string) float64 {
	var mblock MicroBlock
	c:=CurrentHeight()
	result:=0.0
	for i:=1;i<c;i++ {
		mblock.Index=i
		mblock.Read()
		result+=mblock.Balance(addr)
	}
	return result
}

func GetTransactionCount(addr string) (int,int) {
	var mblock MicroBlock
	c:=CurrentHeight()
	result1,result2:=0,0
	for i:=1;i<c;i++ {
		mblock.Index=i
		mblock.Read()
		tmp1,tmp2:=mblock.TxCount(addr)
		result1+=tmp1
		result2+=tmp2
	}
	return result1,result2
}

func GetTransactionList(addr string) []TxData {
	var mblock MicroBlock
	var result []TxData
	c:=CurrentHeight()
	for i:=1;i<c;i++ {
		mblock.Index=i
		mblock.Read()
		tmptd,ok:=mblock.TxList(addr)
		if ok {
			result=append(result,tmptd)
		}
	}
	return result	
}

func GetTransactionDetail(txhash string) TxData {
	var mblock MicroBlock
	var result TxData
	c:=CurrentHeight()
	for i:=1;i<c;i++ {
		mblock.Index=i
		mblock.Read()
		tmptd,ok:=mblock.TxDetail(txhash)
		if ok {
			result=tmptd
			break
		}
	}
	return result	
}
