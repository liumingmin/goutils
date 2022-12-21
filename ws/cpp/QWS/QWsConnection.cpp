#include "QWsConnection.h"
#include <QtNetwork/QAbstractSocket>
#include "msg.pb.h"

BYTE QWsConnection::m_packetHeadFlag[2] = { 254, 238 };

QWsConnection::QWsConnection(const QString& url, uint32_t retryInterval, QObject* parent)
    : QObject(parent)
    , m_pWs(nullptr)
    , m_bConnected(false)
    , m_strUrl(url)
    , m_nRetryInterval(retryInterval)
{
    m_pWs = new QWebSocket;
    m_pWs->setParent(this);


    connect(m_pWs, &QWebSocket::connected, [this]() {
        m_bConnected = true;
        if (m_establishHandler)
        {
            m_establishHandler(m_pWs);
        }
    });

    connect(m_pWs, &QWebSocket::disconnected, [this]() {
        m_bConnected = false;
        if(m_closeHandler)
        {
            m_closeHandler(m_pWs);
        }

        if (m_nRetryInterval>0)
            QTimer::singleShot(m_nRetryInterval, this, &QWsConnection::Connect);
    });

    connect(m_pWs, QOverload<QAbstractSocket::SocketError>::of(&QWebSocket::error), [this](QAbstractSocket::SocketError err) {
        if (m_errHandler)
        {
            m_errHandler(m_pWs, err);
        }
    });

    connect(m_pWs, &QWebSocket::binaryMessageReceived, [this](const QByteArray &message) {
        auto msgPack = _UnpackMsg(message);
        //todo valid && errHanlder

        if(m_mapMsgHandler.contains(msgPack.protocolId))
        {
            auto handler = m_mapMsgHandler[msgPack.protocolId];
            handler(m_pWs, msgPack.dataBuffer);
        }
    });

    m_mapMsgHandler[ws::P_BASE::s2c_err_displace] = [this](QWebSocket* ws, const QByteArray& data) {_OnDisplaced(ws, data); };
}


QWsConnection::~QWsConnection()
{

}

void QWsConnection::RegisterMsgHandler(uint32_t protocolId, MsgHandler handler)
{
    m_mapMsgHandler.insert(protocolId, handler);
}

void QWsConnection::SendMsg(uint32_t protocolId, const QByteArray& data)
{
    m_pWs->sendBinaryMessage(_PackMsg(protocolId, data));
}

void QWsConnection::Connect()
{
    if (m_pWs != nullptr) {
        m_pWs->close();
        m_bConnected = false;
    }

    m_pWs->open(QUrl(m_strUrl));
}

QByteArray QWsConnection::_PackMsg(uint32_t protocolId, const QByteArray& dataBuffer)
{
    BYTE        packetLength[4];
    BYTE        protocolIdArray[4];

    qToLittleEndian(dataBuffer.size()+4, &packetLength);
    qToLittleEndian(protocolId, &protocolIdArray);

    QByteArray byteArray;
    byteArray.append((const char*)m_packetHeadFlag,2);
    byteArray.append((const char*)packetLength, 4);
    byteArray.append((const char*)protocolIdArray, 4);
    byteArray.append(dataBuffer);

    return byteArray;
}

QWsConnection::innerMsgPack QWsConnection::_UnpackMsg(const QByteArray& rawMsg)
{
    BYTE        packetLength[4];
    BYTE        protocolIdArray[4];
    innerMsgPack msgPack;

    memcpy(&msgPack.packetHeadFlag, rawMsg.mid(0, 2).data(), 2);
    memcpy(&packetLength, rawMsg.mid(2, 4).data(), 4);
    memcpy(&protocolIdArray, rawMsg.mid(6, 4).data(), 4);

    msgPack.packetLength = qFromLittleEndian<quint32>(&packetLength);
    msgPack.protocolId = qFromLittleEndian<quint32>(&protocolIdArray);
    msgPack.dataBuffer = rawMsg.mid(10);
    return msgPack;
}

void QWsConnection::_OnDisplaced(QWebSocket* ws, QByteArray msgData)
{
    ws::P_DISPLACE displacedMsg;
    bool result = displacedMsg.ParseFromArray(msgData.data(), msgData.size());
    if (!result)
        return;

    if(m_displacedHandler)
    {
        m_displacedHandler(m_pWs, QString::fromStdString(displacedMsg.old_ip()), 
            QString::fromStdString(displacedMsg.new_ip()), 
            displacedMsg.ts());
    }
}
