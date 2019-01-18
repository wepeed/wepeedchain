package main

import (
	"fmt"
	"time"
	"strconv"
	 "./core"
)
var echo=fmt.Println
var Configure core.Configuration
var mstr="miningtesting!!..."
var td core.TxData
var mbk core.MicroBlock
var pt core.POT

func init() {
	Configure=core.LoadConfiguration("./core/wepeed.conf")
	core.Configure=Configure
}

func main() {
    echo("start")
	WalletTest1()
    echo("end")
}


func main3() {
    echo("start")
	ApiTest1()
    echo("end")
}

func main1() {
    echo("start")
	PotTest6()
    echo("end")
}

func main2() {
    echo("start")
	TxTest()

    echo()
    echo()
	FileTest()

    echo()
    echo()
	MBlockTest()

    echo("end")
}

func TxTest() {
	str:="TX|0x6c01c09e1da9164439e73e3266a582a27942dc4d|0x39f7db90fc811b5c3539d572f4dacfd1676ecde4|3010|0.02|SSDDF3u9445453487u8q83wruf9ujiew938"
	td.Input(str)
	td.Print()
	txdata:=td.Transaction()
	echo ("------------------")
	echo (txdata)
	echo ("------------------")
}

func FileTest() {
	filepath:=core.FilePath(2)
	echo (filepath)
}

func MBlockTest() {
	//str:="1|1566678019|1355|039093uXXEEu32843050943i90i9fidsifjids"
	mbk.Input(td)
	mbk.Print()
	//mbk.Save()
}

func PotTest() {
	pt.SetAddress("Cwejrowr8383j48291212");
	pt.Start()
}

func PotTest2() {
	pt.SetAddress("Cwejrowr8383j48291212");
	pt.MidTime()
}

func PotTest3() {
	pt.SetAddress("Cwejrowr8383j48291212");
	pt.End()
}

func PotTest4() {
	pt.SetAddress("Cwejrowr8383j48291212");
	pt.Start()
	time.Sleep(100 * time.Second)
	pt.MidTime()
}

func PotTest5() {
	pt.SetAddress("Cwejrowr8383j48291212");
	pt.Start()
	time.Sleep(100 * time.Second)
	pt.MidTime()
	time.Sleep(70 * time.Second)
	pt.End()
}

func PotTest6() {
	pt.SetAddress("Cwejrowr8383j48291212");
	pt.Start()
	time.Sleep(3 * time.Second)
	pt.MidTime()
	time.Sleep(4 * time.Second)
	pt.End()
}

func ApiTest1() {
	addr:="0x39f7db90fc811b5c3539d572f4dacfd1676ecde4"
	gb:=core.GetBalance(addr)
	gtc1,gtc2:=core.GetTransactionCount(addr)
	gtl:=core.GetTransactionList(addr)

	echo (gb)
	echo (gtc1)
	echo (gtc2)
	echo (gtl)
}

func WalletTest1() {
	addr:="WaatAor45QpXGURz348Nr1WXcwby7vHtqo"
	var wep core.WepWallet
	wep.Address=addr
	wep.Load()
	t:=wep.Verify("1234")
	echo (wep)
	echo (t)
}

func WalletTest2() {
	pass:="1234"
	var wep core.WepWallet
	wep.Create(pass)

	echo (wep)
}

func WalletTest3() {
	pass:="1234"
	var wep core.WepWallet
	for i:=1;i<=100;i++ {
		wep.Create(pass)
		echo (strconv.Itoa(i)+" "+wep.Address)
	}
}

func WalletTest4() {
	pass:="1234"
	var wep core.WepWallet
	for i:=1;i<=10;i++ {
		wep.Create(pass)
		wep.Save()
	}
}

func WalletTest5() {
	addr:="WaeHZc4uzjuAieiBidXgCzmsL7yLebzeWN"
	var wep core.WepWallet
	wep.Address=addr
	wep.Load()
	echo (wep)
}
