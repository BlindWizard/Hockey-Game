import React from "react";
import { ErrorTypes } from "../errorTypes";

const ErrorAlert = ({ error, resetErrorBoundary }) => {
    if (!error.cause) {
        throw error;
    }

    return (
        <div className="app__error">
            <h1 className="app__logo">Hockey!</h1>

            <div className="app__error-message">
                {error.message}
            </div>

            {error.cause !== ErrorTypes.badBrowser &&
                <button className="app__button" onClick={resetErrorBoundary}>Return to the main menu</button>
            }
        </div>
    )
};

export default ErrorAlert;