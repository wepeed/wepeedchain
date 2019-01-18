package wallet

import (
    "fmt"
	"strings"
	"strconv"
)
 
const (
	encodeBase="abdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	encodeKey="c"
	encodeCnt=61
)

func Cencode(num int) string {
	result:=""
	if num>100 {
		result="#"+encodeBase[num-100:num-99];
	} else if num>=encodeCnt {
		rnum:=num-encodeCnt
		result="c"+encodeBase[rnum:rnum+1]
	} else {
		result=encodeBase[num:num+1]
	}
	return result
}

func Cdecode(str string) string {
	r:=strings.Index(encodeBase,str)
	result:=fmt.Sprintf("%d",r)
	if(r<10) {
		result="0"+result
	}
	return result
}

func CencodeKey(num string) string {
	result:=""
	lnum:=len(num)
	ncount:=(lnum+1)/2
	for i:=0;i<ncount;i++ {
		d:=num[0:1]
		if i*2+2>lnum {
			d=num[i*2:i*2+1]
		} else {
			d=num[i*2:i*2+2]
		}
		n,_:=strconv.Atoi(d)

		if len(d)==1 {
			n+=100
		}
		p:=Cencode(n)
		//echo (fmt.Sprintf("%d  --->  %s",n,p))
		result+=p
	}
	return result
}

func CdecodeKey(str string) string {
	result:=""
	tmpshp:=0
	ncount:=len(str);
	for i:=0;i<ncount;i++ {
		d:=str[i:i+1]

		if d=="#" {
			tmpshp=100
			continue
		} else if d=="c" {
			tmpshp=encodeCnt
			continue
		}
		p:=Cdecode(d)


		if tmpshp>0 {
			n,_:=strconv.Atoi(p)
			if tmpshp==encodeCnt { n+=tmpshp }
			result+=strconv.Itoa(n)
		} else {
			result+=p
		}
		tmpshp=0;
	}
	return result
}

