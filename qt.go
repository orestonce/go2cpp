package go2cpp

import (
	"bytes"
)

func (this *Go2cppContext) appendQtIncludeH(qtH *bytes.Buffer) {
	qtH.WriteString(`#include <QThreadPool>
#include <QObject>
`)
	for _, hName := range this.req.QtIncludeList {
		qtH.WriteString("#include <" + hName + ">\n")
	}
}

func (this *Go2cppContext) appendQtDotHDefine(qtH *bytes.Buffer) {
	declear1 := func(qtCfg qtMethodCfg, buf *bytes.Buffer) { // signal
		buf.WriteString("signal_" + qtCfg.methodName + "_start(")
		for idx := 0; idx < qtCfg.methodType.NumIn(); idx++ {
			in := qtCfg.methodType.In(idx)
			if idx > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(this.goType2Cpp(in))
		}
		buf.WriteString(")")
	}

	qtH.WriteString(`
// qt相关接口
class GoCallObject : public ` + this.req.QtExtendBaseClass + `
{
    Q_OBJECT
public:
    explicit GoCallObject(QObject *parent = 0);
    ~GoCallObject();
    static void RegisterMetaType();
`)
	qtH.WriteString("signals:	// 暴露给外部的接口\n")
	for _, qtCfg := range this.qtMethodList {
		qtH.WriteString("\tvoid ")
		declear1(qtCfg, qtH)
		qtH.WriteString(";\n")
	}

	qtH.WriteString("protected slots:  // 已自动确保以下调用在ui线程, 由子类重写可直接操作ui\n")
	for _, qtCfg := range this.qtMethodList {
		qtH.WriteString("\tvirtual void ")
		this.qtDeclear3(qtCfg, qtH)
		qtH.WriteString("{}\n")
	}
	qtH.WriteString("signals:	// 内部实现需要的信号\n")
	for _, qtCfg := range this.qtMethodList {
		qtH.WriteString("\tvoid ")
		this.qtDeclear2(qtCfg, qtH)
		qtH.WriteString(";\n")
	}
	qtH.WriteString(`private:
	QThreadPool m_pool;
};

GoCallObject* GetDefaultGoCallObject();
`)
}

func (this *Go2cppContext) appendQtIncludeCpp(buf *bytes.Buffer) {
	buf.WriteString(`#include <QtConcurrent/QtConcurrent>
#include <QMetaType>
`)

}

func (this *Go2cppContext) appendQtDotCppDefine(qtCpp *bytes.Buffer) {
	qtCpp.WriteString(`

GoCallObject::GoCallObject(QObject *parent) : `+this.req.QtExtendBaseClass+`(parent)
{
`)
	for _, qtCfg := range this.qtMethodList {
		qtCpp.WriteString("\t" + `connect(this, SIGNAL(`)
		this.qtDeclear2(qtCfg, qtCpp)
		qtCpp.WriteString(`), this, SLOT(`)
		this.qtDeclear3(qtCfg, qtCpp)
		qtCpp.WriteString(`));
    connect(this, &GoCallObject::signal_` + qtCfg.methodName + `_start, [this](`)
		for idx := 0; idx < qtCfg.methodType.NumIn(); idx++ {
			if idx > 0 {
				qtCpp.WriteString(", ")
			}
			qtCpp.WriteString(this.goType2Cpp(qtCfg.methodType.In(idx)) + " " + this.inArgName(idx))
		}
		qtCpp.WriteString(`){
        QtConcurrent::run(&this->m_pool, [this`)
		for idx := 0; idx < qtCfg.methodType.NumIn(); idx++ {
			qtCpp.WriteString(", " + this.inArgName(idx))
		}
		qtCpp.WriteString(`](){
` + "\t\t\t")
		if qtCfg.returnType != nil {
			qtCpp.WriteString(this.goType2Cpp(qtCfg.returnType) + " result = ")
		}
		qtCpp.WriteString(qtCfg.methodName + "_Block(")
		for idx := 0; idx < qtCfg.methodType.NumIn(); idx++ {
			if idx > 0 {
				qtCpp.WriteString(",")
			}
			qtCpp.WriteString(this.inArgName(idx))
		}
		qtCpp.WriteString(");\n")
		qtCpp.WriteString("\t\t\temit this->_" + qtCfg.methodName + "_end(")
		if qtCfg.returnType != nil {
			qtCpp.WriteString("result")
		}
		qtCpp.WriteString(");\n")
		qtCpp.WriteString("\t\t" + `});
    });
`)
	}
	qtCpp.WriteString(`}

GoCallObject::~GoCallObject()
{
    m_pool.waitForDone();
}

GoCallObject* GetDefaultGoCallObject()
{
	static GoCallObject obj;
    return &obj;
}

void GoCallObject::RegisterMetaType()
{
`)
	existsMap := map[string]struct{}{}
	for _, qtCfg := range this.qtMethodList {
		if qtCfg.returnType == nil {
			continue
		}
		name := this.goType2Cpp(qtCfg.returnType)
		_, ok := existsMap[name]
		if ok {
			continue
		}
		qtCpp.WriteString("\t" + `qRegisterMetaType<` + this.goType2Cpp(qtCfg.returnType) + `>("` + name + `");` + "\n")
	}
	qtCpp.WriteString(`}
`)
}

func (this *Go2cppContext) qtDeclear2(qtCfg qtMethodCfg, buf *bytes.Buffer) { // signal
	buf.WriteString("_" + qtCfg.methodName + "_end(")
	if qtCfg.returnType != nil {
		buf.WriteString(this.goType2Cpp(qtCfg.returnType))
	}
	buf.WriteString(")")
}

func (this *Go2cppContext) qtDeclear3(qtCfg qtMethodCfg, buf *bytes.Buffer) { // slot
	buf.WriteString("slot_" + qtCfg.methodName + "_end(")
	if qtCfg.returnType != nil {
		buf.WriteString(this.goType2Cpp(qtCfg.returnType))
	}
	buf.WriteString(")")
}
