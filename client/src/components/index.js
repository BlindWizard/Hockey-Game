import React, { useMemo, useState } from "react";
import Menu from "./menu";
import GameField from "./gameField";
import { Pages } from "./pages";
import { Network } from "../engine/network";
import { useAsyncError } from "../errorHandler";
import { ClientMessages } from "../protocol/messages";
import { KeyBoard } from "../engine/keyboard";

const Index = ({ currentPage, setPage, playerId }) => {
    const [isConnectionEstablished, setConnectionEstablished] = useState(false);
    const throwError = useAsyncError();
  
    const network = useMemo(() => {
        const network = new Network(throwError);
        network.openConnection();

        network.setOnOpen(() => {
            network.send(ClientMessages.Hello(playerId.getId()));
        });

        network.setOnHello((message) => {
            playerId.storeId(message.playerId);
            setConnectionEstablished(true);
        });

        network.setOnGame(() => {
            setPage(Pages.room);
        })

        return network;
    }, []);

    const keyboard = useMemo(() => {
        const keyboard = new KeyBoard();
        keyboard.setHandlers();

        return keyboard;
    }, []);

    return (
        <div className="app__page">
            {currentPage === Pages.menu && <Menu network={network} isConnectionEstablished={isConnectionEstablished} />}
            {currentPage === Pages.room && <GameField setPage={setPage} network={network} keyboard={keyboard} />}
        </div>
    );
};

export default Index;