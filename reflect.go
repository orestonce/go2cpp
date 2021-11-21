package go2cpp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type Go2cppContext struct {
	req               NewGo2cppContextReq
	importMap         map[string]struct{}
	cppTypeDeclareMap map[string]struct{}
	dotH              bytes.Buffer
	dotCpp            bytes.Buffer
	dotGo             bytes.Buffer
	qtMethodList      []qtMethodCfg
}

type qtMethodCfg struct {
	methodType reflect.Type
	methodName string // TrimSuffix(x, "_Block")
	returnType reflect.Type
}

type NewGo2cppContextReq struct {
	CppBaseName       string
	EnableQt          bool
	QtExtendBaseClass string
	QtIncludeList     []string
}

func NewGo2cppContext(req NewGo2cppContextReq) *Go2cppContext {
	return &Go2cppContext{
		req:               req,
		importMap:         map[string]struct{}{},
		cppTypeDeclareMap: map[string]struct{}{},
	}
}

func (this *Go2cppContext) Generate1(methodFn interface{}) {
	fnType := reflect.TypeOf(methodFn)
	pkgName, shortPkgName, methodName := splitAndValidatePkg(runtime.FuncForPC(reflect.ValueOf(methodFn).Pointer()).Name())

	if strings.HasSuffix(methodName, "_Block") { // Qt阻塞式调用
		this.qtMethodList = append(this.qtMethodList, qtMethodCfg{
			methodType: fnType,
			methodName: strings.TrimSuffix(methodName, "_Block"),
			returnType: func() reflect.Type {
				if fnType.NumOut() != 1 {
					return nil
				}
				return fnType.Out(0)
			}(),
		})
	}

	this.goFnDeclare(pkgName, shortPkgName, methodName, fnType)
	this.cppTypeDeclare(fnType)
	this.cppFnDeclare(fnType, methodName)
	this.dotCpp.WriteString(`{
	std::string in;
`)
	for idx := 0; idx < fnType.NumIn(); idx++ {
		this.cppEncode("\t", this.inArgName(idx), fnType.In(idx))
	}
	this.dotCpp.WriteString(`	char *out = NULL;
	int outLen = 0;
	` + this.goFnName(methodName) + "((char *)in.data(), in.length(), &out, &outLen);\n")

	if fnType.NumOut() == 1 {
		out := fnType.Out(0)
		kind := out.Kind()
		buf := &this.dotCpp
		switch kind {
		case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.String, reflect.Struct, reflect.Slice:
			buf.WriteString("\t" + this.goType2Cpp(out) + " retValue;\n")
			buf.WriteString("\t" + "int outIdx = 0;\n")
			this.cppDecode("\t", "retValue", out)
		default:
			panic("genGo " + kind.String())
		}
	}
	this.dotCpp.WriteString("\t" + `if (out != NULL) {
		free(out);
	}
`)
	if fnType.NumOut() == 1 {
		this.dotCpp.WriteString(`	return retValue;` + "\n")
	}
	this.dotCpp.WriteString("}\n\n")
	this.dotGo.WriteString("}\n\n")
}

func (this *Go2cppContext) GetDotCppContent(implDotHContent []byte) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(`#include ` + strconv.Quote(this.req.CppBaseName+".h") + "\n")
	buf.Write(implDotHContent)
	buf.WriteString("\n\n")
	if this.req.EnableQt {
		this.appendQtIncludeCpp(buf)
	}
	buf.WriteString(this.dotCpp.String())
	if this.req.EnableQt {
		this.appendQtDotCppDefine(buf)
	}
	return buf.Bytes()
}

func (this *Go2cppContext) GetDotHContent() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("#pragma once\n\n")
	buf.WriteString("#include <string>\n")
	buf.WriteString("#include <vector>\n")
	buf.WriteString("#include <cstdint>\n")
	if this.req.EnableQt {
		buf.WriteString("//在xxx.pro 内部增加静态库的链接声明\n")
		buf.WriteString("//LIBS += -L$$PWD -l" + this.req.CppBaseName + "-impl\n")
		this.appendQtIncludeH(buf)
	}
	buf.WriteString("\n")
	buf.WriteString(this.dotH.String())
	if this.req.EnableQt {
		this.appendQtDotHDefine(buf)
	}
	return buf.Bytes()
}

