const Spritesheet = require('spritesheet/game.png');

interface GameCanvas {
	bgCanvas: HTMLCanvasElement | null;
	fgCanvas: HTMLCanvasElement | null;
	bgCtx: CanvasRenderingContext2D | null;
	fgCtx: CanvasRenderingContext2D | null;
}

export interface GameState {
	spritesheet: HTMLImageElement;
	canvas: GameCanvas;
	animFrame: number;
	entities: object;
}

// Initialize the game.
export const Init = (): void => {
	const bgCanvas = <HTMLCanvasElement>document.getElementById('bg-canvas');
	const fgCanvas = <HTMLCanvasElement>document.getElementById('fg-canvas');

	if (!bgCanvas || !fgCanvas) {
		throw new Error('Canvas is null');
	}

	// Set canvas size.
	bgCanvas.style.width = '1500px';
	bgCanvas.style.height = '650px';

	fgCanvas.style.width = '1500px';
	fgCanvas.style.height = '650px';

	const canvas: GameCanvas = {
		bgCanvas: bgCanvas,
		fgCanvas: fgCanvas,
		bgCtx: bgCanvas.getContext('2d'),
		fgCtx: fgCanvas.getContext('2d')
	};

	var spritesheet = new Image();
	// Spritesheet is converted to data URL by url-loader,
	// so it is immediately available without another round trip to the server.
	spritesheet.src = Spritesheet;
	// Initialize game state.
	const state: GameState = {
		canvas: canvas,
		spritesheet: spritesheet,
		animFrame: 0,
		entities: {}
	};
	console.log(state);
	run(state);
};

// Running the game.
const run = (state: GameState): void => {
	const loop = () => {
		state.animFrame++;

		window.requestAnimationFrame(loop);
	};
	loop();
};
