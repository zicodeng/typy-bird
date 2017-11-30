import { GameState } from 'game/game';

export const AnimUpdate = (state: GameState): void => {
	typiesAnim(state);
};

const typiesAnim = (state: GameState): void => {
	state.entities.typies.forEach(typie => {
		typie.currentState.anim(state);
	});
};