func (this *Go2cppContext) GetDotGoContent() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(`package main

`)
	var importList []string
	for key := range this.importMap {
		importList = append(importList, key)
	}
	sort.Strings(importList)
	for _, key := range importList {
		buf.WriteString(`import ` + strconv.Quote(key) + "\n")
	}
	buf.WriteString(`import "C"` + "\n\n")
	buf.WriteString("func main(){}\n")
	buf.WriteString(this.dotGo.String())
	return buf.Bytes()
}

func (this *Go2cppContext) MustCreate386LibraryInDir(dir string) {
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(err)
		}
	}

	writeFile(filepath.Join(dir, this.req.CppBaseName+"-impl.go"), this.GetDotGoContent())
	cmd := exec.Command("go", "build", "-buildmode=c-archive", this.req.CppBaseName+"-impl.go")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1", "GOARCH=386")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
	writeFile(filepath.Join(dir, this.req.CppBaseName+".h"), this.GetDotHContent())
	implDotHContent, err := ioutil.ReadFile(filepath.Join(dir, this.req.CppBaseName+"-impl.h"))
	if err != nil {
		panic(err)
	}
	writeFile(filepath.Join(dir, this.req.CppBaseName+".cpp"), this.GetDotCppContent(implDotHContent))
	err = os.Remove(filepath.Join(dir, this.req.CppBaseName+"-impl.go"))
	if err != nil {
		panic(err)
	}
	err = os.Remove(filepath.Join(dir, this.req.CppBaseName+"-impl.h"))
	if err != nil {
		panic(err)
	}
}

func writeFile(dst string, content []byte) {
	err := ioutil.WriteFile(dst, content, 0777)
	if err != nil {
		panic(err)
	}
}

func splitAndValidatePkg(goMethod string) (pkgName string, shortPkgName string, methodName string) {
	tmp := strings.Split(goMethod, ".")
	if len(tmp) < 2 {
		panic("invalid goMethod: " + strconv.Quote(goMethod))
	}
	pkgName, methodName = strings.Join(tmp[0:len(tmp)-1], "."), tmp[len(tmp)-1]
	if strings.Contains(pkgName, "/") {
		tmp = strings.Split(pkgName, "/")
		shortPkgName = tmp[len(tmp)-1]
	} else {
		shortPkgName = pkgName
	}
	return pkgName, shortPkgName, methodName
}

func (this *Go2cppContext) genInArg(writer io.Writer, idx int, in reflect.Type) {
	if idx > 0 {
		writer.Write([]byte(", "))
	}
	writer.Write([]byte(this.goType2Cpp(in) + " " + this.inArgName(idx)))
}

func (this *Go2cppContext) genOut(out reflect.Type) {
	fmt.Println("genOut", out.Name())
}

func (this *Go2cppContext) goType2Cpp(out reflect.Type) string {
	switch out.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int8:
		return "int8_t"
	case reflect.Uint8:
		return "uint8_t"
	case reflect.Int32, reflect.Int:
		return "int32_t"
	case reflect.Uint32:
		return "uint32_t"
	case reflect.String:
		return "std::string"
	case reflect.Struct:
		return out.Name()
	case reflect.Slice:
		return "std::vector<" + this.goType2Cpp(out.Elem()) + ">"
	default:
		panic("goType2Cpp: " + out.Kind().String())
	}
}

