import React, { useState, useRef, useEffect } from "react";
import { Graphics } from "../engine/graphics";
import { Position } from "../engine/position";
import { Constants } from "../engine/constants";
import { ClientMessages } from "../protocol/messages";
import { useAsyncError } from "../errorHandler";

const Menu = ({ network, isConnectionEstablished }) => {
    const [playersOnline, setPlayersOnline] = useState(null);
    const [isFindingGame, setIsFindingGame] = useState(false);
    const [latency, setLatency] = useState(0);
    const throwError = useAsyncError();

    const canvasRef = useRef(null);

    if (isConnectionEstablished) {
        network.runPing();
    }

    const findGame = () => {
        try {
            if (!isFindingGame) {
                network.send(ClientMessages.Queue());
                setIsFindingGame(true);
            } else {
                network.send(ClientMessages.UnQueue());
                setIsFindingGame(false);
            }        
        } catch (e) {
            throwError(e);
        }
    }

    useEffect(() => {
        const graphics = new Graphics(canvasRef.current);
        graphics.drawField();
        graphics.drawPlayer(new Position(
            Constants.gameWidth / 2,
            Constants.gameHeight - 2 * Constants.entityRadius,
        ));
        graphics.drawPlayer(new Position(
            Constants.gameWidth / 2,
            2 * Constants.entityRadius,
        ));
    }, []);

    useEffect(() => {
        network.setOnPing((message) => setLatency(message.latency));
        network.setOnOnline((message) => setPlayersOnline(message.count));
    }, []);

    let btnCap = () => {
        if (isFindingGame) {
            return <span>
                <span>Finding game...</span>
                <span className="app__button-notice">Cancel</span>
            </span>;
        }

        if (isConnectionEstablished) {
            return 'Find game';
        }

        return 'Connecting...';
    }

    return (
        <div className="app__menu-page">
            <h1 className="app__logo app__logo_with-latency" data-latency={latency}>Hockey!</h1>
            <div className="app__online">
                {isConnectionEstablished ? <div className="app__online-badge"></div> : ''}
                Players online: {null !== playersOnline ? playersOnline : '...'}
                </div>
            <button className="app__button" onClick={findGame} disabled={!isConnectionEstablished}>{btnCap()}</button>
            <canvas className="app__game" ref={canvasRef}></canvas>
        </div>
    )
};

export default Menu;