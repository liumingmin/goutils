class WsConnection {
    constructor() {
        this.ws = null;
        this.connected = false;
        this.msgHandler = {};
        this.establishHandler = null;
        this.errHandler = null;
        this.closeHandler = null;
        this.displacedHandler = null;
        this.packetHeadFlag = new Uint8Array([254, 238]);

        this.msgHandler[proto.ws.P_BASE.S2C_ERR_DISPLACE] = (ws, buffer) => { this.onDisplaced(ws, buffer); };
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

    setDisplacedHandler(displacedHandler) {
        this.displacedHandler = displacedHandler;
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
            
            this.establishHandler?.(this.ws, e);
        };
        this.ws.onerror = (error) => {
            this.errHandler?.(this.ws, error);
        };

        this.ws.onmessage = (e) => {
            let msgPack = this.unpackMsg(e.data);
            //console.log(msgPack);
            let wsMessage = proto.ws.P_MESSAGE.deserializeBinary(msgPack.dataBuffer);
            let handler = this.msgHandler[wsMessage.getProtocolId()];
            handler?.(this.ws, wsMessage.getData());
        };
        this.ws.onclose = (e) => {
            this.closeHandler?.(this.ws, e);
            this.connected = false;
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

    onDisplaced(ws, buffer) {
        let displacedMsg = proto.ws.P_DISPLACE.deserializeBinary(buffer);
        let oldIp = new TextDecoder().decode(displacedMsg.getOldIp());
        let newIp = new TextDecoder().decode(displacedMsg.getNewIp());
        console.log(oldIp, " displaced by ", newIp, " at ", displacedMsg.getTs().toString());
        this.displacedHandler?.(this.ws, oldIp, newIp, displacedMsg.getTs());
    }
}