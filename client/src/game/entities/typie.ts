import Sprite, { Entity } from 'game/entities';
import { GameState } from 'game/game';

class Typie implements Entity {
	private width = 70;
	private height = 50;

	public type: string;
	public id: number;
	public sprite: Sprite;
	public targetX: number;
	public targetY: number;
	public targetWidth: number;
	public targetHeight: number;
	public isReady: boolean;

	public initTargetX;
	public initTargetY;

	public velocityX: number;
	public velocityY: number;

	public spriteAnims;
	public states;
	public currentState;

	constructor(
		spritesheet: HTMLImageElement,
		id: number,
		targetX: number,
		targetY: number,
		targetWidth: number = 70,
		targetHeight: number = 50,
		isReady: boolean = false
	) {
		this.type = 'Typie';
		this.id = id;
		this.targetX = targetX;
		this.targetY = targetY;
		this.targetWidth = targetWidth;
		this.targetHeight = targetHeight;
		this.isReady = isReady;

		this.initTargetX = targetX;
		this.initTargetY = targetY;

		this.velocityX = 0;
		this.velocityY = 0;

		this.spriteAnims = {
			move: {
				frames: [
					new Sprite(spritesheet, 166, 0, this.width, this.height),
					new Sprite(spritesheet, 236, 0, this.width, this.height),
					new Sprite(spritesheet, 96, 0, this.width, this.height)
				],
				currentFrame: 0
			},
			stand: new Sprite(spritesheet, 166, 0, this.width, this.height)
		};

		this.states = {
			moving: {
				movement: (state: GameState) => {
					this.targetX += this.velocityX;
					// Jump
					if (this.velocityY === 0) {
						this.velocityY = -14;
					}
				},
				anim: (state: GameState) => {
					if (state.animFrame % 5 === 0) {
						this.sprite = this.spriteAnims.move.frames[
							this.spriteAnims.move.currentFrame
						];
						this.spriteAnims.move.currentFrame++;

						// Reset current frame.
						if (
							this.spriteAnims.move.currentFrame >
							this.spriteAnims.move.frames.length - 1
						) {
							this.spriteAnims.move.currentFrame = 0;
						}
					}
				}
			},
			standing: {
				movement: (state: GameState) => {
					// Game state remain the same on standing.
					return;
				},
				anim: (state: GameState) => {
					this.sprite = this.spriteAnims.stand;
				}
			}
		};

		this.sprite = this.spriteAnims.stand;
		this.currentState = this.states.standing;
	}
}

export default Typie;
