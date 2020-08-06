package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha1"
    "flag"
    "fmt"
    "golang.org/x/crypto/pbkdf2"
    "io/ioutil"
    "log"
)


var (
    wxid string
    iv string
    salt string
    wxapkgPath string
    decWxapkgPath string
)

func main(){
    flag.StringVar(&wxid,"wxid","","小程序的id")
    flag.StringVar(&iv,"iv","the iv: 16 bytes","AES加密的IV,默认不需要设置，如果版本有变化，设置")
    flag.StringVar(&salt,"salt","saltiest","pbkdf2用到的salt,默认不需要设置")
    flag.StringVar(&wxapkgPath,"in","__APP__.wxapkg","需要解密的wxapkg的文件路径")
    flag.StringVar(&decWxapkgPath,"out","dec.wxapkg","解密后的wxapkg的文件路径")
    flag.Parse()

    if wxid == ""{
        fmt.Println("缺少wx小程序id(wxid)，该id在pc存放wxapkg包的路径上寻找")
        return
    }
    dec()
}

func dec(){
    dataByte,err := ioutil.ReadFile(wxapkgPath)
    if err != nil{
        log.Fatal(err)
    }

    dk := pbkdf2.Key([]byte(wxid),[]byte(salt),1000,32,sha1.New)
    block,_ := aes.NewCipher(dk)
    blockMode := cipher.NewCBCDecrypter(block,[]byte(iv))
    originData := make([]byte,1024)
    blockMode.CryptBlocks(originData,dataByte[6:1024+6])

    afData := make([]byte,len(dataByte) - 1024 - 6)
    var xorKey = byte(0x66)
    if len(wxid) >= 2 {
        xorKey = wxid[len(wxid) - 2]
    }
    for i,b := range dataByte[1024+6:]{
        afData[i] = b ^ xorKey
    }

    originData = append(originData[:1023],afData...)

    err = ioutil.WriteFile(decWxapkgPath,originData,0666)
    if err != nil{
        fmt.Println("写文件失败")
        return
    }
    fmt.Println("解密成功")
}



