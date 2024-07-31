import { ClientMessages } from "../protocol/messages";
import { Constants } from "./constants";
import { Position } from "./position";
import { Utils } from "./utils";

export class Game {
    constructor(graphics, network) {
        this.graphics = graphics;
        this.network = network;
        this.updateTimer = null;
        this.captureTimer = null;
        this.setDefaultPositions();
        this.countA = 0;
        this.countB = 0;
        this.debug = false;
    }

    toggleDebug() {
        this.debug = !this.debug;
    }

    setDefaultPositions() {
        this.playerPosition = new Position(
            Constants.gameWidth / 2, 
            Constants.gameHeight - 2 * Constants.entityRadius,
        );

        this.opponentPosition = new Position(
            Constants.gameWidth / 2, 
            2 * Constants.entityRadius,
        );

        this.puckPosition = new Position(
            Constants.gameWidth / 2,
            Constants.gameHeight / 2
        )
    }

    updatePlayerPosition(e) {
        const rect = this.graphics.canvas.getBoundingClientRect();

        const mousePosition = new Position(
            ((e.clientX - rect.left) / (rect.right - rect.left)) * Constants.canvasWidth,
            ((e.clientY - rect.top) / (rect.bottom - rect.top)) * Constants.canvasHeight,
        );

        const gamePlayerPosition = Utils.convertToGameXY(mousePosition);

        this.playerPosition.x = Math.round(Utils.clamp(
            gamePlayerPosition.x, 
            Constants.entityRadius + Constants.lineWidth, 
            Constants.gameWidth - Constants.entityRadius - Constants.lineWidth
        ));

        this.playerPosition.y = Math.round(Utils.clamp(
            gamePlayerPosition.y, 
            Constants.gameHeight / 2 + Constants.entityRadius + Constants.lineWidth,
             Constants.gameHeight - Constants.entityRadius - Constants.lineWidth
        ));
    }

    run() {
        window.addEventListener("mouseenter", this.updatePlayerPosition.bind(this));
        window.addEventListener("mousemove", this.updatePlayerPosition.bind(this));
        this.captureTimer = setInterval(this.capturePlayerInput.bind(this), 20);
        this.updateTimer = setInterval(this.updateWorld.bind(this), 20);
    }

    exit(reason) {
        clearInterval(this.captureTimer);
        clearInterval(this.updateTimer);
        this.network.send(ClientMessages.ExitGame(reason));
    }

    capturePlayerInput() {
        this.network.send(ClientMessages.PlayerAction(this.playerPosition));
    }

    receiveWorld(opponentPosition, puckPosition, countA, countB) {
        this.opponentPosition = opponentPosition;
        this.puckPosition = puckPosition;
        this.countA = countA;
        this.countB = countB;
    }

    updateWorld() {
        this.graphics.clear();
        this.graphics.drawField();
        this.graphics.drawPlayer(this.playerPosition);
        this.graphics.drawPlayer(this.opponentPosition);
        this.graphics.drawPuck(this.puckPosition);
        this.graphics.drawGoals(this.countA, this.countB);

        if (this.debug) {
            this.drawDebug();
        }
    }

    drawDebug() {
        this.graphics.drawText(
            `A:${this.playerPosition.x}:${this.playerPosition.y}`,
            Utils.convertToCanvasXY(new Position(20, Constants.gameHeight - Constants.fontSize)),
            Constants.secondaryFont
        );

        this.graphics.drawText(
            `B:${this.opponentPosition.x}:${this.opponentPosition.y}`,
            Utils.convertToCanvasXY(new Position(20, Constants.fontSize)),
            Constants.secondaryFont
        );

        this.graphics.drawText(
            `Puck:${this.puckPosition.x}:${this.puckPosition.y}`, 
            Utils.convertToCanvasXY(new Position(20, Constants.gameHeight / 2 + Constants.fontSize / 2)),
            Constants.secondaryFont,
        );
    }
}