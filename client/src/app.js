import "./styles/reset.css";
import "./styles/style.scss";

import React, { useState } from "react";
import { createRoot } from "react-dom/client";
import Index from "./components/index";
import { ErrorBoundary } from "react-error-boundary";
import ErrorAlert from "./components/errorAlert";
import { Pages } from "./components/pages";
import PlayerId from "./playerId";

const root = createRoot(document.getElementById('app'));
const playerId = new PlayerId();

const App = ({playerId}) => {
    let [currentPage, setPage] = useState(Pages.menu);

    return (
        <ErrorBoundary
            FallbackComponent={ErrorAlert}
            onReset={() => setPage(Pages.menu)}
        >
            <Index currentPage={currentPage} setPage={setPage} playerId={playerId} />
        </ErrorBoundary>
    )
}


root.render(
    <App playerId={playerId}/>
);
