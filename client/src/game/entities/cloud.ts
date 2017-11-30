import Sprite, { Entity } from 'game/entities';

class Cloud implements Entity {
	public type: string;
	public sprite: Sprite;
	public targetX: number;
	public targetY: number;
	public targetWidth: number;
	public targetHeight: number;

	constructor(
		spritesheet: HTMLImageElement,
		targetX: number,
		targetY: number,
		targetWidth: number = 160,
		targetHeight: number = 80
	) {
		this.type = 'Cloud';
		this.sprite = new Sprite(spritesheet, 0, 64, 160, 80);
		this.targetX = targetX;
		this.targetY = targetY;
		this.targetWidth = targetWidth;
		this.targetHeight = targetHeight;
	}
}

export default Cloud;
