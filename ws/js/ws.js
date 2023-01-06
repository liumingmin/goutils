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
            let handler = this.msgHandler[msgPack.protocolId];
            handler?.(this.ws, msgPack.dataBuffer);
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
        this.ws.send(this.packMsg(protocolId, data));
    }

    unpackMsg(buffer) {
        let packetHeadFlag = buffer.slice(0, 2);
        const packetLenDv = new DataView(buffer.slice(2, 6));
        const packetLength = packetLenDv.getUint32(0, /* little endian data */ true);
        const protocolIdDv = new DataView(buffer.slice(6, 10));
        const protocolId = protocolIdDv.getUint32(0, /* little endian data */ true);
        let dataBuffer = buffer.slice(10);
        return {
            packetHeadFlag,
            packetLength,
            protocolId,
            dataBuffer
        };
    }

    packMsg(protocolId, buffer) {
        let dataArray = new Uint8Array(buffer);

        let packetLength = new Uint8Array(4);
        new DataView(packetLength.buffer).setUint32(0, dataArray.byteLength+4, true /* littleEndian */);

        let protocolIdArray = new Uint8Array(4);
        new DataView(protocolIdArray.buffer).setUint32(0, protocolId, true /* littleEndian */);

        let packet = new Uint8Array([...this.packetHeadFlag, ...packetLength, ...protocolIdArray, ...dataArray]);
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