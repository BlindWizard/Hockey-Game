import { Position } from "../engine/position";

const MessageType = {
    Hello: "HELLO",
    Ping: "PING",
    Pong: "PONG",
    Online: "ONLINE",
    Queue: "QUEUE",
    UnQueue: "UNQUEUE",
    Game: "GAME",
    ExitGame: "EXITGAME",
    World: "WORLD",
    PlayerAction: "PLAYERACTION",
};

const ServerMessages = {
    Hello: (playerId) => {
        return {
            playerId,
            type: MessageType.Hello,
        };
    },
    Ping: (timestamp, latency) => {
        return {
            timestamp,
            latency,
            type: MessageType.Ping,
        };
    },
    Online: (count) => {
        return {
            count,
            type: MessageType.Online,
        };
    },
    Game: (roomId) => {
        return {
            roomId,
            type: MessageType.Game,
        } 
    },
    World: (playerPosition, opponentPosition, puckPosition, countA, countB) => {
        return {
            playerPosition,
            opponentPosition,
            puckPosition,
            countA,
            countB,
            type: MessageType.World,
        };
    },
    ExitGame: (reason) => {
        return {
            reason,
            type:MessageType.ExitGame,
        };
    },
    parse(body) {
        let parts = body.split(':');
        if (parts.length <= 0) {
            return null;
        }

        let type = parts.shift();

        switch (type) {
            case MessageType.Hello:
                return this.Hello(parts.shift());
            case MessageType.Ping:
                return this.Ping(parts.shift(), parts.shift());
            case MessageType.Online: 
                return this.Online(parts.shift());
            case MessageType.Game: 
                return this.Game(parts.shift());
            case MessageType.ExitGame:
                return this.ExitGame(parts.shift())
            case MessageType.World:
                const playerX = parts.shift();
                const playerY = parts.shift();
                const opponentX = parts.shift();
                const opponentY = parts.shift();
                const puckX = parts.shift();
                const puckY = parts.shift();
                const countA = parts.shift();
                const countB = parts.shift();

                return this.World(
                    new Position(Number(playerX), Number(playerY)), 
                    new Position(Number(opponentX), Number(opponentY)), 
                    new Position(Number(puckX), Number(puckY)),
                    countA,
                    countB
                );
            default: return null
        }
    }
};

const ClientMessages = {
    Hello: (playerId) => {
        return {
            playerId,
            stringify: () => {
                return `${MessageType.Hello}:${playerId || ''}`;
            }
        };
    },
    Queue: () => {
        return {
            stringify: () => {
                return `${MessageType.Queue}`;
            }
        };
    },
    UnQueue: () => {
        return {
            stringify: () => {
                return `${MessageType.UnQueue}`;
            }
        };
    },
    Ping: () => {
        return {
            stringify: () => {
                return `${MessageType.Ping}`;
            }
        };
    },
    Pong: (timestamp) => {
        return {
            timestamp,
            stringify: () => {
                return `${MessageType.Pong}:${timestamp}`;
            }
        }; 
    },
    PlayerAction: (position) => {
        return {
            x: position.x,
            y: position.y,
            stringify: () => {
                return `${MessageType.PlayerAction}:${position.x}:${position.y}`;
            }
        };
    },
    ExitGame: (reason) => {
        return {
            reason,
            stringify: () => {
                return `${MessageType.ExitGame}:${reason}`;
            }
        };
    }
}

export {MessageType, ServerMessages, ClientMessages};