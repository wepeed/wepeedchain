package core

/*
type version struct {
    Version    int
    BestHeight int
    AddrFrom   string
}

type getblocks struct {
    AddrFrom string
}

type inv struct {
    AddrFrom string
    Type     string
    Items    [][]byte
}
type getdata struct {
    AddrFrom string
    Type     string
    ID       []byte
}


var nodeAddress string
var knownNodes = []string{"localhost:3000"}

func StartServer(nodeID, minerAddress string) {
    nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
    miningAddress = minerAddress
    ln, err := net.Listen(protocol, nodeAddress)
    defer ln.Close()
    bc:=NewBlockchain(nodeID)
    if nodeAddress != knownNodes[0] {
        sendVersion(knownNodes[0],bc)
    }
    for {
        conn, err := ln.Accept()
        go handleConnection(conn, bc)
    }
}

func sendVersion(addr string, bc *Blockchain) {
    bestHeight:=bc.GetBestHeight()
    payload:=gobEncode(version{nodeVersion,bestHeight,nodeAddress})
    request:= append(commandToBytes("version"), payload...)
    sendData(addr,request)
}

func commandToBytes(command string) []byte {
    var bytes [commandLength]byte
    for i, c := range command {
        bytes[i] = byte(c)
    }
    return bytes[:]
}

func bytesToCommand(bytes []byte) string {
    var command []byte
    for _, b := range bytes {
        if b != 0x0 {
            command = append(command, b)
        }
    }
    return fmt.Sprintf("%s", command)
}

func handleConnection(conn net.Conn, bc *Blockchain) {
    request, err := ioutil.ReadAll(conn)
    command := bytesToCommand(request[:commandLength])
    fmt.Printf("Received %s command\n", command)
    switch command {
    case "version":
        handleVersion(request,bc)
    default:
        fmt.Println("Unknown command!")
    }
    conn.Close()
}

func handleVersion(request []byte, bc *Blockchain) {
    var buff bytes.Buffer
    var payload version
    buff.Write(request[commandLength:])
    dec := gob.NewDecoder(&buff)
    err := dec.Decode(&payload)
    myBestHeight := bc.GetBestHeight()
    foreignerBestHeight := payload.BestHeight
    if myBestHeight<foreignerBestHeight {
        sendGetBlocks(payload.AddrFrom)
    } else if myBestHeight>foreignerBestHeight {
        sendVersion(payload.AddrFrom, bc)
    }
    if !nodeIsKnown(payload.AddrFrom) {
        knownNodes = append(knownNodes, payload.AddrFrom)
    }
}

func handleGetBlocks(request []byte, bc *Blockchain) {
    ...
    blocks := bc.GetBlockHashes()
    sendInv(payload.AddrFrom, "block", blocks)
}


func handleInv(request []byte, bc *Blockchain) {
    ...
    fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

    if payload.Type == "block" {
        blocksInTransit = payload.Items

        blockHash := payload.Items[0]
        sendGetData(payload.AddrFrom, "block", blockHash)

        newInTransit := [][]byte{}
        for _, b := range blocksInTransit {
            if bytes.Compare(b, blockHash) != 0 {
                newInTransit = append(newInTransit, b)
            }
        }
        blocksInTransit = newInTransit
    }

    if payload.Type == "tx" {
        txID := payload.Items[0]

        if mempool[hex.EncodeToString(txID)].ID == nil {
            sendGetData(payload.AddrFrom, "tx", txID)
        }
    }
}

func handleGetData(request []byte, bc *Blockchain) {
    ...
    if payload.Type == "block" {
        block, err := bc.GetBlock([]byte(payload.ID))

        sendBlock(payload.AddrFrom, &block)
    }

    if payload.Type == "tx" {
        txID := hex.EncodeToString(payload.ID)
        tx := mempool[txID]

        sendTx(payload.AddrFrom, &tx)
    }
}
*/


