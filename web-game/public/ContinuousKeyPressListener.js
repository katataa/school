class ContinuousKeyPressListener {
    /**
     * Creates a listener that calls the callback repeatedly while the key is held down.
     * @param {string} keyCode - The key code (e.g., "ArrowUp").
     * @param {function} callback - The function to call repeatedly.
     * @param {number} interval - The delay (in ms) between repeated calls.
     */
    constructor(keyCode, callback, interval = 100) {
      this.keyCode = keyCode;
      this.callback = callback;
      this.interval = interval;
      this.intervalId = null;
  
      this.keydownFunction = (event) => {
        if (event.code === this.keyCode && this.intervalId === null) {
          this.callback();
          this.intervalId = setInterval(() => {
            this.callback();
          }, this.interval);
        }
      };
  
      this.keyupFunction = (event) => {
        if (event.code === this.keyCode) {
          if (this.intervalId !== null) {
            clearInterval(this.intervalId);
            this.intervalId = null;
          }
        }
      };
  
      document.addEventListener("keydown", this.keydownFunction);
      document.addEventListener("keyup", this.keyupFunction);
    }
  
    unbind() {
      document.removeEventListener("keydown", this.keydownFunction);
      document.removeEventListener("keyup", this.keyupFunction);
      if (this.intervalId !== null) {
        clearInterval(this.intervalId);
        this.intervalId = null;
      }
    }
  }
  
  export default ContinuousKeyPressListener;
  