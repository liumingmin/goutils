#include "QWsConnection.h"
#include <QtNetwork/QAbstractSocket>
#include "msg.pb.h"

#include "QWsSignal.h"

BYTE QWsConnection::m_packetHeadFlag[2] = {254, 239};

QWsConnection::QWsConnection(const QString &url, const HttpHeaders &mapHeaders, uint32_t retryInterval, QObject *parent)
    : QObject(parent), m_pWs(nullptr), m_bConnected(false), m_strUrl(url), m_nRetryInterval(retryInterval), m_nSnCounter(0), m_mapHeaders(mapHeaders)
{
    m_pWs = new QWebSocket;
    m_pWs->setParent(this);

    connect(m_pWs, &QWebSocket::connected, this, [this]()
            {
        m_connDeadline = std::chrono::steady_clock::now() + std::chrono::seconds(m_nTimeWait);
        m_bConnected = true;
        if (m_establishHandler)
        {
            m_establishHandler(m_pWs);
        } });

    connect(m_pWs, &QWebSocket::disconnected, this, [this]()
            {
        m_bConnected = false;
        if (m_closeHandler)
        {
            m_closeHandler(m_pWs);
        } });

    if (m_nRetryInterval > 0)
    {
        m_pReconnTimer = new QTimer(this);
        m_pReconnTimer->setInterval(m_nRetryInterval);
        connect(m_pReconnTimer, &QTimer::timeout, this, [this]()
                {
            if (!m_bConnected)
                Connect(); });
        m_pReconnTimer->start();
    }

    connect(m_pWs, QOverload<QAbstractSocket::SocketError>::of(&QWebSocket::error), this,
            [this](QAbstractSocket::SocketError err)
            {
                if (m_errHandler)
                {
                    m_errHandler(m_pWs, err, "");
                }
            });

    connect(m_pWs, &QWebSocket::binaryMessageReceived, this, [this](const QByteArray &message)
            {
        const auto msgPack = _UnpackMsg(message);

        bool bIsValidMsg = _CheckMsgPackErr(msgPack);
        if (!bIsValidMsg)
            return;

        if (msgPack.sn > 0)
        {
            std::promise<QByteArray> promise;
            bool exists = false;

            {
                std::lock_guard<std::mutex> lk(m_mapSnPromiseMutex);
                auto iter = m_mapSnPromise.find(msgPack.sn);
                if (iter != m_mapSnPromise.end())
                {
                    promise = std::move(iter->second);
                    m_mapSnPromise.erase(msgPack.sn);
                    exists = true;
                }
            }

            if (exists)
                promise.set_value(msgPack.dataBuffer);
        }

        bool exist = m_mapMsgHandler.contains(msgPack.protocolId);
        if (exist)
        {
            // send message receive confirm
            SendResponseMsg(msgPack.protocolId, msgPack.sn, "");
            const auto handler = m_mapMsgHandler[msgPack.protocolId];
            handler(m_pWs, msgPack.dataBuffer);
        }

        if (msgPack.sn == 0 && !exist && m_errHandler)
        {
            m_errHandler(m_pWs, QAbstractSocket::SocketError::UnknownSocketError, "protocolId's handler not found");
        } });

    m_pPingTimer = new QTimer(this);
    m_pPingTimer->setInterval(24000);
    connect(m_pPingTimer, &QTimer::timeout, this, [this]()
            {
        if (m_pWs->isValid())
            m_pWs->ping(); });

    m_nTimeWait = 60;
    m_pDeadlineTimer = new QTimer(this);
    m_pDeadlineTimer->setInterval(1000);
    m_connDeadline = std::chrono::steady_clock::now() + std::chrono::seconds(m_nTimeWait);
    connect(m_pDeadlineTimer, &QTimer::timeout, this, [this]()
            {
        if(m_pPingTimer->isActive() &&  m_connDeadline < std::chrono::steady_clock::now())
        {
            _Reset();
        } });
    m_pDeadlineTimer->start();

    connect(m_pWs, &QWebSocket::pong, this, [this](quint64 elapsed, const QByteArray &)
            { m_connDeadline = std::chrono::steady_clock::now() + std::chrono::seconds(m_nTimeWait); });
    // m_mapMsgHandler[ws::P_BASE::s2c_err_displace] = [this](QWebSocket* ws, const QByteArray& data)
    //{
    //     _OnDisplaced(ws, data);
    // };
}

