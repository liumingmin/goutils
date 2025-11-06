#pragma once
#include <QtCore>
#include <QtWebSockets/QWebSocket>

typedef unsigned char BYTE;

class  QWsConnection : public QObject
{
    Q_OBJECT

    using MsgHandler = std::function<void(QWebSocket*, const QByteArray&)>;
    using EvtHandler = std::function<void(QWebSocket*)>;
    using ErrHandler = std::function<void(QWebSocket*, QAbstractSocket::SocketError, const QString&)>;
    using DisplacedHandler = std::function<void(QWebSocket*, QString, QString, int64_t)>;
    using HttpHeaders = std::map<std::string, std::string>;

public:
    enum State
    {
        STATE_OK = 1,
        STATE_SEND_FAILED = -1,
        STATE_RESP_TIMEOUT = -2
    };

public:
    explicit QWsConnection(const QString& url, const HttpHeaders& mapHeaders = {}, uint32_t retryInterval=0, QObject* parent = nullptr);
    virtual  ~QWsConnection();

    void AcceptAllSelfSignCert();
    void AcceptSelfSignCert(const QString&  caCertPath);

    void RegisterMsgHandler(uint32_t protocolId, MsgHandler handler);
    State SendMsg(uint32_t protocolId, const QByteArray& data);
    State SendRequestMsg(uint32_t protocolId, const QByteArray& request, uint32_t nTimeoutMs, QByteArray& response);
    State SendResponseMsg(uint32_t protocolId, uint32_t reqSn, const QByteArray& data);

    inline bool IsConnected() { return m_bConnected; }
    inline void SetUrl(const QString& url) { m_strUrl = url; }
    inline void SetRetryInterval(uint32_t retryInterval) { m_nRetryInterval = retryInterval; }

    inline void SetEstablishHandler(EvtHandler establishHandler) { m_establishHandler = establishHandler; }
    inline void SetCloseHandler(EvtHandler closeHandler) { m_closeHandler = closeHandler; }
    inline void SetErrHandler(ErrHandler errHandler) { m_errHandler = errHandler; }
    inline void SetDisplacedHandler(DisplacedHandler displacedHandler) { m_displacedHandler = displacedHandler; }

    inline void SetTimeWait(uint64_t timeWait) { m_nTimeWait = timeWait; };
    void SetPingInterval(uint64_t nInterval);
    void SetPing(bool bEnable);

public slots:
    void Connect();
    void Close();

protected:
    struct innerMsgPack
    {
        BYTE        packetHeadFlag[2];
        uint32_t    packetLength;
        uint32_t    protocolId;
        uint32_t    sn;
        QByteArray  dataBuffer;

        innerMsgPack():packetLength(0), protocolId(0), sn(0)
        {
            packetHeadFlag[2] = { 0 };
        }
    };

    QByteArray _PackMsg(uint32_t protocolId, uint32_t sn, const QByteArray& dataBuffer);
    innerMsgPack _UnpackMsg(const QByteArray& rawMsg);
    void _OnDisplaced(QWebSocket* ws, const QByteArray& msgData);
    bool _CheckMsgPackErr(const innerMsgPack& msgPack);
    void _Reset();
    uint32_t _GetNextSn();


private:
    QWebSocket*                                     m_pWs = NULL;
    QTimer*                                         m_pReconnTimer = NULL;
    QTimer*                                         m_pPingTimer = NULL;
    QTimer*                                         m_pDeadlineTimer = NULL;
    uint64_t                                        m_nTimeWait = 0;
    std::chrono::time_point<std::chrono::steady_clock>  m_connDeadline;
    bool                                            m_bConnected = false;
    QString                                         m_strUrl;
    uint32_t                                        m_nRetryInterval = 0;
    HttpHeaders                                     m_mapHeaders;

    QHash<uint32_t, MsgHandler>                     m_mapMsgHandler;
    EvtHandler                                      m_establishHandler;
    EvtHandler                                      m_closeHandler;
    ErrHandler                                      m_errHandler;
    DisplacedHandler                                m_displacedHandler;

    std::atomic_uint32_t                            m_nSnCounter = 0;   //sn counter, atomic
    std::map<uint32_t, std::promise<QByteArray>>    m_mapSnPromise;     //sn channel store map
    std::mutex                                      m_mapSnPromiseMutex;//sn channel store map lock

    static BYTE                                     m_packetHeadFlag[2];
};
