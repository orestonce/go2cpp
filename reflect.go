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
	req               NewGo2cppContext_Req
	importMap         map[string]struct{}
	cppTypeDeclareMap map[string]struct{}
	dotH              bytes.Buffer
	dotCpp            bytes.Buffer
	dotGo             bytes.Buffer
	idx               int
}

type NewGo2cppContext_Req struct {
	CppBaseName                 string
	EnableQtClass_RunOnUiThread bool
	EnableQtClass_Toast         bool
}

func NewGo2cppContext(req NewGo2cppContext_Req) *Go2cppContext {
	return &Go2cppContext{
		req:               req,
		importMap:         map[string]struct{}{},
		cppTypeDeclareMap: map[string]struct{}{},
	}
}

func (this *Go2cppContext) getNextVarName() string {
	this.idx++
	return "tmp" + strconv.Itoa(this.idx)
}

func (this *Go2cppContext) Generate1(methodFn interface{}) {
	this.idx = 0
	fnType := reflect.TypeOf(methodFn)
	pkgName, shortPkgName, methodName := splitAndValidatePkg(runtime.FuncForPC(reflect.ValueOf(methodFn).Pointer()).Name())

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
		case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.String, reflect.Struct, reflect.Slice, reflect.Map, reflect.Float64, reflect.Float32:
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
	buf.WriteString(this.dotCpp.String() + "\n")
	if this.req.EnableQtClass_RunOnUiThread {
		buf.WriteString(runOnUiThread_dotCpp)
	}
	if this.req.EnableQtClass_Toast {
		buf.WriteString(toast_dotCpp)
	}
	return buf.Bytes()
}

func (this *Go2cppContext) GetDotHContent() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("#pragma once\n\n")
	buf.WriteString("#include <string>\n")
	buf.WriteString("#include <vector>\n")
	buf.WriteString("#include <cstdint>\n")
	buf.WriteString("#include <map>\n")
	buf.WriteString("//Qt Creator 需要在xxx.pro 内部增加静态库的链接声明\n")
	buf.WriteString("//LIBS += -L$$PWD -l" + this.req.CppBaseName + "-impl\n")
	buf.WriteString("\n")
	buf.WriteString(this.dotH.String())
	if this.req.EnableQtClass_RunOnUiThread {
		buf.WriteString(runOnUiThread_dotH)
	}
	if this.req.EnableQtClass_Toast {
		buf.WriteString(toast_dotH)
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
	this.mustCreateLibrary(dir, "386")
}

func (this *Go2cppContext) MustCreateAmd64LibraryInDir(dir string) {
	this.mustCreateLibrary(dir, "amd64")
}

func (this *Go2cppContext) mustCreateLibrary(dir string, goarch string) {
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			panic(err)
		}
	}

	writeFile(filepath.Join(dir, this.req.CppBaseName+"-impl.go"), this.GetDotGoContent())
	cmd := exec.Command("go", "build", "-buildmode=c-archive", this.req.CppBaseName+"-impl.go")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1", "GOARCH="+goarch)
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
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.Struct:
		return out.Name()
	case reflect.Slice:
		return "std::vector<" + this.goType2Cpp(out.Elem()) + ">"
	case reflect.Map:
		return "std::map<" + this.goType2Cpp(out.Key()) + ", " + this.goType2Cpp(out.Elem()) + ">"
	default:
		panic("goType2Cpp: " + out.Kind().String())
	}
}

