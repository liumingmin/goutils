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
}