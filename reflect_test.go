package go2cpp

import (
	"github.com/orestonce/go2cpp/internal/testdata"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGo2cppContext_Generate1(t *testing.T) {
	//err := os.MkdirAll("internal/testdata", 0777)
	//if err != nil {
	//	panic(err)
	//}
	//err = os.WriteFile("internal/testdata/main.go", []byte(str), 0777)
	//if err != nil {
	//	panic(err)
	//}
	ctx := NewGo2cppContext(NewGo2cppContextReq{
		CppBaseName: "InProcessRpc",
		GoLibName:   "InProcessRpc-impl",
	})
	ctx.Generate1(testdata.Hello_BoolTrue)
	ctx.Generate1(testdata.Hello_BoolFalse)

	ctx.Generate1(testdata.Hello_Int8Max)
	ctx.Generate1(testdata.Hello_Int8Min)
	ctx.Generate1(testdata.Hello_Int8Common)

	ctx.Generate1(testdata.Hello_Uint8Max)
	ctx.Generate1(testdata.Hello_Uint8Min)
	ctx.Generate1(testdata.Hello_Uint8Common)

	ctx.Generate1(testdata.Hello_Int32Max)
	ctx.Generate1(testdata.Hello_Int32Min)
	ctx.Generate1(testdata.Hello_Int32Common)

	ctx.Generate1(testdata.Hello_Uint32Max)
	ctx.Generate1(testdata.Hello_Uint32Min)
	ctx.Generate1(testdata.Hello_Uint32Common)

	ctx.Generate1(testdata.Hello_IntMax)
	ctx.Generate1(testdata.Hello_IntMin)
	ctx.Generate1(testdata.Hello_IntCommon)

	ctx.Generate1(testdata.Hello_StringEmpty)
	ctx.Generate1(testdata.Hello_StringCommon0)
	ctx.Generate1(testdata.Hello_StringCommon1)

	ctx.Generate1(testdata.Hello_Slice0)

	ctx.Generate1(testdata.Hello_Struct0)

	ctx.MustCreate386LibraryInDir("tmp/temp1")
	
	err := os.WriteFile("tmp/temp1/main.cpp", []byte(mainCpp), 0777)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("g++", "InProcessRpc.cpp", "main.cpp", "-std=c++11", "-m32", "InProcessRpc-impl.a")
	cmd.Dir = "tmp/temp1/"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd = exec.Command(filepath.Join(pwd, "tmp/temp1/a.exe"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

const mainCpp = `
#include "InProcessRpc.h"
#include <iostream>

void assert_ok(bool ok, std::string tag)
{
	if(ok)
	{
		//std::cout << "pass: " << tag << std::endl;
		return;
	}
	std::cout << "assert_ok failed: " << tag << std::endl;
}

void testHello_StringXXX();
void testHello_Slice0();
void testHello_Struct0();

int main()
{
	std::cout << "begin go2cpp test:" << std::endl;
	assert_ok(Hello_BoolTrue(true) == true, "Hello_BoolTrue");
	assert_ok(Hello_BoolFalse(false) == false, "Hello_BoolFalse");

	assert_ok(Hello_Int8Max(127) == 127, "Hello_Int8Max");
	assert_ok(Hello_Int8Min(-128) == -128, "Hello_Int8Min");
	assert_ok(Hello_Int8Common(12) == 12, "Hello_Int8Common");

	assert_ok(Hello_Uint8Max(uint8_t(255)) == uint8_t(255), "Hello_Uint8Max");
	assert_ok(Hello_Uint8Min(uint8_t(0)) == uint8_t(0), "Hello_Uint8Min");
	assert_ok(Hello_Uint8Common(uint8_t(95)) == uint8_t(95), "Hello_Uint8Common");

	assert_ok(Hello_Int32Max(int32_t(2147483647)) == int32_t(2147483647), "Hello_Int32Max");
	assert_ok(Hello_Int32Min(int32_t(-2147483648)) == int32_t(-2147483648), "Hello_Int32Min");
	assert_ok(Hello_Int32Common(int32_t(10086)) == int32_t(10086), "Hello_Int32Common");

	assert_ok(Hello_Uint32Max(uint32_t(4294967295)) == uint32_t(4294967295), "Hello_Uint32Max");
	assert_ok(Hello_Uint32Min(uint32_t(0)) == uint32_t(0), "Hello_Uint32Min");
	assert_ok(Hello_Uint32Common(uint32_t(1001011)) == uint32_t(1001011), "Hello_Uint32Common");
	
	assert_ok(Hello_IntMax(int(2147483647)) == int(2147483647), "Hello_IntMax");
	assert_ok(Hello_IntMin(int(-2147483648)) == int(-2147483648), "Hello_IntMin");
	assert_ok(Hello_IntCommon(int(0x12345678)) == int(0x12345678), "Hello_IntCommon");
	
	testHello_StringXXX();

	testHello_Slice0();
	
	testHello_Struct0();
	
	std::cout << "end go2cpp test." << std::endl;
}

void testHello_StringXXX()
{
	assert_ok(Hello_StringEmpty("") =="", "Hello_StringEmpty");
	std::string origin = std::string("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2a\x2b\x2c\x2d\x2e\x2f\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3a\x3b\x3c\x3d\x3e\x3f\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4a\x4b\x4c\x4d\x4e\x4f\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5a\x5b\x5c\x5d\x5e\x5f\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6a\x6b\x6c\x6d\x6e\x6f\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7a\x7b\x7c\x7d\x7e\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff", 256);
	std::string after = Hello_StringCommon0(origin);
	assert_ok(origin==after && origin.length() == 256,"Hello_StringCommon0");
	
	std::string origin1;
	for (int i=0; i <1024; i++)
	{
		origin1.append(origin);
	}
	std::string after1 = Hello_StringCommon1(origin1);
	assert_ok(origin1==after1 && origin1.length() == 256*1024,"Hello_StringCommon1");
}

void testHello_Slice0()
{
	std::vector<std::string> slice0;
	slice0.push_back("1");
	slice0.push_back("2");
	slice0.push_back("34567");
	
	std::vector<std::string> slice0_after = Hello_Slice0(slice0);
	bool slice0_ok = true;
	if (slice0_after.size() != slice0.size()) {
		slice0_ok = false;
	} else {
		for (uint32_t idx =0; idx<slice0.size(); idx++)
		{
			if (slice0[idx] != slice0_after[idx])
			{
				slice0_ok = false;
				break;
			}
		}
	}
	assert_ok(slice0_ok, "Hello_Slice0");
}

void testHello_Struct0()
{
	Hello_Struct0Req origin;
	origin.Name = "name0";
	origin.Age = 192;
	origin.U8 = 76;
	origin.I8 = -3;
	origin.L1.Name = "name1";
	origin.L1.Age = 9;
	origin.Int32Slice.push_back(8);
	origin.Int32Slice.push_back(3);
	origin.Int32Slice.push_back(-4);
	origin.Int32Slice.push_back(90);
	origin.Name2 = "Name2";

	Hello_Struct0Req after = Hello_Struct0(origin);
	bool ok= true;
	if (after.Name != "name0") { ok = false; }
	if (after.Age != 192) { ok = false; }
	if (after.U8 != 76) { ok = false; }
	if (after.I8 != -3) { ok = false; }
	if (after.L1.Name != "name1") { ok = false; }
	if (after.L1.Age != 9) { ok = false; }
	if (after.Int32Slice.size() != 4) { 
		ok = false;
	} else {
		if (after.Int32Slice[0] != 8) { ok = false; }
		if (after.Int32Slice[1] != 3) { ok = false; }
		if (after.Int32Slice[2] != -4) { ok = false; }
		if (after.Int32Slice[3] != 90) { ok = false; }
	}
	if (after.Name2 != "Name2") { ok = false; }
	
	assert_ok(ok, "Hello_Struct0");
}
`

const str = `package testdata

import (
	"reflect"
	"strings"
)

func Hello_BoolTrue(a bool) bool {
	if !a {
		panic("Hello_BoolTrue")
	}
	return a
}

func Hello_BoolFalse(a bool) bool {
	if a {
		panic("Hello_BoolFalse")
	}
	return a
}

func Hello_Int8Max(a int8) int8 {
	if a != 127 {
		panic("Hello_Int8Max")
	}
	return a
}

func Hello_Int8Min(a int8) int8 {
	if a != -128 {
		panic("Hello_Int8Min")
	}
	return a
}

func Hello_Int8Common(a int8) int8 {
	if a != 12 {
		panic("Hello_Int8Common")
	}
	return a
}

func Hello_Uint8Max(a uint8) uint8 {
	if a != 255 {
		panic("Hello_Uint8Max")
	}
	return a
}

func Hello_Uint8Min(a uint8) uint8 {
	if a != 0 {
		panic("Hello_Uint8Min")
	}
	return a
}

func Hello_Uint8Common(a uint8) uint8 {
	if a != 95 {
		panic("Hello_Uint8Common")
	}
	return a
}

func Hello_Int32Max(a int32) int32 {
	if a != 2147483647 {
		panic("Hello_Int32Max")
	}
	return a
}

func Hello_Int32Min(a int32) int32 {
	if a != -2147483648 {
		panic("Hello_Int32Min")
	}
	return a
}

func Hello_Int32Common(a int32) int32 {
	if a != 10086 {
		panic("Hello_Int32Common")
	}
	return a
}

func Hello_Uint32Max(a uint32) uint32 {
	if a != 4294967295 {
		panic("Hello_Uint32Max")
	}
	return a
}

func Hello_Uint32Min(a uint32) uint32 {
	if a != 0 {
		panic("Hello_Uint32Min")
	}
	return a
}

func Hello_Uint32Common(a uint32) uint32 {
	if a != 1001011 {
		panic("Hello_Uint32Common")
	}
	return a
}

func Hello_IntMax(a int) int {
	if a != 2147483647 {
		panic("Hello_IntMax")
	}
	return a
}

func Hello_IntMin(a int) int {
	if a != -2147483648 {
		panic("Hello_IntMin")
	}
	return a
}

func Hello_IntCommon(a int) int {
	if a != 0x12345678 {
		panic("Hello_IntCommon")
	}
	return a
}

func Hello_StringEmpty(s string) string {
	if s != "" {
		panic("Hello_StringEmpty")
	}
	return s
}

func Hello_StringCommon0(s string) string {
	if s != string([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f, 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf, 0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf, 0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef, 0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}) {
		panic("Hello_StringCommon0")
	}
	return s
}

func Hello_StringCommon1(s string) string {
	expect := strings.Repeat(string([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f, 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf, 0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf, 0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf, 0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef, 0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}), 1024)
	if s != expect {
		panic("Hello_StringCommon1")
	}
	return s
}

func Hello_Slice0(data []string) []string {
	if !reflect.DeepEqual(data, []string{"1", "2", "34567"}) {
		panic("Hello_Slice0")
	}
	return data
}

type Hello_Struct0ReqL1 struct {
	Name string
	Age  byte
}

type Hello_Struct0Req struct {
	Name       string
	Age        int
	U8         uint8
	I8         int8
	L1         Hello_Struct0ReqL1
	Int32Slice []int32
	Name2      string
}

func Hello_Struct0(req Hello_Struct0Req) Hello_Struct0Req {
	expect := Hello_Struct0Req{
		Name: "name0",
		Age:  192,
		U8:   76,
		I8:   -3,
		L1: Hello_Struct0ReqL1{
			Name: "name1",
			Age:  9,
		},
		Int32Slice: []int32{
			8, 3, -4, 90,
		},
		Name2: "Name2",
	}
	if !reflect.DeepEqual(req, expect) {
		panic("Hello_Struct0")
	}
	return req
}
`
