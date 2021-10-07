package go2cpp

import (
	"github.com/orestonce/go2cpp/testdata"
	"testing"
)

func TestGo2cppContext_Generate1(t *testing.T) {
	//err := os.MkdirAll("testdata", 0777)
	//if err != nil {
	//	panic(err)
	//}
	//err = os.WriteFile("testdata/main.go", []byte(str), 0777)
	//if err != nil {
	//	panic(err)
	//}
	ctx := NewGo2cppContext(NewGo2cppContextReq{
		CppBaseName: "InProcessRpc",
		GoLibName:   "InProcessRpc-impl",
	})
	ctx.Generate1(testdata.HelloWorld2)
	ctx.Generate1(testdata.HelloWorld)
	ctx.Generate1(testdata.Ed25519)
	ctx.Generate1(testdata.HelloWorld4)
	ctx.Generate1(testdata.HelloWorld5)
	ctx.Generate1(testdata.HelloWorld6)
	ctx.MustCreate386LibraryInDir("tmp/temp1")
}

const str = `package testdata

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
}`
