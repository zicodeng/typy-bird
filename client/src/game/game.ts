import { PhysicsUpdate } from './physics';
import { MovementUpdate } from './movement';
import { EntitiesInit } from './entities';
import { RenderInit, RenderUpdate } from 'game/render';
import { AnimUpdate } from './anim';

import Typie from './entities/typie';
import Heart from './entities/heart';

import axios from 'axios';

const Spritesheet = require('spritesheet/game.png');

interface GameCanvas {
	bgCanvas: HTMLCanvasElement;
	fgCanvas: HTMLCanvasElement;
	bgCtx: CanvasRenderingContext2D;
	fgCtx: CanvasRenderingContext2D;
}

export interface GameState {
	startTime: number;
	spritesheet: HTMLImageElement;
	canvas: GameCanvas;
	animFrame: number;
	entities: any;
}

interface GameRoom {
	available: boolean;
	players: any;
}

// Initialize the game.
export const Init = (websocket: WebSocket, initGameRoom: GameRoom): void => {
	const bgCanvas = <HTMLCanvasElement>document.getElementById('bg-canvas');
	const fgCanvas = <HTMLCanvasElement>document.getElementById('fg-canvas');

	if (!bgCanvas || !fgCanvas) {
		throw new Error('Canvas is null');
	}

	const bgCtx = bgCanvas.getContext('2d');
	const fgCtx = fgCanvas.getContext('2d');

	if (!bgCtx || !fgCtx) {
		throw new Error('Canvas context is null');
	}

	// Don't use CSS style to set Canvas size,
	// because it will cause scaling issues to entities.
	const canvasWidth = window.innerWidth;
	const canvasHeight = window.innerHeight;

	// Set canvas size.
	bgCanvas.width = canvasWidth;
	bgCanvas.height = canvasHeight;

	fgCanvas.width = canvasWidth;
	fgCanvas.height = canvasHeight;

	const canvas: GameCanvas = {
		bgCanvas: bgCanvas,
		fgCanvas: fgCanvas,
		bgCtx: bgCtx,
		fgCtx: fgCtx
	};

	var spritesheet = new Image();
	// Spritesheet is converted to data URL by url-loader,
	// so it is immediately available without another round trip to the server.
	spritesheet.src = Spritesheet;
	// Initialize game state.
	const state: GameState = {
		startTime: 0,
		canvas: canvas,
		spritesheet: spritesheet,
		animFrame: 0,
		entities: {}
	};

	state.entities.typies = [];
	state.entities.hearts = [];

	// Load players in current game room first.
	initGameRoom.players.forEach((player, i) => {
		renderTypie(state, player.id, i);
	});

	// Update game state based on the server's response.
	websocket.addEventListener('message', event => {
		// Change state that will get passed to update and render functions.
		const data = JSON.parse(event.data);
		const gameRoom = data.payload;
		console.log(gameRoom);
		switch (data.type) {
			case 'Ready':
				gameRoom.players.forEach(player => {
					state.entities.typies.forEach(typie => {
						if (player.id === typie.id) {
							if (player.isReady) {
								typie.isReady = true;
								typie.currentState = typie.states.moving;
							} else {
								typie.isReady = false;
								typie.currentState = typie.states.standing;
							}
						}
					});
				});
				break;

			case 'NewTypie':
				const playerID = data.players[gameRoom.players.length].ID;
				// If this data we received is related to creating a new Typie.
				renderTypie(state, playerID, gameRoom.players.length);
				break;

			case 'Position':
				let isGameEnded = true;
				gameRoom.players.forEach(player => {
					state.entities.typies.forEach(typie => {
						if (player.id === typie.id) {
							typie.targetX = calcPos(player.position);
							if (player.position === 20) {
								reachFinishLine(state, player.id);
							}
						}
						if (player.position !== 20) {
							isGameEnded = false;
						}
					});
				});
				if (isGameEnded) {
					// Clear all typies in game state.
					state.entities.typies = [];
					endGame();
				}
				break;

			case 'GameStart':
				const startTime = data.startTime;
				state.startTime = Date.parse(startTime);
				break;

			default:
				break;
		}
	});

	console.log(state.entities.typies);

	EntitiesInit(state);

	run(state);
};

// Running the game.
const run = (state: GameState): void => {
	const loop = () => {
		update(state);
		render(state);

		state.animFrame++;
		window.requestAnimationFrame(loop);
	};
	loop();
};

const update = (state: GameState): void => {
	AnimUpdate(state);
	MovementUpdate(state);
	PhysicsUpdate(state);
};

const render = (state: GameState): void => {
	RenderUpdate(state);
	RenderInit(state);
};

const leftMargin = 50;
const rightMargin = 100;

const renderTypie = (state: GameState, playerID: number, i: number): void => {
	const maxPlayer = 4;
	const canvasWidth = window.innerWidth;
	const canvasHeight = window.innerHeight;

	state.entities.typies.push(
		new Typie(
			state.spritesheet,
			playerID,
			leftMargin,
			canvasHeight / maxPlayer * i + canvasHeight / maxPlayer / 2
		)
	);
	state.entities.hearts.push(
		new Heart(
			state.spritesheet,
			canvasWidth - rightMargin,
			canvasHeight / maxPlayer * i + canvasHeight / maxPlayer / 2
		)
	);
};

const calcPos = (pos: number) => {
	const canvasWidth = window.innerWidth - leftMargin - rightMargin;
	return canvasWidth / 20 * pos;
};

const reachFinishLine = (state: GameState, playerID: number): void => {
	const url = `http://${getCurrentHost()}/typie/me?auth=${playerID}`;
	// Send this player's record to server.
	const record = {
		record: (Date.now() - state.startTime) / 1000
	};
	axios.patch(url, record).catch(error => {
		console.log(error.response.data);
	});
};

const endGame = () => {
	const url = `http://${getCurrentHost()}/gameroom`;
	axios.post(url).catch(err => {
		console.log(err);
	});
	window.location.replace('index.html');
};

const getCurrentHost = (): string => {
	let host: string;
	if (window.location.hostname === 'typy-bird.zicodeng.me') {
		host = 'typy-bird-api.zicodeng.me';
	} else {
		host = 'localhost:3000';
	}
	return host;
};
