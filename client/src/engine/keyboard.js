export class KeyBoard {
    onDebugHandler(handler) {
        this.debugHangler = handler;
    }

    setHandlers() {
        document.onkeyup = ((e) => {
            if (e.code == 'KeyO') {
                this.debugHangler && this.debugHangler()
            }
        })
    }
}