func (this *Go2cppContext) cppDecode(prefix string, name string, fType reflect.Type) {
	buf := &this.dotCpp
	kind := fType.Kind()
	decodeUint32 := func(prefix, name0 string) {
		varNameA := this.getNextVarName()
		varNameB := this.getNextVarName()
		varNameC := this.getNextVarName()
		varNameD := this.getNextVarName()
		buf.WriteString(prefix + "uint32_t " + varNameA + " = uint32_t(uint8_t(out[outIdx+0]) << 24);\n")
		buf.WriteString(prefix + "uint32_t " + varNameB + " = uint32_t(uint8_t(out[outIdx+1]) << 16);\n")
		buf.WriteString(prefix + "uint32_t " + varNameC + " = uint32_t(uint8_t(out[outIdx+2]) << 8);\n")
		buf.WriteString(prefix + "uint32_t " + varNameD + " = uint32_t(uint8_t(out[outIdx+3]) << 0);\n")
		buf.WriteString(prefix + name0 + " = " + varNameA + " | " + varNameB + " | " + varNameC + " | " + varNameD + ";\n")
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
		varName := this.getNextVarName()
		buf.WriteString(prefix + "\t" + "uint32_t " + varName + " = 0;\n")
		decodeUint32(prefix+"\t", varName)
		buf.WriteString(prefix + "\t" + name + " = std::string(out+outIdx, out+outIdx+" + varName + ");\n")
		buf.WriteString(prefix + "\t" + "outIdx+=" + varName + ";\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Float32, reflect.Float64:
		buf.WriteString(prefix + "{\n")
		varName := this.getNextVarName()
		buf.WriteString(prefix + "\t" + "uint32_t " + varName + " = 0;\n")
		decodeUint32(prefix+"\t", varName)
		sName := this.getNextVarName()
		buf.WriteString(prefix + "\tstd::string " + sName + " = std::string(out+outIdx, out+outIdx+" + varName + ");\n")
		buf.WriteString(prefix + "\t" + "outIdx+=" + varName + ";\n")
		fmtS := "%f"
		if kind == reflect.Float64 {
			fmtS = "%lf"
		}
		buf.WriteString("sscanf(" + sName + ".c_str(), \"" + fmtS + "\", &" + name + ");\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Struct:
		buf.WriteString(prefix + "{\n")
		for idx2 := 0; idx2 < fType.NumField(); idx2++ {
			this.cppDecode(prefix+"\t", name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
		}
		buf.WriteString(prefix + "}\n")
	case reflect.Slice:
		buf.WriteString(prefix + "{\n")
		varName0 := this.getNextVarName()
		buf.WriteString(prefix + "\t" + "uint32_t " + varName0 + " = 0;\n")
		decodeUint32(prefix+"\t", varName0)
		varName1 := this.getNextVarName()
		buf.WriteString(prefix + "\t" + "for (uint32_t " + varName1 + " = 0; " + varName1 + " < " + varName0 + "; " + varName1 + "++) {\n")
		varName2 := this.getNextVarName()
		buf.WriteString(prefix + "\t\t" + this.goType2Cpp(fType.Elem()) + " " + varName2 + ";\n")
		this.cppDecode(prefix+"\t\t", varName2, fType.Elem())
		buf.WriteString(prefix + "\t\t" + name + ".push_back(" + varName2 + ");\n")
		buf.WriteString(prefix + "\t" + "}\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Map:
		buf.WriteString(prefix + "{\n")
		varName0 := this.getNextVarName()
		buf.WriteString(prefix + "\t" + "uint32_t " + varName0 + " = 0;\n")
		decodeUint32(prefix+"\t", varName0)
		varName1 := this.getNextVarName()
		buf.WriteString(prefix + "\t" + "for (uint32_t " + varName1 + " = 0; " + varName1 + " < " + varName0 + "; " + varName1 + "++) {\n")
		kName := this.getNextVarName()
		buf.WriteString(prefix + "\t\t" + this.goType2Cpp(fType.Key()) + " " + kName + ";\n")
		this.cppDecode(prefix+"\t\t", kName, fType.Key())
		vName := this.getNextVarName()
		buf.WriteString(prefix + "\t\t" + this.goType2Cpp(fType.Elem()) + " " + vName + ";\n")
		this.cppDecode(prefix+"\t\t", vName, fType.Elem())
		buf.WriteString(prefix + "\t\t" + name + "[" + kName + "] = " + vName + ";\n")
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
			case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.Float32, reflect.Float64:
				buf.WriteString("var " + this.inArgName(idx) + " " + kind.String() + "\n")
				this.goDecode(this.inArgName(idx), fnType.In(idx))
			case reflect.String:
				buf.WriteString("var " + this.inArgName(idx) + " string\n")
				this.goDecode(this.inArgName(idx), fnType.In(idx))
			case reflect.Struct, reflect.Slice:
				buf.WriteString("var " + this.inArgName(idx) + " " + in.String() + "\n")
				this.goDecode(this.inArgName(idx), fnType.In(idx))
			case reflect.Map:
				buf.WriteString("var " + this.inArgName(idx) + " = " + in.String() + "{}\n")
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
		varName := this.getNextVarName()
		buf.WriteString(varName + " := make([]byte, 4)\n")
		buf.WriteString("inBuf.Read(" + varName + ")\n")
		buf.WriteString(name0 + " = " + kind0 + "(binary.BigEndian.Uint32(" + varName + "))\n")
	}
	if fType.PkgPath() != `` {
		this.importMap[fType.PkgPath()] = struct{}{}
	}
	switch kind {
	case reflect.Bool:
		buf.WriteString("{\n")
		varName := this.getNextVarName()
		buf.WriteString(varName + ", _ := inBuf.ReadByte()\n")
		buf.WriteString(name + " = (" + varName + " != 0)\n")
		buf.WriteString("}\n")
	case reflect.Int8, reflect.Uint8:
		buf.WriteString("{\n")
		varName := this.getNextVarName()
		buf.WriteString(varName + ", _ := inBuf.ReadByte()\n")
		buf.WriteString(name + " = " + kind.String() + "(" + varName + ")\n")
		buf.WriteString("}\n")
	case reflect.Int, reflect.Int32, reflect.Uint32:
		buf.WriteString("{\n")
		decodeInt(name, kind.String())
		buf.WriteString("}\n")
	case reflect.String:
		buf.WriteString("{\n")
		varName := this.getNextVarName()
		buf.WriteString("var " + varName + " int\n")
		decodeInt(varName, "int")
		varName2 := this.getNextVarName()
		buf.WriteString(varName2 + " := make([]byte, " + varName + ")\n")
		buf.WriteString("inBuf.Read(" + varName2 + ")\n")
		buf.WriteString(name + " = string(" + varName2 + ")\n")
		buf.WriteString("}\n")
	case reflect.Float32, reflect.Float64:
		buf.WriteString("{\n")
		varName := this.getNextVarName()
		buf.WriteString("var " + varName + " int\n")
		decodeInt(varName, "int")
		varName2 := this.getNextVarName()
		buf.WriteString(varName2 + " := make([]byte, " + varName + ")\n")
		buf.WriteString("inBuf.Read(" + varName2 + ")\n")
		this.importMap["strconv"] = struct{}{}
		nameV, nameErr := this.getNextVarName(), this.getNextVarName()
		buf.WriteString(nameV + ", " + nameErr + " := strconv.ParseFloat(string(" + varName2 + "), 64)\n")
		buf.WriteString("if " + nameErr + " != nil {\n")
		buf.WriteString("panic(" + nameErr + ")\n")
		buf.WriteString("}\n")
		buf.WriteString(name + " = " + kind.String() + "(" + nameV + ")\n")
		buf.WriteString("}\n")
	case reflect.Struct:
		buf.WriteString("{\n")
		if fType.NumField() == 0 {
			buf.WriteString(" _ = inBuf // 防止inBuf无引用错误\n")
		} else {
			for idx2 := 0; idx2 < fType.NumField(); idx2++ {
				this.goDecode(name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
			}
		}
		buf.WriteString("}\n")
	case reflect.Slice:
		buf.WriteString("{\n")
		varName := this.getNextVarName()
		buf.WriteString("var " + varName + " int\n")
		decodeInt(varName, "int")
		buf.WriteString(name + " = make(" + fType.String() + ", " + varName + ")\n")
		varName2 := this.getNextVarName()
		buf.WriteString("for " + varName2 + " := 0; " + varName2 + " < " + varName + "; " + varName2 + "++ {\n")
		this.goDecode(name+"["+varName2+"]", fType.Elem())
		buf.WriteString("}\n")
		buf.WriteString("}\n")
	case reflect.Map:
		buf.WriteString("{\n")
		varName := this.getNextVarName()
		buf.WriteString("var " + varName + " int\n")
		decodeInt(varName, "int")
		varName2 := this.getNextVarName()
		buf.WriteString("for " + varName2 + " := 0; " + varName2 + " < " + varName + "; " + varName2 + "++ {\n")
		kName := this.getNextVarName()
		vName := this.getNextVarName()
		buf.WriteString("var " + kName + " " + fType.Key().String() + ";\n")
		buf.WriteString("var " + vName + " " + fType.Elem().String() + ";\n")
		this.goDecode(kName, fType.Key())
		this.goDecode(vName, fType.Elem())
		buf.WriteString(name + "[" + kName + "] = " + vName + ";\n")
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
		varName := this.getNextVarName()
		buf.WriteString(varName + " := make([]byte, 4)\n")
		buf.WriteString("binary.BigEndian.PutUint32(" + varName + ", uint32(" + name0 + "))\n")
		buf.WriteString("outBuf.Write(" + varName + ")\n")
	}
	switch kind {
	case reflect.Bool:
		buf.WriteString("{\n")
		varName := this.getNextVarName()
		buf.WriteString("var " + varName + " byte\n")
		buf.WriteString("if " + name + " {\n")
		buf.WriteString(varName + " = 1\n")
		buf.WriteString("}\n")
		buf.WriteString("outBuf.WriteByte(" + varName + ")\n")
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
	case reflect.Float64, reflect.Float32:
		buf.WriteString("{\n")
		nameS := this.getNextVarName()
		buf.WriteString(nameS + " := strconv.FormatFloat(float64(" + name + "), 'f', -1, 64)\n")
		encodeUint32("len(" + nameS + ")")
		buf.WriteString("outBuf.WriteString(" + nameS + ")\n")
		buf.WriteString("}\n")
	case reflect.Struct:
		buf.WriteString("{\n")
		if fType.NumField() == 0 {
			buf.WriteString("_ = outArg // 防止outArg无引用错误\n")
		} else {
			for idx2 := 0; idx2 < fType.NumField(); idx2++ {
				this.goEncode(name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
			}
		}
		buf.WriteString("}\n")
	case reflect.Slice:
		buf.WriteString("{\n")
		encodeUint32("len(" + name + ")")
		varName := this.getNextVarName()
		buf.WriteString("for " + varName + " := 0; " + varName + " < len(" + name + "); " + varName + "++{\n")
		this.goEncode(name+"["+varName+"]", fType.Elem())
		buf.WriteString("}\n")
		buf.WriteString("}\n")
	case reflect.Map:
		buf.WriteString("{\n")
		encodeUint32("len(" + name + ")")
		kName := this.getNextVarName()
		vName := this.getNextVarName()
		buf.WriteString("for " + kName + ", " + vName + " := range " + name + "{\n")
		this.goEncode(kName, fType.Key())
		this.goEncode(vName, fType.Elem())
		buf.WriteString("}\n")
		buf.WriteString("}\n")
	default:
		panic("genGo " + kind.String())
	}
}

func (this *Go2cppContext) cppEncode(prefix string, name string, fType reflect.Type) {
	buf := &this.dotCpp
	encodeUint32 := func(prefix, name0 string) {
		varName := this.getNextVarName()
		buf.WriteString(prefix + `char ` + varName + `[4];
` + prefix + varName + `[0] = (uint32_t(` + name0 + `) >> 24) & 0xFF;
` + prefix + varName + `[1] = (uint32_t(` + name0 + `) >> 16) & 0xFF;
` + prefix + varName + `[2] = (uint32_t(` + name0 + `) >> 8) & 0xFF;
` + prefix + varName + `[3] = (uint32_t(` + name0 + `) >> 0) & 0xFF;
` + prefix + `in.append(` + varName + `, 4);` + "\n")
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
		varName := this.getNextVarName()
		buf.WriteString(prefix + "\t" + `uint32_t ` + varName + ` = ` + name + ".length();\n")
		encodeUint32(prefix+"\t", varName)
		buf.WriteString(prefix + "\t" + `in.append(` + name + ");\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Float32, reflect.Float64:
		buf.WriteString(prefix + "{\n")
		bufName := this.getNextVarName()
		buf.WriteString("char " + bufName + "[64] = \"\";\n")
		nName := this.getNextVarName()
		buf.WriteString("int " + nName + " = snprintf(" + bufName + ", 64, \"%f\", " + name + ");\n")
		encodeUint32(prefix+"\t", nName)
		buf.WriteString(prefix + "\t" + `in.append(` + bufName + ", size_t(" + nName + "));\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Struct:
		buf.WriteString(prefix + "{\n")
		for idx2 := 0; idx2 < fType.NumField(); idx2++ {
			this.cppEncode(prefix+"\t", name+"."+fType.Field(idx2).Name, fType.Field(idx2).Type)
		}
		buf.WriteString(prefix + "}\n")
	case reflect.Slice:
		buf.WriteString(prefix + "{\n")
		varName := this.getNextVarName()
		buf.WriteString(prefix + "\t" + `uint32_t ` + varName + ` = ` + name + ".size();\n")
		encodeUint32(prefix+"\t", varName)
		varName2 := this.getNextVarName()
		buf.WriteString(prefix + "\tfor (uint32_t " + varName2 + "=0; " + varName2 + " < " + varName + "; ++" + varName2 + ") {\n")
		this.cppEncode(prefix+"\t\t", name+"["+varName2+"]", fType.Elem())
		buf.WriteString(prefix + "\t}\n")
		buf.WriteString(prefix + "}\n")
	case reflect.Map:
		buf.WriteString(prefix + "{\n")
		varName := this.getNextVarName()
		buf.WriteString(prefix + "\t" + `uint32_t ` + varName + ` = ` + name + ".size();\n")
		encodeUint32(prefix+"\t", varName)
		itName := this.getNextVarName()
		buf.WriteString(prefix + "\tfor(" + this.goType2Cpp(fType) + "::iterator " + itName + " = " + name + ".begin(); " + itName + " != " + name + ".end(); ++" + itName + ") {\n")
		this.cppEncode(prefix+"\t\t", itName+"->first", fType.Key())
		this.cppEncode(prefix+"\t\t", itName+"->second", fType.Elem())
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
		case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32, reflect.String, reflect.Float32, reflect.Float64:
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
			declareTypeBefore(in.Elem())
		case reflect.Map:
			for _, one := range []reflect.Type{in.Key(), in.Elem()} {
				typeName = this.goType2Cpp(one)
				_, ok := this.cppTypeDeclareMap[typeName]
				if ok {
					continue
				}
				declareTypeBefore(one)
			}
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
			baseTypeNameValue := map[string]string{}
			buf.WriteString("struct " + this.goType2Cpp(in) + "{\n")
			for idx1 := 0; idx1 < in.NumField(); idx1++ {
				t := in.Field(idx1).Type
				buf.WriteString("\t" + this.goType2Cpp(t) + " " + in.Field(idx1).Name + ";\n")
				switch t.Kind() {
				case reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int32, reflect.Uint32:
					baseTypeNameValue[in.Field(idx1).Name] = "0"
				case reflect.Bool:
					baseTypeNameValue[in.Field(idx1).Name] = "false"
				}
			}
			if len(baseTypeNameValue) > 0 {
				isFirst := true
				buf.WriteString("\t" + this.goType2Cpp(in) + "(): ")
				for idx1 := 0; idx1 < in.NumField(); idx1++ {
					name := in.Field(idx1).Name
					value := baseTypeNameValue[name]
					if value == "" {
						continue
					}
					if isFirst {
						isFirst = false
					} else {
						buf.WriteString(",")
					}
					buf.WriteString(name + "(" + value + ")")
				}
				buf.WriteString("{}\n")
			}
			buf.WriteString("};\n")
		}
	}
}
