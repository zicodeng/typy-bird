import { GameState } from './game';

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

export const EntitiesInit = (state: GameState): void => {};
