package testdata

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func Hello_EmptyArg() {
}

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

		panic("Hello_IntMin " + strconv.Itoa(a))
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

type HelloStruct1Req struct{}
type HelloStruct1Resp struct{}

func HelloStruct1(req HelloStruct1Req) (resp HelloStruct1Resp) {
	return resp
}

type Hello_Struct2Req struct {
	Data []Hello_Struct0ReqL1
}

func Hello_Struct2(req Hello_Struct2Req) Hello_Struct2Req {
	expect := Hello_Struct2Req{
		Data: []Hello_Struct0ReqL1{
			{
				Name: "n2",
				Age:  1,
			},
			{
				Name: "n8",
				Age:  9,
			},
		},
	}
	if !reflect.DeepEqual(req, expect) {
		panic("Hello_Struct2")
	}
	return req
}

func Hello_Block(s string) (i int) {
	time.Sleep(time.Second)
	return len(s)
}

func Hello_OutPkg(inArg sort.IntSlice) (outArg unicode.Range32) {
	return outArg
}

type Struct3_L2 struct {
	Id int
}

type Struct3_L1 struct {
	L2 Struct3_L2
}

type Struct3_L0 struct {
	S []Struct3_L1
	B bool
}

func Hello_Struct3(req Struct3_L0) (resp Struct3_L0) {
	if len(req.S) != 2 {
		panic("Hello_StructL0 1")
	}
	if req.S[0].L2.Id != 1 {
		panic(`Hello_StructL0 2`)
	}
	if req.S[1].L2.Id != 2 {
		panic(`Hello_StructL0 3`)
	}
	return req
}

type Struct4 struct {
	V string
	I int
}

func Hello_Map(m map[string]Struct4) (m2 map[int]Struct4) {
	if len(m) != 2 {
		panic("len(m) != 2")
	}
	m2 = map[int]Struct4{}
	for _, v := range m {
		m2[v.I] = v
		if v.V != strconv.Itoa(v.I) {
			panic("v.V " + v.V)
		}
	}
	return m2
}
