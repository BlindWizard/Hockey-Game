import { ErrorTypes } from "../errorTypes";
import { ClientMessages, MessageType, ServerMessages } from "../protocol/messages";
import { Constants } from "./constants";

export class Network {

    constructor(errorHandler) {
        this.messageHandlers = new Map();
        this.errorHandler = errorHandler;
    }

    setOnOpen(onOpen) {
        this.onOpen = onOpen;
    }

    setOnHello(onHello) {
        this.messageHandlers.set(MessageType.Hello, onHello);
    }

    setOnPing(onPing) {
        this.messageHandlers.set(MessageType.Ping, onPing);
    }

    setOnGame(onGame) {
        this.messageHandlers.set(MessageType.Game, onGame);
    }

    setOnOnline(onOnline) {
        this.messageHandlers.set(MessageType.Online, onOnline);
    }

    setOnWorld(onWorld) {
        this.messageHandlers.set(MessageType.World, onWorld)
    }

    setOnExitGame(onExitGame) {
        this.messageHandlers.set(MessageType.ExitGame, onExitGame)
    }

    openConnection() {
        this.socket = new WebSocket(this.url());

        this.socket.onopen = () => {
            this.onOpen && this.onOpen();
        }

        this.socket.onmessage = (e) => {
            const data = e.data;
            const message = ServerMessages.parse(data);

            if (!message) {
                return;
            }

            if (message.type === MessageType.Ping) {
                this.send(ClientMessages.Pong(message.timestamp));
            }

            if (this.messageHandlers.has(message.type)) {
                this.messageHandlers.get(message.type)(message);
            }
        }

        this.socket.onerror = () => {
            this.errorHandler(new Error('Connection error', { cause: ErrorTypes.lostConnection }));
        }

        this.socket.onClose = (e) => {
            this.errorHandler(new Error(
                e.wasClean ? 'Server closed connection' : 'Connection to server lost',
                { cause: ErrorTypes.lostConnection }
            ));
        }
    }

    send(message) {
        try {
            if (!this.socket || this.socket.readyState === WebSocket.CONNECTING) {
                throw new Error('You need open connection before messaging');
            }
    
            if (this.socket.readyState === WebSocket.CLOSING || this.socket.readyState === WebSocket.CLOSED) {
                throw new Error('Server closed connection', { cause: ErrorTypes.lostConnection });
            }
            
            this.socket.send(message.stringify());
        } catch (e) {
            this.errorHandler(e);
        }
    }

    runPing() {
        if (this.ping) {
            return;
        }

        this.ping = setInterval(() => {
            try {
                this.send(ClientMessages.Ping());
            } catch (e) {
                clearInterval(this.ping);
                this.errorHandler(e);
            }
        }, 1000);
    }

    url() {
        return (window.location.protocol === "https:" ? "wss://" : "ws://") + window.location.hostname + ':' + Constants.socketPort + Constants.socketPath;
    }
}