import { GameState } from 'game/game';

export const MovementUpdate = (state: GameState): void => {
	typiesMovement(state);
};

const typiesMovement = (state: GameState): void => {
	state.entities.typies.forEach(typie => {
		typie.currentState.movement(state);
	});
};
