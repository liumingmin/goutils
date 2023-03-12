#include "QWsConnection.h"
#include <QtCore>
#include <QtGui>

enum ProtocolEnum
{
    C2S_REQ = 2,
    S2C_RESP = 3,

    C2S_REQ_TIMEOUT = 4,
    S2C_RESP_TIMEOUT = 5
};

int main(int argc, char* argv[])
{
    QCoreApplication app(argc, argv);

    std::future<void> asyncFuture;

    QWsConnection conn("ws://127.0.0.1:8003/join?uid=y10000", 10000); //hosts 127.0.0.1 <- test.com
    conn.SetEstablishHandler([&](QWebSocket*)
    {
        qDebug() << "connected";

        int bRetCode = 0;
        bRetCode = conn.SendMsg(C2S_REQ, QByteArray::fromStdString("cpp request"));
        qDebug() << "SendMsg retCode: " << bRetCode;

        for (int i = 0; i < 3; i++)
        {
            QByteArray response;
            bRetCode = conn.SendRequestMsg(C2S_REQ, QByteArray::fromStdString("cpp rpc request"), 2000, response);
            qDebug() << "SendRequestMsg retCode: " << bRetCode << " resp:" << QString(response);
        }

        QByteArray responseTimout;
        bRetCode = conn.SendRequestMsg(C2S_REQ_TIMEOUT, QByteArray::fromStdString("cpp rpc  request timeout response"),
                                       2000, responseTimout);
        qDebug() << "SendRequestMsgTimeout retCode: " << bRetCode << " resp:" << QString(responseTimout);
    });

    conn.SetCloseHandler([&](QWebSocket*)
    {
        qDebug() << "closed";
    });
    conn.SetErrHandler([&](QWebSocket*, QAbstractSocket::SocketError err, const QString& errMsg)
    {
        qDebug() << "err:" << err << " " << errMsg;
    });

    conn.SetDisplacedHandler([&](QWebSocket*, QString oldIp, QString newIp, int64_t ts)
    {
        qDebug() << oldIp << " displaced by " << newIp << " at " << ts;
    });

    conn.RegisterMsgHandler(3, [](QWebSocket*, const QByteArray& data)
    {
        qDebug() << data;
    });
    conn.AcceptSelfSignCert("xxx.crt");
    conn.AcceptAllSelfSignCert();
    conn.Connect();

    app.exec();
}
