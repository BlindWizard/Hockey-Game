import { ClientMessages } from "../protocol/messages";
import { Constants } from "./constants";
import { Position } from "./position";
import { Utils } from "./utils";

const FRAME_TIME = 20;
const LERP_STEPS = 5;

export class Game {
    constructor(graphics, network) {
        this.graphics = graphics;
        this.network = network;
        this.updateTimer = null;
        this.captureTimer = null;
        this.interpolateCounter = 0;
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

        this.opponentPositionPrev = structuredClone(this.opponentPosition);
        this.puckPositionPrev = structuredClone(this.puckPosition);

        this.opponentPositionInter = structuredClone(this.opponentPosition);
        this.puckPositionInter = structuredClone(this.puckPosition);
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
        this.captureTimer = setInterval(this.capturePlayerInput.bind(this), FRAME_TIME);
        this.updateTimer = setInterval(this.updateWorld.bind(this), FRAME_TIME / LERP_STEPS);
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
        this.interpolateCounter = 0;

        this.opponentPositionPrev.x = this.opponentPosition.x;
        this.opponentPositionPrev.y = this.opponentPosition.y;
        this.puckPositionPrev.x = this.puckPosition.x;
        this.puckPositionPrev.y = this.puckPosition.y;

        this.opponentPosition.x = opponentPosition.x;
        this.opponentPosition.y = opponentPosition.y;

        this.puckPosition.x = puckPosition.x;
        this.puckPosition.y = puckPosition.y;

        this.countA = countA;
        this.countB = countB;
    }

    interpolatePositions() {
        if (this.interpolateCounter > LERP_STEPS) {
            return;
        }

        this.interpolateCounter++;
        let stepSize = this.interpolateCounter / LERP_STEPS;

        this.opponentPositionInter.x = Utils.lerp(this.opponentPositionPrev.x, this.opponentPosition.x, stepSize);
        this.opponentPositionInter.y = Utils.lerp(this.opponentPositionPrev.y, this.opponentPosition.y, stepSize);

        this.puckPositionInter.x = Utils.lerp(this.puckPositionPrev.x, this.puckPosition.x, stepSize);
        this.puckPositionInter.y = Utils.lerp(this.puckPositionPrev.y, this.puckPosition.y, stepSize);
    }

    updateWorld() {
        this.interpolatePositions();

        this.graphics.clear();
        this.graphics.drawField();
        this.graphics.drawPlayer(this.playerPosition);
        this.graphics.drawPlayer(this.opponentPositionInter);
        this.graphics.drawPuck(this.puckPositionInter);
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
            `B:${this.opponentPositionInter.x}:${this.opponentPositionInter.y}`,
            Utils.convertToCanvasXY(new Position(20, Constants.fontSize)),
            Constants.secondaryFont
        );

        this.graphics.drawText(
            `Puck:${this.puckPositionInter.x}:${this.puckPositionInter.y}`, 
            Utils.convertToCanvasXY(new Position(20, Constants.gameHeight / 2 + Constants.fontSize / 2)),
            Constants.secondaryFont,
        );
    }
}