import React, { useEffect, useRef, useState } from "react";
import { Game } from "../engine/game";
import { Graphics } from "../engine/graphics";
import { useAsyncError } from "../errorHandler";
import { ErrorTypes } from "../errorTypes";

const GameField = ({ network, keyboard }) => {
    const canvasRef = useRef(null);
    const throwError = useAsyncError();
    
    let game;

    useEffect(() => {
        const graphics = new Graphics(canvasRef.current);
        graphics.drawField();

        game = new Game(graphics, network);
        game.setDefaultPositions();
        game.run();

        network.setOnWorld((message) => {
            game.receiveWorld(message.opponentPosition, message.puckPosition, message.countA, message.countB);
        });
        network.setOnExitGame((message) => {
            throwError(new Error(message.reason, {cause: ErrorTypes.gameExit}))
        })

        keyboard.onDebugHandler(() => {
            game.toggleDebug();
        })

        return () => game.exit('Player leaved the room');
    }, []);

    const leaveGame = () => {
        game && game.exit('Player leaved the room');
    };

    return (
        <div className="app__game-page">
            <button className="app__button" onClick={() => leaveGame()}>Leave game</button>
            <canvas className="app__game" ref={canvasRef}></canvas>
        </div>
    )
}

export default GameField;