QWsConnection::~QWsConnection()
{
    if (m_pReconnTimer)
        m_pReconnTimer->stop();

    m_pDeadlineTimer->stop();
    m_pPingTimer->stop();

    Close();
}

void QWsConnection::AcceptAllSelfSignCert()
{
    QSslConfiguration sslConfiguration = m_pWs->sslConfiguration();
    sslConfiguration.setPeerVerifyMode(QSslSocket::VerifyNone);
    m_pWs->setSslConfiguration(sslConfiguration);
    m_pWs->ignoreSslErrors();
}

void QWsConnection::AcceptSelfSignCert(const QString &caCertPath)
{
    QList<QSslCertificate> certs = QSslCertificate::fromPath(caCertPath);
    QSslConfiguration sslConfiguration = m_pWs->sslConfiguration();
    sslConfiguration.addCaCertificates(certs);

    QList<QSslError> expectedSslErrors;
    for (auto &cert : certs)
    {
        QSslError ignoreError(QSslError::InvalidPurpose, cert);
        expectedSslErrors.append(ignoreError);
    }
    m_pWs->ignoreSslErrors(expectedSslErrors);
}

void QWsConnection::RegisterMsgHandler(uint32_t protocolId, MsgHandler handler)
{
    m_mapMsgHandler.insert(protocolId, handler);
}

QWsConnection::State QWsConnection::SendMsg(uint32_t protocolId, const QByteArray &data)
{
    auto sent = m_pWs->sendBinaryMessage(_PackMsg(protocolId, 0, data));
    if (sent > 0)
    {
        return STATE_OK;
    }
    return STATE_SEND_FAILED;
}

QWsConnection::State QWsConnection::SendRequestMsg(uint32_t protocolId, const QByteArray &request,
                                                   uint32_t nTimeoutMs, QByteArray &response)
{
    State nRetCode = STATE_SEND_FAILED;

    uint32_t sn = _GetNextSn();

    std::promise<QByteArray> promise;
    std::future<QByteArray> future = promise.get_future();

    {
        std::lock_guard<std::mutex> lk(m_mapSnPromiseMutex);
        m_mapSnPromise.emplace(sn, std::move(promise));
    }

    auto sent = m_pWs->sendBinaryMessage(_PackMsg(protocolId, sn, request));
    if (sent <= 0)
    {
        nRetCode = STATE_SEND_FAILED;
        goto Exit0;
    }

    {
        QWS_DECLARE_SIGNAL(signal);
        auto waitFuture = std::async(std::launch::async, [&signal, &future, &nRetCode, &response, nTimeoutMs]
        {
            QWS_GUARD_SIGNAL(signal);

            auto status = future.wait_for(std::chrono::milliseconds(nTimeoutMs));
            if (status == std::future_status::timeout)
            {
                nRetCode = STATE_RESP_TIMEOUT;
            }
            else if (status == std::future_status::ready)
            {
                nRetCode = STATE_OK;
                response = future.get();
            }
        });

        QWS_WAIT_SIGNAL(signal);
    }

Exit0:
{
    std::lock_guard<std::mutex> lk(m_mapSnPromiseMutex);
    m_mapSnPromise.erase(sn);
}
    return nRetCode;
}

QWsConnection::State QWsConnection::SendResponseMsg(uint32_t protocolId, uint32_t reqSn, const QByteArray &data)
{
    auto sent = m_pWs->sendBinaryMessage(_PackMsg(protocolId, reqSn, data));
    if (sent > 0)
    {
        return STATE_OK;
    }
    return STATE_SEND_FAILED;
}

void QWsConnection::SetPingInterval(uint64_t nInterval)
{
    m_pPingTimer->setInterval(nInterval);
}

void QWsConnection::SetPing(bool bEnable)
{
    if (bEnable)
    {
        if (m_pPingTimer->isActive())
        {
            m_pPingTimer->stop();
        }

        m_pPingTimer->start();
    }
    else
    {
        m_pPingTimer->stop();
    }
}

