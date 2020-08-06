# windows PC端wxapkg解密

## 说明
由于不想安装安卓模拟器去提取wxapkg包，windows PC端的微信也支持小程序，但是PC端的wxapkg是被加密存储的。该项目是把wxapkg解密。目前微信PC版本为：**2.9.5.31**.
## 使用方法
`pc_wxapkg_decrypt.exe -wxid 微信小程序id -in 要解密的wxapkg路径 -out 解密后的路径`

```
pc_wxapkg_decrypt.exe -h

Usage of pc_wxapkg_decrypt.exe:
   -in string
         需要解密的wxapkg的文件路径 (default "__APP__.wxapkg")
   -iv string
         AES加密的IV,默认不需要设置，如果版本有变化，设置 (default "the iv: 16 bytes")
   -out string
         解密后的wxapkg的文件路径 (default "dec.wxapkg")
   -salt string
         pbkdf2用到的salt,默认不需要设置 (default "saltiest")
   -wxid string
         小程序的id
```

wxapkg路径为：C:\Users\xxxx\Documents\WeChat Files\Applet\\**wx2xxx84w9w7a3xxxx**\\\__APP__.wxapkg,小程序id为：wx2xxx84w9w7a3xxxx

解密完成后，就可以用wxappUnpacker（https://github.com/gudqs7/wxappUnpacker）解包了。
## 原理
PC端微信把wxapkg给加密，加密后的文件的起始为**V1MMWX**。

加密方法为：
1. 首先pbkdf2生成AES的key。利用微信小程序id字符串为pass，salt为**saltiest** 迭代次数为1000。调用pbkdf2生成一个32位的key
2. 首先取原始的wxapkg的包得前1023个字节通过AES通过1生成的key和iv(**the iv: 16 bytes**),进行加密
3. 接着利用微信小程序id字符串的倒数第2个字符为xor key，依次异或1023字节后的所有数据，如果微信小程序id小于2位，则xorkey 为 **0x66**
4. 把AES加密后的数据（1024字节）和xor后的数据一起写入文件，并在文件头部添加**V1MMWX**标识
