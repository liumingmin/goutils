#pragma once

#include <QObject>

#define QWS_DECLARE_SIGNAL(DeclSignal) QWsSignal DeclSignal;QEventLoop DeclSignal##EvtLoop;QObject::connect(&DeclSignal, &QWsSignal::finished, &DeclSignal##EvtLoop, &QEventLoop::quit,Qt::QueuedConnection)
#define QWS_EMIT_SIGNAL(DeclSignal)  emit DeclSignal.finished()
#define QWS_GUARD_SIGNAL(DeclSignal) QWsSignalGuard DeclSignal##Guard(&DeclSignal)
#define QWS_WAIT_SIGNAL(DeclSignal)  DeclSignal##EvtLoop.exec()

class QWsSignal : public QObject
{
    Q_OBJECT

public:
    QWsSignal(QObject* parent = nullptr) : QObject(parent)
    {
    };

    ~QWsSignal()
    {
    };

Q_SIGNALS:
    void finished();
};

class QWsSignalGuard
{
public:
    QWsSignalGuard(QWsSignal* pSignal) : m_pSignal(pSignal)
    {
    }

    ~QWsSignalGuard()
    {
        emit(*m_pSignal).finished();
    }

private:
    QWsSignal* m_pSignal;
};
