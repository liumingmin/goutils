class WsConnection {
    constructor() {
        this.ws = null;
        this.connected = false;
        this.msgHandler = {};
        this.establishHandler = null;
        this.errHandler = null;
        this.closeHandler = null;
    }

    registerMsgHandler(protocolId, handler) {
        this.msgHandler[protocolId] = handler;
    }

    setEstablishHandler(establishHandler) {
        this.establishHandler = establishHandler;
    }

    setErrHandler(errHandler) {
        this.errHandler = errHandler;
    }

    setCloseHandler(closeHandler) {
        this.closeHandler = closeHandler;
    }

    connect(url) {
        this.ws = new WebSocket(url);

        this.ws.onopen = () => {
            this.ws.binaryType = 'arraybuffer';
            this.connected = true;
            if (this.establishHandler) {
                this.establishHandler(this.ws);
            }
        };
        this.ws.onerror = (error) => {
            if (this.errHandler) {
                this.errHandler(this.ws, error);
            }
        };

        this.ws.onmessage = (e) => {
            let wsMessage = proto.ws.P_MESSAGE.deserializeBinary(e.data);
            let handler = this.msgHandler[wsMessage.getProtocolId()];
            if (handler) {
                handler(this.ws, wsMessage.getData());
            }
        };
        this.ws.onclose = (e) => {
            this.connected = false;
            if (this.closeHandler) {
                this.closeHandler(this.ws, e);
            }
        };
    }

    sendMsg(protocolId, data) {
        let wsMessage = new proto.ws.P_MESSAGE;
        wsMessage.setProtocolId(protocolId);
        wsMessage.setData(data);
        this.ws.send(wsMessage.serializeBinary());
    }

    stringToArrayBuffer(s) {
        var i = s.length;
        var n = 0;
        var ba = new Array()
        for (var j = 0; j < i;) {
            var c = s.codePointAt(j);
            if (c < 128) {
                ba[n++] = c;
                j++;
            } else if ((c > 127) && (c < 2048)) {
                ba[n++] = (c >> 6) | 192;
                ba[n++] = (c & 63) | 128;
                j++;
            } else if ((c > 2047) && (c < 65536)) {
                ba[n++] = (c >> 12) | 224;
                ba[n++] = ((c >> 6) & 63) | 128;
                ba[n++] = (c & 63) | 128;
                j++;
            } else {
                ba[n++] = (c >> 18) | 240;
                ba[n++] = ((c >> 12) & 63) | 128;
                ba[n++] = ((c >> 6) & 63) | 128;
                ba[n++] = (c & 63) | 128;
                j += 2;
            }
        }
        return new Uint8Array(ba);
    }
}