import { Constants } from "./constants";
import { Position } from "./position";

const Utils = {
    convertToCanvasXY: (gamePosition) => {
        return new Position(
            gamePosition.x + Constants.lineWidth / 2,
            gamePosition.y + Constants.gatesHeight + Constants.lineWidth / 2
        );
    },

    convertToGameXY: (canvasPosition) => {
        return new Position(
            canvasPosition.x - Constants.lineWidth / 2,
            canvasPosition.y - Constants.gatesHeight - Constants.lineWidth / 2
        );
    },

    clamp: (num, min, max) => Math.min(Math.max(num, min), max),

    lerp: (start, end, amt) => (1 - amt) * start + amt * end,
}

export { Utils };