void QWsConnection::Connect()
{
    _Reset();

    QNetworkRequest request(m_strUrl);

    for (auto iter = m_mapHeaders.begin(); iter != m_mapHeaders.end(); ++iter)
    {
        request.setRawHeader(QByteArray::fromStdString(iter->first), QByteArray::fromStdString(iter->second));
    }

    m_connDeadline = std::chrono::steady_clock::now() + std::chrono::seconds(10);
    m_pWs->open(request);
}

void QWsConnection::Close()
{
    if (m_pReconnTimer)
        m_pReconnTimer->stop();

    _Reset();
}

void QWsConnection::_Reset()
{
    if (m_pWs != nullptr)
    {
        m_pWs->close();
        m_bConnected = false;
    }

    m_nSnCounter = 0;

    std::lock_guard<std::mutex> lk(m_mapSnPromiseMutex);
    m_mapSnPromise.clear();

    // must clear the old header data
    m_mapHeaders.clear();
}

uint32_t QWsConnection::_GetNextSn()
{
    uint32_t sn = m_nSnCounter.fetch_add(2);
    if (sn == 0)
    {
        sn = m_nSnCounter.fetch_add(2);
    }

    return sn;
}

QByteArray QWsConnection::_PackMsg(uint32_t protocolId, uint32_t sn, const QByteArray &dataBuffer)
{
    BYTE packetLength[4];
    BYTE protocolIdArray[4];
    BYTE snArray[4];

    qToLittleEndian(dataBuffer.size() + 8, &packetLength);
    qToLittleEndian(protocolId, &protocolIdArray);
    qToLittleEndian(sn, &snArray);

    QByteArray byteArray;
    byteArray.append((const char *)m_packetHeadFlag, 2);
    byteArray.append((const char *)packetLength, 4);
    byteArray.append((const char *)protocolIdArray, 4);
    byteArray.append((const char *)snArray, 4);
    byteArray.append(dataBuffer);

    return byteArray;
}

QWsConnection::innerMsgPack QWsConnection::_UnpackMsg(const QByteArray &rawMsg)
{
    BYTE packetLength[4];
    BYTE protocolIdArray[4];
    BYTE snArray[4];

    innerMsgPack msgPack;

    memcpy(&msgPack.packetHeadFlag, rawMsg.mid(0, 2).data(), 2);
    memcpy(&packetLength, rawMsg.mid(2, 4).data(), 4);
    memcpy(&protocolIdArray, rawMsg.mid(6, 4).data(), 4);
    memcpy(&snArray, rawMsg.mid(10, 4).data(), 4);

    msgPack.packetLength = qFromLittleEndian<quint32>(&packetLength);
    msgPack.protocolId = qFromLittleEndian<quint32>(&protocolIdArray);
    msgPack.sn = qFromLittleEndian<quint32>(&snArray);
    msgPack.dataBuffer = rawMsg.mid(14);
    return msgPack;
}

void QWsConnection::_OnDisplaced(QWebSocket *ws, const QByteArray &msgData)
{
    // ws::P_DISPLACE displacedMsg;
    // bool result = displacedMsg.ParseFromArray(msgData.data(), msgData.size());
    // if (!result)
    //     return;

    // if (m_displacedHandler)
    //{
    //     m_displacedHandler(m_pWs, QString::fromStdString(displacedMsg.old_ip()),
    //                        QString::fromStdString(displacedMsg.new_ip()),
    //                        displacedMsg.ts());
    // }
}

bool QWsConnection::_CheckMsgPackErr(const innerMsgPack &msgPack)
{
    if (msgPack.packetHeadFlag[0] != m_packetHeadFlag[0] || msgPack.packetHeadFlag[1] != m_packetHeadFlag[1])
    {
        if (m_errHandler)
        {
            m_errHandler(m_pWs, QAbstractSocket::SocketError::UnknownSocketError, "packetHeadFlag error");
        }
        return false;
    }

    if (msgPack.packetLength != msgPack.dataBuffer.size() + 8)
    {
        if (m_errHandler)
        {
            m_errHandler(m_pWs, QAbstractSocket::SocketError::UnknownSocketError, "packetLength error");
        }
        return false;
    }

    return true;
}
