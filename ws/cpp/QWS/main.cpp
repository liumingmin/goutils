#include "QWsConnection.h"
#include <QtCore>
#include <QtGui>

int main(int argc, char *argv[])
{
    QCoreApplication app(argc, argv);

    QWsConnection conn("wss://test.com:8003/join?uid=y10000",10000); //hosts 127.0.0.1 <- test.com
    conn.SetEstablishHandler([&](QWebSocket*)
    {
        qDebug() << "connected";
        conn.SendMsg(2, QByteArray::fromStdString("js request"));
    });

    conn.SetCloseHandler([&](QWebSocket*)
    {
        qDebug() << "closed";
    });
    conn.SetErrHandler([&](QWebSocket*, QAbstractSocket::SocketError err)
    {
        qDebug() << "err:" << err;
    });

    conn.SetDisplacedHandler([&](QWebSocket*, QString oldIp, QString newIp, int64_t ts)
    {
        qDebug() << oldIp << " displaced by " << newIp << " at " << ts;
    });

    conn.RegisterMsgHandler(3, [&](QWebSocket*, QByteArray data){
        qDebug() << data;
    });
    conn.AcceptSelfSignCert("xxx.crt");
    conn.Connect();

    app.exec();
}