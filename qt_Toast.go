package go2cpp

const toast_dotH = `
// Thanks: https://github.com/live-in-a-dream/Qt-Toast

#include <QString>
#include <QObject>

class QTimer;
class QLabel;
class QWidget;

namespace Ui {
class Toast;
}

class Toast : public QObject
{
    Q_OBJECT

public:
    explicit Toast(QObject *parent = nullptr);

    static Toast* Instance();
    //错误
    void SetError(const QString &text,const int & mestime = 3000);
    //成功
    void SetSuccess(const QString &text,const int & mestime = 3000);
    //警告
    void SetWaring(const QString &text,const int & mestime = 3000);
    //提示
    void SetTips(const QString &text,const int & mestime = 3000);
private slots:
    void onTimerStayOut();
private:
    void setText(const QString &color="FFFFFF",const QString &bgcolor = "000000",const int & mestime=3000,const QString &textconst="");
private:
    QWidget *m_myWidget;
    QLabel *m_label;
    QTimer *m_timer;
    Ui::Toast *ui;
};
`

const toast_dotCpp = `
#include <QTimer>
#include <QLabel>
#include <QWidget>
#include <QPropertyAnimation>
#include <QPainter>
#include <QScreen>
#include <QHBoxLayout>
#include <QDesktopWidget>
#include <QApplication>

QString StringToRGBA(const QString &color);

Toast::Toast(QObject *parent) : QObject(parent)
{
    m_myWidget = new QWidget;
    m_myWidget->setFixedHeight(60);
    m_label = new QLabel;
    m_label->setFixedHeight(30);
    m_label->move(0,0);
    QFont ft;
    ft.setPointSize(10);
    m_label->setFont(ft);
    m_label->setAlignment(Qt::AlignCenter);
    m_label->setStyleSheet("color:white");
    m_myWidget->setStyleSheet("border: none;background-color:black;border-radius:10px");
    QHBoxLayout * la = new QHBoxLayout;
    la->addWidget(m_label);
    la->setContentsMargins(0,0,0,0);
    m_myWidget->setLayout(la);
    m_myWidget->hide();
    m_myWidget->setWindowFlags(Qt::FramelessWindowHint | Qt::Tool | Qt::WindowStaysOnTopHint);
    m_myWidget->setAttribute(Qt::WA_TranslucentBackground,true);
    m_timer = new QTimer();
    m_timer->setInterval(3000);
    connect(m_timer,SIGNAL(timeout()),this,SLOT(onTimerStayOut()));
}

Toast *Toast::Instance()
{
    static Toast instance;
    return &instance;
}

void Toast::setText(const QString &color,const QString &bgcolor,const int & mestime,const QString &text){
    QApplication::beep();
    QFontMetrics fm(m_label->font());
    int width = fm.boundingRect(text).width() + 30;
    m_myWidget->setFixedWidth(width);
    m_label->setFixedWidth(width);
    m_label->setText(text);
    QString style = QString("color:").append(StringToRGBA(color));
    m_label->setStyleSheet(style);

    m_myWidget->setStyleSheet(QString("border: none;border-radius:10px;")
                            .append("background-color:").append(StringToRGBA(bgcolor)));

    QDesktopWidget *pDesk = QApplication::desktop();
    m_myWidget->move((pDesk->width() - m_myWidget->width()) / 2, (pDesk->height() - m_myWidget->height()) / 2);
    m_myWidget->show();
    m_timer->setInterval(mestime);
    m_timer->stop();
    m_timer->start();
}

void Toast::SetError(const QString &text,const int & mestime){
     setText("FFFFFF","FF0000",mestime,text);
}

void Toast::SetSuccess(const QString &text,const int & mestime){
     setText("000000","00FF00",mestime,text);
}

void Toast::SetWaring(const QString &text,const int & mestime){
     setText("FF0000","FFFF00",mestime,text);
}

void Toast::SetTips(const QString &text,const int & mestime){
     setText("FFFFFF","0080FF",mestime,text);
}

QString StringToRGBA(const QString &color){
    int r = color.mid(0,2).toInt(nullptr,16);
    int g = color.mid(2,2).toInt(nullptr,16);
    int b = color.mid(4,2).toInt(nullptr,16);
    int a = color.length()>=8?color.mid(6,2).toInt(nullptr,16):QString("FF").toInt(nullptr,16);
    QString style = QString("rgba(").append(QString::number(r)).append(",")
            .append(QString::number(g)).append(",")
            .append(QString::number(b)).append(",")
            .append(QString::number(a))
            .append(");");
    return style;
}

void Toast::onTimerStayOut()
{
    m_timer->stop();
    m_myWidget->hide();
}
`
