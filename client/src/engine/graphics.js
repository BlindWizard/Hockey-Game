import { Utils } from "./utils";
import { Constants } from "./constants";
import { Position } from "./position";
import { ErrorTypes } from "../errorTypes";

export class Graphics {
    constructor(canvas) {
        this.canvas = canvas;
        this.context = this.canvas.getContext("2d");
        if (!this.context) {
            throw new Error("Your browser is kinda old and does not support HTML5. Please update.", { cause: ErrorTypes.badBrowser });
        }

        this.canvas.width = Constants.canvasWidth;
        this.canvas.height = Constants.canvasHeight;
        this.context.lineWidth = Constants.lineWidth;
        this.context.strokeStyle = Constants.lineStyle;
        this.context.fillStyle = Constants.lineStyle;
        this.context.font = Constants.font;
    }

    clear() {
        this.context.clearRect(0, 0, Constants.canvasWidth, Constants.canvasHeight);
    }

    drawField() {
        this.context.beginPath();
        const leftTop = Utils.convertToCanvasXY(new Position(0, 0));
        const leftTopGate = Utils.convertToCanvasXY(new Position(Constants.gameWidth / 2 - Constants.gatesWidth / 2, 0));
        const rightTopGate = Utils.convertToCanvasXY(new Position(Constants.gameWidth / 2 + Constants.gatesWidth / 2, 0));
        const rightTop = Utils.convertToCanvasXY(new Position(Constants.gameWidth, 0));
        const leftBottom = Utils.convertToCanvasXY(new Position(0, Constants.gameHeight));
        const leftBottomGate = Utils.convertToCanvasXY(new Position(Constants.gameWidth / 2 - Constants.gatesWidth / 2, Constants.gameHeight));
        const rightBottomGate = Utils.convertToCanvasXY(new Position(Constants.gameWidth / 2 + Constants.gatesWidth / 2, Constants.gameHeight));
        const rightBottom = Utils.convertToCanvasXY(new Position(Constants.gameWidth, Constants.gameHeight));

        this.context.moveTo(leftTop.x, leftTop.y);
        this.context.lineTo(leftTopGate.x, leftTopGate.y);
        this.context.moveTo(rightTopGate.x, rightTopGate.y);
        this.context.lineTo(rightTop.x, rightTop.y);
        this.context.lineTo(rightBottom.x, rightBottom.y);
        this.context.lineTo(rightBottomGate.x, rightBottomGate.y);
        this.context.moveTo(leftBottomGate.x, leftBottomGate.y);
        this.context.lineTo(leftBottom.x, leftBottom.y);
        this.context.lineTo(leftTop.x, leftTop.y);

        this.context.stroke();
    }

    drawPlayer(position) {
        const canvasPosition = Utils.convertToCanvasXY(position);

        this.context.beginPath();
        this.context.arc(canvasPosition.x, canvasPosition.y, Constants.entityRadius, 0, 2 * Math.PI);
        this.context.stroke();
    }

    drawPuck(position) {
        const canvasPosition = Utils.convertToCanvasXY(position);

        this.context.beginPath();
        this.context.arc(canvasPosition.x, canvasPosition.y, Constants.entityRadius, 0, 2 * Math.PI);
        this.context.fill();
        this.context.stroke();
    }

    drawGoals(countA, countB) {
        this.drawText(countA, new Position(Constants.canvasWidth / 2, Constants.canvasHeight), Constants.font, "center");
        this.drawText(countB, new Position(Constants.canvasWidth / 2, Constants.fontSize), Constants.font, "center");
    }

    drawText(text, position, font, align) {
        this.context.font = font || Constants.font;
        this.context.textAlign = align || "start";
        this.context.fillText(text, position.x, position.y);
    }
}