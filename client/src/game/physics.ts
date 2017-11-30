import Typie from './entities/typie';
import { GameState } from 'game/game';
import { Entity } from 'game/entities';

export const PhysicsUpdate = (state: GameState): void => {
	state.entities.typies.forEach(typie => {
		gravity(typie);
	});
};

const gravity = (entity: Entity): void => {
	if (!entity.velocityY) {
		return;
	}
	entity.velocityY += 1.5;
	entity.targetY += entity.velocityY;
	if (entity.initTargetY && entity.targetY > entity.initTargetY) {
		entity.targetY = entity.initTargetY;
		entity.velocityY = 0;
	}
};
