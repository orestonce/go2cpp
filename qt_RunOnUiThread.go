package go2cpp

const runOnUiThread_dotH = `#include <vector>
#include <functional>
#include <QMutex>
#include <QObject>
#include <QThreadPool>
#include <QMutexLocker>

class RunOnUiThread : public QObject
{
    Q_OBJECT
public:
    virtual ~RunOnUiThread();

    void AddRunFnOn_OtherThread(std::function<void()> fn);
    // !!!注意,fn可能被调用,也可能由于RunOnUiThread被析构不被调用
    // 依赖于在fn里delete回收内存, 关闭文件等操作可能造成内存泄露
    void AddRunFnOn_UiThread(std::function<void ()> fn);
    bool IsDone();
private slots:
    void slot_newFn();
private:
    bool m_done = false;
    std::vector<std::function<void()>> m_funcList;
    QMutex m_mutex;
    QThreadPool m_pool;
};
`

const runOnUiThread_dotCpp = `// Qt:
#include <QMutexLocker>
#include <QtConcurrent/QtConcurrent>

RunOnUiThread::~RunOnUiThread()
{
    {
        QMutexLocker lk(&m_mutex);
        m_done = true;
        m_funcList.clear();
    }
    m_pool.clear();
    m_pool.waitForDone();
}

void RunOnUiThread::AddRunFnOn_OtherThread(std::function<void ()> fn)
{
    QMutexLocker lk(&m_mutex);
    if (m_done) {
        return;
    }
    QtConcurrent::run(&m_pool, fn);
}

void RunOnUiThread::slot_newFn()
{
    std::vector<std::function<void ()>> funcList;
    {
        QMutexLocker lk(&m_mutex);
        funcList.swap(m_funcList);
    }

    for(std::function<void ()>& fn : funcList)
    {
        if (IsDone()) { // 快速结束
            return;
        }
        fn();
    }
}

void RunOnUiThread::AddRunFnOn_UiThread(std::function<void ()> fn)
{
    {
        QMutexLocker lk(&m_mutex);
        if (m_done) {
            return;
        }
        m_funcList.push_back(fn);
    }

    QMetaObject::invokeMethod(this, "slot_newFn", Qt::QueuedConnection);
}

bool RunOnUiThread::IsDone()
{
    QMutexLocker lk(&m_mutex);
    return m_done;
}
`
