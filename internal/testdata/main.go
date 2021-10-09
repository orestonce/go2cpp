package testdata

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

func HelloWorld2(i int, j int32, name string) int32 {
	fmt.Println("Hello world: ", i, j, i+int(j), name)
	return int32(i) + j + int32(len(name))
}

func Ed25519(pk string, req2 HelloWorldReq) (pub string) {
	key := ed25519.NewKeyFromSeed([]byte(pk))
	return hex.EncodeToString(key.Public().(ed25519.PublicKey))
}

type HelloWorldL0 struct {
	Name string
	Age  byte
}

type HelloWorldReq struct {
	Name string
	Age  int
	U8   uint8
	I8   int8
	L0   HelloWorldL0
}

type HelloWorldResp struct {
	Data string
	U8   uint8
	I8   int8
	L0   HelloWorldL0
}

func HelloWorld(req HelloWorldReq) HelloWorldResp {
	fmt.Println("Hello, " + req.Name)
	return HelloWorldResp{
		Data: "data1 " + req.Name,
	}
}

func HelloWorld4(i byte) (o byte) {
	return i + 1
}

func HelloWorld5(i []byte) (o []string) {
	return nil
}

func HelloWorld6(i []HelloWorldReq) (o []HelloWorldResp) {
	return nil
}

func HelloWorld7(b bool) bool {
	return !b
}
