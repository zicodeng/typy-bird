import { Entity } from './entities';
import { GameState } from './game';

interface Text {
	content: string;
	color: string;
	font: string;
	xPos: number;
	yPos: number;
}

// Initialize game environment entities on background canvas.
// These entities will not interact with each other,
// but they might interact with foreground entities.
export const RenderInit = (state: GameState): void => {
	// Draw clouds.
	state.entities.clouds.forEach(cloud => {
		drawEntity(cloud, state.canvas.bgCtx);
	});
};

// Render game entities on foreground canvas.
// These entities will frequently interact with each other,
// and update game state.
export const RenderUpdate = (state: GameState): void => {
	state.canvas.fgCtx.clearRect(0, 0, state.canvas.fgCanvas.width, state.canvas.fgCanvas.height);

	// Draw Typies.
	state.entities.typies.forEach(typie => {
		drawEntity(typie, state.canvas.fgCtx);
		const typieName: Text = {
			content: typie.userName,
			color: '#000',
			font: '15px sans-serif',
			xPos: typie.targetX,
			yPos: typie.targetY + typie.targetHeight + 20
		};
		drawText(typieName, state.canvas.fgCtx);
	});

	// Draw hearts.
	state.entities.hearts.forEach(heart => {
		drawEntity(heart, state.canvas.fgCtx);
	});
};

const drawEntity = (entity: Entity, ctx: CanvasRenderingContext2D): void => {
	ctx.drawImage(
		entity.sprite.spritesheet,
		entity.sprite.srcX,
		entity.sprite.srcY,
		entity.sprite.srcWidth,
		entity.sprite.srcHeight,
		entity.targetX,
		entity.targetY,
		entity.targetWidth,
		entity.targetHeight
	);
};

const drawText = (text: Text, ctx: CanvasRenderingContext2D): void => {
	ctx.fillStyle = text.color;
	ctx.font = text.font;
	ctx.fillText(text.content, text.xPos, text.yPos);
};
