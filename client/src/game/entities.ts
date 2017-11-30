import { GameState } from './game';
import Typie from './entities/typie';
import Cloud from './entities/cloud';

export interface Entity {
	type: string;
	sprite: Sprite;
	targetX: number;
	targetY: number;
	targetWidth: number;
	targetHeight: number;
}

export default class Sprite {
	constructor(
		public spritesheet: HTMLImageElement,
		public srcX: number,
		public srcY: number,
		public srcWidth: number,
		public srcHeight: number
	) {}
}

export const EntitiesInit = (state: GameState): void => {
	const canvasWidth = state.canvas.fgCanvas.width;
	const canvasHeight = state.canvas.fgCanvas.height;

	const typieCount = 4;
	state.entities.typies = [];
	for (let i = 0; i < typieCount; i++) {
		let posY = canvasHeight / typieCount * i + canvasHeight / typieCount / 2;
		state.entities.typies.push(new Typie(state.spritesheet, 100, posY));
	}

	// Define cloud positions.
	var cloudPos = [
		[Math.floor(canvasWidth / 10 * 0.5), Math.floor(canvasHeight / 10 * 6.5)],
		[Math.floor(canvasWidth / 10 * 2), Math.floor(canvasHeight / 10 * 2)],
		[Math.floor(canvasWidth / 10 * 5), Math.floor(canvasHeight / 10 * 5)],
		[Math.floor(canvasWidth / 10 * 8), Math.floor(canvasHeight / 10 * 4)],
		[Math.floor(canvasWidth / 10 * 4), Math.floor(canvasHeight / 10 * 9)],
		[Math.floor(canvasWidth / 10 * 9), Math.floor(canvasHeight / 10 * 0.5)]
	];

	state.entities.clouds = [];
	cloudPos.forEach(function(pos) {
		state.entities.clouds.push(new Cloud(state.spritesheet, pos[0], pos[1]));
	});
};
