export interface Sprite {
	spritesheet: HTMLImageElement;
	srcX: number;
	srcY: number;
	srcWidth: number;
	srcHeight: number;
}

export interface Entity {
	type: string;
	sprite: Sprite;
	targetX: number;
	targetY: number;
	targetWidth: number;
	targetHeight: number;
}
