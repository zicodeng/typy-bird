import { GameState } from './game';
import Typie from './entities/typie';

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
};
