package go2cpp

import (
	"github.com/orestonce/go2cpp/internal/testdata"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGo2cppContext_Generate1(t *testing.T) {
	ctx := NewGo2cppContext(NewGo2cppContext_Req{
		CppBaseName:                 "InProcessRpc",
		EnableQtClass_RunOnUiThread: false,
		EnableQtClass_Toast:         false,
		NotRemoveImplDotGo:          true,
	})

	ctx.Generate1(testdata.Hello_EmptyArg)
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

	ctx.Generate1(testdata.Hello_Float32)
	ctx.Generate1(testdata.Hello_Float64)

	ctx.Generate1(testdata.Hello_IntCommon)

	ctx.Generate1(testdata.Hello_StringEmpty)
	ctx.Generate1(testdata.Hello_StringCommon0)
	ctx.Generate1(testdata.Hello_StringCommon1)

	ctx.Generate1(testdata.Hello_Slice0)

	ctx.Generate1(testdata.Hello_Struct0)
	ctx.Generate1(testdata.HelloStruct1)
	ctx.Generate1(testdata.Hello_Struct2)
	ctx.Generate1(testdata.Hello_Block)

	ctx.Generate1(testdata.Hello_OutPkg)
	ctx.Generate1(testdata.Hello_Struct3)
	ctx.Generate1(testdata.Hello_Map)

	ctx.MustCreateAmd64LibraryInDir("tmp/temp1")
	//ctx.MustCreate386LibraryInDir("tmp/temp1")

	err := ioutil.WriteFile("tmp/temp1/main.cpp", []byte(mainCpp), 0777)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("g++", "InProcessRpc.cpp", "main.cpp", "-std=c++11", "-m64", "InProcessRpc-impl.a")
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
	os.RemoveAll("tmp")
}

const mainCpp = `
#include "InProcessRpc.h"
#include <math.h>
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
void testHello_Struct2();
void testHello_Struct3();
void testHello_Map();

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
	
	assert_ok(Hello_IntCommon(int(0x12345678)) == int(0x12345678), "Hello_IntCommon");

	float v = Hello_Float32(0.5678);
	assert_ok(fabs(v - 1234.56) < 0.01, "Hello_Float32");	// float32 有效位数为6位
	assert_ok(fabs(Hello_Float64(0.5678) - 1234.5678) < 1e-6, "Hello_Float64");
	
	testHello_StringXXX();
	
	testHello_Slice0();
	
	testHello_Struct0();
	testHello_Struct2();
	testHello_Struct3();
	testHello_Map();
	
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

void testHello_Struct2()
{
	Hello_Struct2Req origin;
	Hello_Struct0ReqL1 data;
	data.Name = "n2";
	data.Age  = 1;
	origin.Data.push_back(data);
	
	data.Name = "n8";
	data.Age  = 9;
	origin.Data.push_back(data);
	
	Hello_Struct2Req after = Hello_Struct2(origin);

	assert_ok(after.Data.size() == 2, "Hello_Struct2.size");
	assert_ok(after.Data[0].Name == "n2", "Hello_Struct2.0.Name");
	assert_ok(after.Data[0].Age == 1, "Hello_Struct2.0.Age");

	assert_ok(after.Data[1].Name == "n8", "Hello_Struct2.1.Name");
	assert_ok(after.Data[1].Age == 9, "Hello_Struct2.1.Name");
}

void testHello_Struct3()
{
	Struct3_L0 L0;
	Struct3_L1 tmp0;
	tmp0.L2.Id = 1;
	L0.S.push_back(tmp0);
	
	Struct3_L1 tmp1;
	tmp1.L2.Id = 2;
	L0.S.push_back(tmp1);
	
	Struct3_L0 resp = Hello_Struct3(L0);
	
	assert_ok(resp.S.size() == 2, "testHello_Struct3.Size");
	assert_ok(resp.S[0].L2.Id == 1, "testHello_Struct3.[1]");
	assert_ok(resp.S[1].L2.Id == 2, "testHello_Struct3.[2]");
}

void testHello_Map()
{
	std::map<std::string, Struct4> in;

	Struct4 v1;
	v1.V = "1";
	v1.I = 1;

	Struct4 v2;
	v2.V = "2";
	v2.I = 2;

	in[v1.V] = v1;
	in[v2.V] = v2;
	Hello_Map_Req req;
	req.MData = in;

	Hello_Map_Resp resp = Hello_Map(req);
	std::map<int, Struct4> out = resp.MData;
	assert_ok(out.size() == 2, "testHello_Map.out.size");
	assert_ok(out[1].I == 1, "testHello_Map.out[1].I");
	assert_ok(out[1].V == std::string("1"), "testHello_Map.out[1].V");
	assert_ok(out[2].I == 2, "testHello_Map.out[2].I");
	assert_ok(out[2].V == std::string("2"), "testHello_Map.out[2].V");
}
`
