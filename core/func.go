package core

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/gob"
	"runtime"
	"fmt"
	"os"
	"strconv"
	"strings"
    "bytes"
	"net"
	"net/http"
	"io"
	"io/ioutil"
    "path/filepath"
)

var echo=fmt.Println
var Configure Configuration
var filepathSeparator=string(filepath.Separator)
var Version="1.1"
var MblockDelim="||"
var TxDelim="|"
var PotDelim="|"


func Err(err error, exit int) int {
	if err != nil {
		fmt.Println(err)	
	}
	if exit>=1 {
		os.Exit(exit)
		return 1
	}
	return 0
}

func netError(err error) {
	if err!=nil && err!=io.EOF {
		fmt.Println("Network Error : ", err)
	}
}

func IpCheck() []string {
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil
	}
	addrs=append(addrs,host)
	if len(addrs)==2 {
		addrs2:=make([]string,3)
		addrs2[0]="mac_linux"
		addrs2[1]=addrs[0]
		addrs2[2]=addrs[1]
		addrs=addrs2
	}
	return addrs
}

func FilePath(idx int) string {
	divn:=idx/Configure.Datanumber
	divm:=idx%Configure.Datanumber
	if divm>0 {
		divn++
	} else if divm==0 {
		divm=Configure.Datanumber
	}
	if divn==0 {
		divn++
		divm=1
	}
	//nhex:=fmt.Sprintf("%x",Configure.Datanumber)
	nstr:=fmt.Sprintf("%0.5x",divn)
	dirname:=Configure.Datafolder+filepathSeparator+nstr
	if dirExist(dirname)==false {
		if err:=os.MkdirAll(dirname, os.FileMode(0755)); err!=nil {
			return "Directory not found.\\1"
		}	
	}	
	return dirname
}

func FileWrite(path string, object interface{}) error {
	file,err:=os.Create(path)
	if err==nil {
		encoder:=gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func FileRead(path string, object interface{}) error {
	file,err:=os.Open(path)
	if err==nil {
		decoder:=gob.NewDecoder(file)
		err=decoder.Decode(object)
	}
	file.Close()
	return err
}

func FileSize(dirpath string) int64 {
	file, err := os.Open(dirpath) 
	if err != nil {
		echo (err)
	}
	fi, err := file.Stat()
	if err != nil {
		echo (err)
	}
	return fi.Size()
}

func fileCheck(e error) {
	if e!=nil {
		_, file, line, _:=runtime.Caller(1)
		fmt.Println(line,  "\t", file, "\n", e)
		os.Exit(1)
	}
}

func fileSearch(dirname string,find string) string{
    result:=""
	d,err:=os.Open(dirname)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
    file, err:=d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi:=range file {
        if fi.Mode().IsRegular() {
			fstr:=fi.Name()
			if strings.Index(fstr,find)>=0 {
				result=fi.Name()
				return result
			}
        }
    }
	return result
}

func dirExist(dirName string) bool{
	result:=true
	_,err:=os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err ) {
			result=false
		}
	}
	return result
}

func MaxFind(dirpath string) string {
	find:="0"
    d, err:=os.Open(dirpath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer d.Close()
	fi, err:=d.Readdir(-1)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, fi:=range fi {
        if fi.Mode().IsRegular() {
        } else {
  			if fi.Name()>find {
				find=fi.Name()
			}
		}
   }
   return find
}


func PathDelete(path string) error {
	err:=os.RemoveAll(path)
	os.MkdirAll(path,0755)
	return err
}

func GetBytes(key interface{}) ([]byte) {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
        return nil
    }
    return buf.Bytes()
}

func ByteToStr(bytes []byte) string {
	var str []byte
	for _, v := range bytes {
		if v != 0x0 {
			str = append(str, v)
		}
	}
	return fmt.Sprintf("%s", str)
}


func setHashV(str interface{}) string {
	h:=sha256.New()
	h.Write(GetBytes(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash(str string) string {
	h:=sha256.New()
	h.Write([]byte(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func setHash2(str string) string {
	h:=sha512.New384()
	h.Write([]byte(str))
	hashed:=h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func CallHash(str string,hno int) string {
	switch hno {
		case 2 : return setHash2(str)
		default : return setHash(str)
	}
}



func CurrentHeight() int {
	result:=0
	f:=MaxFind(Configure.Datafolder+filepathSeparator)
	if f=="0" {
		return 1
	}
	nint,_:=strconv.ParseUint(f,16,32)
	result=(int(nint)-1)*Configure.Datanumber
	dirres:=result
	path:=FilePath(result)
	for i:=0;i<2000;i++ {
		findi:=strconv.Itoa(dirres+i)
		if fileSearch(path,findi+"_")>"" {
			result=dirres+i
		}
	}
	return result	
}


func GetTxCount(addr string) int {
	c,count:=0,0
	var mblock MicroBlock
	if c<=0 {
		c=CurrentHeight()-1
	}
	for i:=0;i<c;i++ {
		mblock.Index=i
		err:=mblock.Read()
		Err(err,0)
		if(mblock.Data.From==addr) {
			count++
		}
	}
	return count
}


func NodeSend(cmode string,data string) string {
	arr:=IpCheck()
	reader:=strings.NewReader("cmode="+cmode+"&_token=9X1rK2Z2sofIeFpqg6VBXI5aUWsPOfGPGyzzztgu&data="+data+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&netname="+Configure.Network+"&netset="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&netport="+strconv.Itoa(Configure.Port)+"&ver="+Version)
	request,_:=http.NewRequest("POST","http://"+Configure.Mainserver+"/"+cmode, reader)
	request.Header.Add("content-type","application/x-www-form-urlencoded")
	request.Header.Add("cache-control","no-cache")
	client:=&http.Client{}
	res, err := client.Do(request)
	Err(err,0)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	Err(err,0)
	s:=string(body)
	return s
}

