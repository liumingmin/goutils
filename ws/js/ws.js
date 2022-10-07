class WsConnection {
    constructor() {
        this.ws = null;
        this.connected = false;
        this.msgHandler = {};
        this.establishHandler = null;
        this.errHandler = null;
        this.closeHandler = null;
        this.packetHeadFlag = new Uint8Array([254, 238]);
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

    connect(url, retryInterval) {
        if (this.ws !== null) {
            this.ws.close();
            this.connected = false;
        }

        this.ws = new WebSocket(url);

        this.ws.onopen = (e) => {
            this.ws.binaryType = 'arraybuffer';
            this.connected = true;
            if (this.establishHandler !== null) {
                this.establishHandler(this.ws, e);
            }
        };
        this.ws.onerror = (error) => {
            if (this.errHandler !== null) {
                this.errHandler(this.ws, error);
            }
        };

        this.ws.onmessage = (e) => {
            let msgPack = this.unpackMsg(e.data);
            //console.log(msgPack);
            let wsMessage = proto.ws.P_MESSAGE.deserializeBinary(msgPack.dataBuffer);
            let handler = this.msgHandler[wsMessage.getProtocolId()];
            if (handler !== null) {
                handler(this.ws, wsMessage.getData());
            }
        };
        this.ws.onclose = (e) => {
            this.connected = false;
            if (this.closeHandler !== null) {
                this.closeHandler(this.ws, e);
            }
            this.ws = null;

            if (retryInterval === undefined) return;

            setTimeout(() => this.connect(url, retryInterval), retryInterval);
        };
    }

    sendMsg(protocolId, data) {
        let wsMessage = new proto.ws.P_MESSAGE;
        wsMessage.setProtocolId(protocolId);
        wsMessage.setData(data);
        this.ws.send(this.packMsg(wsMessage.serializeBinary()));
    }

    unpackMsg(buffer) {
        let packetHeadFlag = buffer.slice(0, 2);
        const dv = new DataView(buffer.slice(2, 6));
        const packetLength = dv.getUint32(0, /* little endian data */ true);
        let dataBuffer = buffer.slice(6);
        return {
            packetHeadFlag,
            packetLength,
            dataBuffer
        };
    }

    packMsg(buffer) {
        let dataArray = new Uint8Array(buffer);

        let packetLength = new Uint8Array(4);
        new DataView(packetLength.buffer).setUint32(0, buffer.byteLength, true /* littleEndian */);

        let packet = new Uint8Array([...this.packetHeadFlag, ...packetLength, ...dataArray]);
        return packet.buffer;
    }
}