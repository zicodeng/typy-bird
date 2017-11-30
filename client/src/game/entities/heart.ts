import { Entity } from 'game/entities';
import { GameState } from 'game/game';
import Sprite from 'game/entities';

class Heart implements Entity {
	public type: string;
	public sprite: Sprite;
	public targetX: number;
	public targetY: number;
	public targetWidth: number;
	public targetHeight: number;

	public spriteAnims;
	public states;
	public currentState;

	constructor(
		spritesheet: HTMLImageElement,
		targetX: number,
		targetY: number,
		targetWidth: number = 32,
		targetHeight: number = 32
	) {
		this.type = 'Heart';
		this.sprite = new Sprite(spritesheet, 0, 0, 32, 32);
		this.targetX = targetX;
		this.targetY = targetY;
		this.targetWidth = targetWidth;
		this.targetHeight = targetHeight;

		this.spriteAnims = {
			float: {
				frames: [(this.targetY = this.targetY + 5), (this.targetY = this.targetY - 5)],
				currentFrame: 0
			}
		};

		this.states = {
			floating: {
				anim: (state: GameState) => {
					if (state.animFrame % 25 === 0) {
						this.targetY = this.spriteAnims.float.frames[
							this.spriteAnims.float.currentFrame
						];
						this.spriteAnims.float.currentFrame++;
						if (
							this.spriteAnims.float.currentFrame >
							this.spriteAnims.float.frames.length - 1
						) {
							this.spriteAnims.float.currentFrame = 0;
						}
					}
				}
			}
		};

		this.currentState = this.states.floating;
	}
}

export default Heart;