/*
import (
	"fmt"
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"encoding/gob"
)



func TcpDial(addr string) (*bufio.ReadWriter, error) {
	log.Print("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func NewNodepoint() *Nodepoint {
	nodeGet:=NodeRegister("node","get")
	log.Print("Node get : "+nodeGet)
	return &Nodepoint{
		handler: map[string]HandleFunc{},
	}
}

func (np *Nodepoint) AddHandleFunc(name string, f HandleFunc) {
	np.m.Lock()
	np.handler[name]=f
	np.m.Unlock()
}

func (np *Nodepoint) Listen() error {
	var err error
	np.listener,err=net.Listen("tcp",":"+strconv.Itoa(Configure.Port))
	if err != nil {
		return err
	}
	log.Println("Listen on", np.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := np.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go np.handleMessages(conn)
	}
}

func (np *Nodepoint) handleMessages(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()
	for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}
		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd + "'")
		np.m.RLock()
		handleCommand, ok := np.handler[cmd]
		np.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw)
	}
}

func NodeServer() error {
	Nodepoint:=NewNodepoint()
	Nodepoint.AddHandleFunc("NODECHECK", handleNodeCheck)
	Nodepoint.AddHandleFunc("SYNC", handleSync)
	Nodepoint.AddHandleFunc("VERSION", handleVersion)
	Nodepoint.AddHandleFunc("DOWNLOAD", handleDownload)
	Nodepoint.AddHandleFunc("TXPOOL", handleTxPool)
	Nodepoint.AddHandleFunc("STRING", handleString)
	Nodepoint.AddHandleFunc("DATA", handleData)
	return Nodepoint.Listen()
}

func NodesGet(nodeGet string) []NodeAddr {
	line := strings.Split(nodeGet, "||")
	for k:=range line {
		result := strings.Split(line[k], "|")
		nodes=append(nodes,NodeAddr{result[0],result[1],result[2]})
	}
	return nodes
}

func handleNodeCheck(rw *bufio.ReadWriter) {
	arr:=IpCheck()
	reader :=strings.NewReader("cmode=get&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&network="+Configure.Network+"&nettype="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&portnum="+strconv.Itoa(Configure.Port)+"&blocktime="+strconv.Itoa(Configure.Blocktime))
	request, _ := http.NewRequest("POST", "http://"+Configure.Mainserver+"/node/node_work.html", reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, _ := client.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	s:=string(body)
	r:=NodesGet(s)
	log.Println(r)
}

func handleTxPool(rw *bufio.ReadWriter) {
	arr:=IpCheck()
	reader :=strings.NewReader("cmode=txpool&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&network="+Configure.Network+"&nettype="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&portnum="+strconv.Itoa(Configure.Port)+"&blocktime="+strconv.Itoa(Configure.Blocktime))
	request, _ := http.NewRequest("POST", "http://"+Configure.Mainserver+"/pow", reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, _ := client.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	s:=string(body)
	//r:=NodeTx(s)

	log.Println(s)
}

func handleDownload(rw *bufio.ReadWriter) {
	log.Println("Download file")
	filename:=chainFileName()
    fileUrl:="http://"+Configure.Mainserver+"/node/download.html?file="+filename
    err := DownloadFile("."+filepathSeparator+Configure.Mainserver+filepathSeparator+filename, fileUrl)
    if err != nil {
		log.Println(err)
    }
}

func handleString(rw *bufio.ReadWriter) {
	log.Print(handelStr)
	err:=rw.Flush()
	if err!=nil {
		log.Println("Flush failed.", err)
	}
}

func handleData(rw *bufio.ReadWriter) {
	log.Print("Receive data:")
	dec:=gob.NewDecoder(rw)
	err:=dec.Decode(&handelData)
	if err!=nil {
		log.Println("Error decoding data:", err)
		return
	}
	log.Printf("Outer data: \n%#v\n", handelData)
}

func handleSync(rw *bufio.ReadWriter) {
	handleVersion(rw)
	handleDownload(rw) 
	cubechainInfo=ChainFileRead()
	handelStr="Download chain file"
	handleString(rw)
	ccnt:=ChainCheck(&cubechainInfo)
	handelStr="Chain check : "+strconv.Itoa(ccnt)
	handleString(rw)
}

func handleVersion(rw *bufio.ReadWriter) {
	log.Println("Version Informaition:"+Version)
	err:=rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

func DownloadFile(filepath string, url string) error {
    out,err:=os.Create(filepath)
    if err!=nil {
        return err
    }
    defer out.Close()
    resp,err:=http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    _, err=io.Copy(out, resp.Body)
    if err!=nil {
        return err
    }
    return nil
}

func NodeRegister(comm string,mode string) string {
	arr:=IpCheck()
	reader :=strings.NewReader("cmode="+mode+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&network="+Configure.Network+"&nettype="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&portnum="+strconv.Itoa(Configure.Port)+"&blocktime="+strconv.Itoa(Configure.Blocktime))
	request, _ := http.NewRequest("POST", "http://"+Configure.Mainserver+"/chain/"+comm, reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, err := client.Do(request)
    if err != nil {
        return ""
    }
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return ""
    }
	return string(body)
}


func NodeTx() {
	arr:=IpCheck()
	reader :=strings.NewReader("cmode=txpool&_token=9X1rK2Z2sofIeFpqg6VBXI5aUWsPOfGPGyzzztgu&blockno="+strconv.Itoa(blockno)+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&network="+Configure.Network+"&nettype="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&portnum="+strconv.Itoa(Configure.Port)+"&blocktime="+strconv.Itoa(Configure.Blocktime))
	request, _ := http.NewRequest("POST", "http://node.cubescan.io/pow", reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, err := client.Do(request)
    if err != nil {
        return TxpoolToBst("")
    }
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return TxpoolToBst("")
    }
	s:=string(body)

	r:=TxpoolToBst(s)
	return r
}

func NodeMiningReturn(b *Block) string {
	arr:=IpCheck()
	blocksize:=BlockSize(b.Index,b.Cubeno)
	
	reader:=strings.NewReader("cmode=mining&_token=9X1rK2Z2sofIeFpqg5VBXI5aUWsPOfGPGyzzztgu&txmine="+TxMine+"&txcnt="+strconv.Itoa(Txcnt)+"&txamount="+strconv.FormatFloat(Txamount,'f',-1,64)+"&merkle="+b.Merkle+"&amount="+strconv.FormatFloat(Pratio.BlockHash+Sumfee,'f',-1,64)+"&fee="+strconv.FormatFloat(Sumfee,'f',-1,64)+"&hashcnt="+strconv.Itoa(b.Nonce)+"&phash="+b.PrevCubeHash+"&hash="+b.Hash+"&maddress="+Configure.Address+"&blocksize="+strconv.FormatInt(blocksize,10)+"&cubeno="+strconv.Itoa(b.Index)+"&blockno="+strconv.Itoa(b.Cubeno)+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&network="+Configure.Network+"&nettype="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&portnum="+strconv.Itoa(Configure.Port)+"&blocktime="+strconv.Itoa(Configure.Blocktime))
	//echo(TxMine)
	request,_:=http.NewRequest("POST","http://node.cubescan.io/pow",reader)
	request.Header.Add("content-type","application/x-www-form-urlencoded")
	request.Header.Add("cache-control","no-cache")
	client := &http.Client{}
	res, err := client.Do(request)
    if err != nil {
        return ""
    }
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return ""
    }
	s:=string(body)
	//echo (s)
	return s
}






func NodeBlockTree(blockno int) []byte {
	txp:=Serialize(NodeTxBST(blockno))
	//echo ("----------------")
	//echo (txp)
	//echo ("----------------")
	//NodeMining(blockno)
	return txp
}


func SearchNode() bool{
	rs:=false
	ii:=10
	for i:=1;i<=ii;i++ {
		nodeGet:=NodeRegister("node","get")	
		if nodeGet=="Nothing Node." {
			fmt.Printf("[ %s] Not found node.\n", time.Now())
			time.Sleep(30*time.Second)
		} else {
			line := strings.Split(nodeGet, "||")
			for k:=range line {
				result := strings.Split(line[k], "|")
				nodes=append(nodes,NodeAddr{result[0],result[1],result[2]})
				rs=true
			}
			i=ii
		}
	}
	return rs
}

func NodeListening() string {
	cb:=make(chan string)
	sb:="Node Listening start"
	arr:=IpCheck()
	myAddr:=arr[1]+":"+strconv.Itoa(Configure.Port)
	addr, err := net.ResolveTCPAddr("tcp4",myAddr)
	listener, err := net.ListenTCP("tcp",addr)
	netError(err)
	go func(l *net.TCPListener) {
		for {
			connection, err := l.AcceptTCP()
			netError(err)
			sb="connection"+strconv.Itoa(Configure.Port)+strconv.Itoa(int(time.Now().Unix()))
			cb <-sb
			fmt.Println(connection)
		}
	}(listener)
	return ""
}
*/