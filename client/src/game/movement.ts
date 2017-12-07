import { GameState } from 'game/game';

export const MovementUpdate = (state: GameState): void => {
	typiesMovement(state);
};

const typiesMovement = (state: GameState): void => {
	state.entities.typies.forEach(typie => {
		// Reach finish line.
		if (typie.targetX > state.canvas.fgCanvas.width - 100) {
			typie.velocityX = 0;
			typie.targetX = state.canvas.fgCanvas.width - 100;
		} else {
			typie.currentState.movement(state);
		}
	});
};
