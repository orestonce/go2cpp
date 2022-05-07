package go2cpp

const runOnUiThread_dotH = `
#include <QObject>
#include <QVector>
#include <QThreadPool>
#include <QMutex>
#include <QMutexLocker>
#include <functional>

class RunOnUiThread : public QObject
{
    Q_OBJECT
public:
    explicit RunOnUiThread(QObject *parent = nullptr);
    virtual ~RunOnUiThread();

    void AddRunFnOn_OtherThread(std::function<void()> fn);
    // !!!注意,fn可能被调用,也可能由于RunOnUiThread被析构不被调用
    // 依赖于在fn里delete回收内存, 关闭文件等操作可能造成内存泄露
    void AddRunFnOn_UiThread(std::function<void ()> fn);
signals:
    void signal_newFn();
private slots:
    void slot_newFn();
private:
    bool m_done;
    QVector<std::function<void()>> m_funcList;
    QMutex m_Mutex;
    QThreadPool m_pool;
};
`

const runOnUiThread_dotCpp = `
// Qt:
#include <QMutexLocker>
#include <QtConcurrent/QtConcurrent>

RunOnUiThread::RunOnUiThread(QObject *parent) : QObject(parent), m_done(false)
{
    // 用signal里的Qt::QueuedConnection 将多线程里面的函数转化到ui线程里调用
    connect(this, SIGNAL(signal_newFn()), this, SLOT(slot_newFn()), Qt::QueuedConnection);
}

RunOnUiThread::~RunOnUiThread()
{
    {
        QMutexLocker lk(&this->m_Mutex);
        this->m_done = true;
        this->m_funcList.clear();
    }
    this->m_pool.clear();
    this->m_pool.waitForDone();
}

void RunOnUiThread::AddRunFnOn_OtherThread(std::function<void ()> fn)
{
    QMutexLocker lk(&this->m_Mutex);
    if (this->m_done) {
        return;
    }
    QtConcurrent::run(&this->m_pool, fn);
}

void RunOnUiThread::slot_newFn()
{
    QVector<std::function<void ()>> fn_vector;
    {
        QMutexLocker lk(&this->m_Mutex);
        if (this->m_funcList.empty() || this->m_done) {
            return;
        }
        fn_vector.swap(this->m_funcList);
    }

    for(std::function<void ()>& fn : fn_vector)
    {
        bool v_done = false;
        {
            QMutexLocker lk(&this->m_Mutex);
            v_done = this->m_done;
        }
        if (v_done) { // 快速结束
            return;
        }
        fn();
    }
}

void RunOnUiThread::AddRunFnOn_UiThread(std::function<void ()> fn)
{
    {
        QMutexLocker lk(&this->m_Mutex);
        if (this->m_done) {
            return;
        }
        m_funcList.push_back(fn);
    }
    emit this->signal_newFn();
}
`