func (this *Go2cppContext) cppDecode(prefix string, name string, fType reflect.Type) {
	buf := &this.dotCpp
	kind := fType.Kind()
	decodeUint32 := func(prefix, name0 string) {
		buf.WriteString(prefix + "uint32_t a = uint32_t(uint8_t(out[outIdx+0]) << 24);\n")
		buf.WriteString(prefix + "uint32_t b = uint32_t(uint8_t(out[outIdx+1]) << 16);\n")
		buf.WriteString(prefix + "uint32_t c = uint32_t(uint8_t(out[outIdx+2]) << 8);\n")
		buf.WriteString(prefix + "uint32_t d = uint32_t(uint8_t(out[outIdx+3]) << 0);\n")
		buf.WriteString(prefix + name0 + " = a | b | c | d;\n")
		buf.WriteString(prefix + "outIdx+=4;\n")
	}
	switch kind {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		buf.WriteString(prefix + name + " = (" + this.goType2Cpp(fType) + ") out[outIdx];\n")
		buf.WriteString(prefix + "outIdx++;\n")
	case reflect.Int, reflect.Int32, reflect.Uint32:
		buf.WriteString(prefix + "{\n")
		decodeUint32(prefix+"\t", name)
		buf.WriteString(prefix + "}\n")
	case reflect.String:
		buf.WriteString(prefix + "{\n")
		buf.WriteString(prefix + "\t" + "uint32_t length = 0;\n")
		decodeUint32(prefix+"\t", "length")
		buf.WriteString(prefix + "\t" + name + " = std::string(out+outIdx, out+outIdx+length);\n")
		buf.WriteString(prefix + "\t" + "outIdx+=length;\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Struct:
		buf.WriteString(prefix + "{\n")
		for idx2 := 0; idx2 < fType.NumField(); idx2++ {
			this.cppDecode(prefix+"\t", name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
		}
		buf.WriteString(prefix + "}\n")
	case reflect.Slice:
		buf.WriteString(prefix + "{\n")
		buf.WriteString(prefix + "\t" + "uint32_t length = 0;\n")
		decodeUint32(prefix+"\t", "length")
		buf.WriteString(prefix + "\t" + "for (uint32_t idx = 0; idx < length; idx++) {\n")
		buf.WriteString(prefix + "\t\t" + this.goType2Cpp(fType.Elem()) + " tmp;\n")
		this.cppDecode(prefix+"\t\t", "tmp", fType.Elem())
		buf.WriteString(prefix + "\t\t" + name + ".push_back(tmp);\n")
		buf.WriteString(prefix + "\t" + "}\n")
		buf.WriteString(prefix + "}\n")
	default:
		panic("genGo " + kind.String())
	}
}

func (this *Go2cppContext) inArgName(idx int) string {
	return "in" + strconv.Itoa(idx)
}

func (this *Go2cppContext) goFnName(name string) string {
	return "Go2cppFn_" + name
}

func (this *Go2cppContext) goFnDeclare(pkgName string, shortPkgName string, methodName string, fnType reflect.Type) {
	this.importMap[pkgName] = struct{}{}
	buf := &this.dotGo

	buf.WriteString(`//export ` + this.goFnName(methodName) + "\n")
	buf.WriteString("func " + this.goFnName(methodName) + "(in *C.char, inLen C.int, out **C.char, outLen *C.int){\n")
	if fnType.NumIn() > 0 {
		this.importMap["strings"] = struct{}{}
		buf.WriteString(`inBuf := strings.NewReader(C.GoStringN(in, inLen))` + "\n")
		for idx := 0; idx < fnType.NumIn(); idx++ {
			in := fnType.In(idx)

			kind := in.Kind()
			switch kind {
			case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32:
				buf.WriteString("var " + this.inArgName(idx) + " " + kind.String() + "\n")
				this.goDecode(this.inArgName(idx), fnType.In(idx))
			case reflect.String:
				buf.WriteString("var " + this.inArgName(idx) + " string\n")
				this.goDecode(this.inArgName(idx), fnType.In(idx))
			case reflect.Struct, reflect.Slice:
				buf.WriteString("var " + this.inArgName(idx) + " " + in.String() + "\n")
				this.goDecode(this.inArgName(idx), fnType.In(idx))
			default:
				panic("genGo " + kind.String())
			}
		}
	}
	if fnType.NumOut() == 1 {
		buf.WriteString(this.outArgName() + " := ")
	}
	buf.WriteString(shortPkgName + "." + methodName + "(")
	for idx := 0; idx < fnType.NumIn(); idx++ {
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(this.inArgName(idx))
	}
	buf.WriteString(")\n")
	if fnType.NumOut() == 1 {
		this.importMap["bytes"] = struct{}{}
		buf.WriteString("var outBuf bytes.Buffer\n")
		buf.WriteString("{\n")
		this.goEncode(this.outArgName(), fnType.Out(0))
		buf.WriteString("*out = C.CString(outBuf.String())\n")
		buf.WriteString("*outLen = C.int(outBuf.Len())\n")
		buf.WriteString("}\n")
	}
}

func (this *Go2cppContext) outArgName() string {
	return "outArg"
}

func (this *Go2cppContext) goDecode(name string, fType reflect.Type) {
	buf := &this.dotGo
	kind := fType.Kind()
	decodeInt := func(name0 string, kind0 string) {
		this.importMap["encoding/binary"] = struct{}{}
		buf.WriteString("tmp := make([]byte, 4)\n")
		buf.WriteString("inBuf.Read(tmp)\n")
		buf.WriteString(name0 + " = " + kind0 + "(binary.BigEndian.Uint32(tmp))\n")
	}
	switch kind {
	case reflect.Bool:
		buf.WriteString("{\n")
		buf.WriteString("value, _ := inBuf.ReadByte()\n")
		buf.WriteString(name + " = (value != 0)\n")
		buf.WriteString("}\n")
	case reflect.Int8, reflect.Uint8:
		buf.WriteString("{\n")
		buf.WriteString("value, _ := inBuf.ReadByte()\n")
		buf.WriteString(name + " = " + kind.String() + "(value)\n")
		buf.WriteString("}\n")
	case reflect.Int, reflect.Int32, reflect.Uint32:
		buf.WriteString("{\n")
		decodeInt(name, kind.String())
		buf.WriteString("}\n")
	case reflect.String:
		buf.WriteString("{\n")
		buf.WriteString("var length int\n")
		decodeInt("length", "int")
		buf.WriteString("data := make([]byte, length)\n")
		buf.WriteString("inBuf.Read(data)\n")
		buf.WriteString(name + " = string(data)\n")
		buf.WriteString("}\n")
	case reflect.Struct:
		buf.WriteString("{\n")
		for idx2 := 0; idx2 < fType.NumField(); idx2++ {
			this.goDecode(name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
		}
		buf.WriteString("}\n")
	case reflect.Slice:
		buf.WriteString("{\n")
		buf.WriteString("var length int\n")
		decodeInt("length", "int")
		buf.WriteString(name + " = make(" + fType.String() + ", length)\n")
		buf.WriteString("for idx := 0; idx < length; idx++ {\n")
		this.goDecode(name+"[idx]", fType.Elem())
		buf.WriteString("}\n")
		buf.WriteString("}\n")
	default:
		panic("genGo " + kind.String())
	}
}

func (this *Go2cppContext) goEncode(name string, fType reflect.Type) {
	kind := fType.Kind()
	buf := &this.dotGo
	encodeUint32 := func(name0 string) {
		this.importMap["encoding/binary"] = struct{}{}
		buf.WriteString("tmp := make([]byte, 4)\n")
		buf.WriteString("binary.BigEndian.PutUint32(tmp, uint32(" + name0 + "))\n")
		buf.WriteString("outBuf.Write(tmp)\n")
	}
	switch kind {
	case reflect.Bool:
		buf.WriteString("{\n")
		buf.WriteString("var value byte\n")
		buf.WriteString("if " + name + " {\n")
		buf.WriteString("value = 1\n")
		buf.WriteString("}\n")
		buf.WriteString("outBuf.WriteByte(value)\n")
		buf.WriteString("}\n")
	case reflect.Int8, reflect.Uint8:
		buf.WriteString("{\n")
		buf.WriteString("outBuf.WriteByte(byte(" + name + "))\n")
		buf.WriteString("}\n")
	case reflect.Int, reflect.Int32, reflect.Uint32:
		buf.WriteString("{\n")
		encodeUint32(name)
		buf.WriteString("}\n")
	case reflect.String:
		buf.WriteString("{\n")
		encodeUint32("len(" + name + ")")
		buf.WriteString("outBuf.WriteString(" + name + ")\n")
		buf.WriteString("}\n")
	case reflect.Struct:
		buf.WriteString("{\n")
		for idx2 := 0; idx2 < fType.NumField(); idx2++ {
			this.goEncode(name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
		}
		buf.WriteString("}\n")
	case reflect.Slice:
		buf.WriteString("{\n")
		encodeUint32("len(" + name + ")")
		buf.WriteString("for idx := 0; idx < len(" + name + "); idx++{\n")
		this.goEncode(name+"[idx]", fType.Elem())
		buf.WriteString("}\n")
		buf.WriteString("}\n")
	default:
		panic("genGo " + kind.String())
	}
}

func (this *Go2cppContext) cppEncode(prefix string, name string, fType reflect.Type) {
	buf := &this.dotCpp
	encodeUint32 := func(prefix, name0 string) {
		buf.WriteString(prefix + `char tmp[4];
` + prefix + `tmp[0] = (uint32_t(` + name0 + `) >> 24) & 0xFF;
` + prefix + `tmp[1] = (uint32_t(` + name0 + `) >> 16) & 0xFF;
` + prefix + `tmp[2] = (uint32_t(` + name0 + `) >> 8) & 0xFF;
` + prefix + `tmp[3] = (uint32_t(` + name0 + `) >> 0) & 0xFF;
` + prefix + `in.append(tmp, 4);` + "\n")
	}
	switch fType.Kind() {
	case reflect.Bool:
		buf.WriteString(prefix + `in.append((char*)(&` + name + "), 1);\n")
	case reflect.Uint8, reflect.Int8:
		buf.WriteString(prefix + `in.append((char*)(&` + name + "), 1);\n")
	case reflect.Int32, reflect.Int, reflect.Uint32:
		buf.WriteString(prefix + "{\n")
		encodeUint32(prefix+"\t", name)
		buf.WriteString(prefix + "}\n")
	case reflect.String:
		buf.WriteString(prefix + "{\n")
		buf.WriteString(prefix + "\t" + `uint32_t length = ` + name + ".length();\n")
		encodeUint32(prefix+"\t", "length")
		buf.WriteString(prefix + "\t" + `in.append(` + name + ");\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Struct:
		buf.WriteString(prefix + "{\n")
		for idx2 := 0; idx2 < fType.NumField(); idx2++ {
			this.cppEncode(prefix+"\t", name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
		}
		buf.WriteString(prefix + "}\n")
	case reflect.Slice:
		buf.WriteString(prefix + "{\n")
		buf.WriteString(prefix + "\t" + `uint32_t length = ` + name + ".size();\n")
		encodeUint32(prefix+"\t", "length")
		buf.WriteString(prefix + "\tfor (uint32_t idx=0; idx < length; idx++) {\n")
		this.cppEncode(prefix+"\t\t", name+"[idx]", fType.Elem())
		buf.WriteString(prefix + "\t}\n")
		buf.WriteString(prefix + "}\n")
	default:
		panic("cppEncode " + fType.Kind().String())
	}
}

func (this *Go2cppContext) cppFnDeclare(fnType reflect.Type, methodName string) {
	writer := io.MultiWriter(&this.dotCpp, &this.dotH)
	if fnType.NumOut() == 0 {
		writer.Write([]byte(`void `))
	} else if fnType.NumOut() == 1 {
		writer.Write([]byte(this.goType2Cpp(fnType.Out(0)) + " "))
	} else {
		panic("fn.NumOut " + strconv.Itoa(fnType.NumOut()))
	}
	writer.Write([]byte(methodName + "("))
	for idx := 0; idx < fnType.NumIn(); idx++ {
		this.genInArg(writer, idx, fnType.In(idx))
	}
	this.dotCpp.WriteString(")")
	this.dotH.WriteString(");\n")
}

func (this *Go2cppContext) cppTypeDeclare(fnType reflect.Type) {
	var typeList []reflect.Type

	var declareTypeBefore func(in reflect.Type)

	declareTypeBefore = func(in reflect.Type) {
		var typeName string
		switch in.Kind() {
		case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.String:
			return
		case reflect.Struct:
			typeName = this.goType2Cpp(in)
			_, ok := this.cppTypeDeclareMap[typeName]
			if ok {
				return
			}
			this.cppTypeDeclareMap[typeName] = struct{}{}
			for idx := 0; idx < in.NumField(); idx++ {
				declareTypeBefore(in.Field(idx).Type)
			}
			typeList = append(typeList, in)
		case reflect.Slice:
			typeName = this.goType2Cpp(in.Elem())
			_, ok := this.cppTypeDeclareMap[typeName]
			if ok {
				return
			}
			this.cppTypeDeclareMap[typeName] = struct{}{}
			declareTypeBefore(in.Elem())
			typeList = append(typeList, in.Elem())
		default:
			panic("cppTypeDeclare " + in.Kind().String())
		}
	}
	for idx := 0; idx < fnType.NumIn(); idx++ {
		declareTypeBefore(fnType.In(idx))
	}
	if fnType.NumOut() > 0 {
		declareTypeBefore(fnType.Out(0))
	}
	buf := &this.dotH
	for _, in := range typeList {
		if in.Kind() == reflect.Struct {
			buf.WriteString("struct " + this.goType2Cpp(in) + "{\n")
			for idx1 := 0; idx1 < in.NumField(); idx1++ {
				buf.WriteString("\t" + this.goType2Cpp(in.Field(idx1).Type) + " " + in.Field(idx1).Name + ";\n")
			}
			buf.WriteString("};\n")
		}
	}
}
