import { GameState } from 'game/game';

export const AnimUpdate = (state: GameState): void => {
	typiesAnim(state);
	heartsAnim(state);
};

const typiesAnim = (state: GameState): void => {
	state.entities.typies.forEach(typie => {
		// Animate the Typie's state if it is "ready".
		typie.currentState.anim(state);
	});
};

const heartsAnim = (state: GameState): void => {
	state.entities.hearts.forEach(heart => {
		heart.currentState.anim(state);
	});
};
