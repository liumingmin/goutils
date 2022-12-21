#include "QWsConnection.h"
#include <QtCore>
#include <QtGui>

int main(int argc, char *argv[])
{
    QCoreApplication app(argc, argv);

    QWsConnection conn;
    conn.Connect("ws://127.0.0.1:8003/join?uid=y10000");
    app.exec();
}