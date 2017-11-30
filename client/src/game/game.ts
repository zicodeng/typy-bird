import { EntitiesInit } from './entities';
import { RenderInit, RenderUpdate } from 'game/render';

const Spritesheet = require('spritesheet/game.png');

interface GameCanvas {
	bgCanvas: HTMLCanvasElement;
	fgCanvas: HTMLCanvasElement;
	bgCtx: CanvasRenderingContext2D;
	fgCtx: CanvasRenderingContext2D;
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

	const bgCtx = bgCanvas.getContext('2d');
	const fgCtx = bgCanvas.getContext('2d');

	if (!bgCtx || !fgCtx) {
		throw new Error('Canvas context is null');
	}

	const canvasWidth = '1200px';
	const canvasHeight = '600px';

	// Set canvas size.
	bgCanvas.style.width = canvasWidth;
	bgCanvas.style.height = canvasHeight;

	fgCanvas.style.width = canvasWidth;
	fgCanvas.style.height = canvasHeight;

	const canvas: GameCanvas = {
		bgCanvas: bgCanvas,
		fgCanvas: fgCanvas,
		bgCtx: bgCtx,
		fgCtx: fgCtx
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
	EntitiesInit(state);
	RenderInit(state);

	run(state);
};

// Running the game.
const run = (state: GameState): void => {
	const loop = () => {
		render(state);

		state.animFrame++;
		window.requestAnimationFrame(loop);
	};
	loop();
};

const render = (state: GameState): void => {
	RenderUpdate(state);
};
