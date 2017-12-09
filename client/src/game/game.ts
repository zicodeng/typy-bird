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
	startTime: Date;
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
		startTime: new Date(),
		canvas: canvas,
		spritesheet: spritesheet,
		animFrame: 0,
		entities: {}
	};

	state.entities.typies = [];
	state.entities.hearts = [];

	// Load players in current game room first.
	initGameRoom.players.forEach((player, i) => {
		renderTypie(state, player.id, player.userName, i);
	});

	// Update game state based on the server's response.
	websocket.addEventListener('message', event => {
		// Change state that will get passed to update and render functions.
		const data = JSON.parse(event.data);
		const gameRoom = data.payload;
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
				const playerID = gameRoom.players[gameRoom.players.length - 1].id;
				const userName = gameRoom.players[gameRoom.players.length - 1].userName;
				// If this data we received is related to creating a new Typie.
				renderTypie(state, playerID, userName, gameRoom.players.length - 1);
				break;

			case 'Position':
				gameRoom.players.forEach(player => {
					state.entities.typies.forEach(typie => {
						if (player.id === typie.id) {
							typie.targetX = calcPos(player.position);
							if (player.position === 20 && !typie.isDone) {
								reachFinishLine(state, gameRoom, player.id);
								typie.isDone = true;
								typie.currentState = typie.states.standing;
							}
						}
					});
				});
				break;

			case 'GameStart':
				state.startTime = new Date(data.startTime);
				break;

			default:
				break;
		}
	});

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

const renderTypie = (state: GameState, playerID: number, userName: string, i: number): void => {
	const maxPlayer = 4;
	const canvasWidth = window.innerWidth;
	const canvasHeight = window.innerHeight;

	state.entities.typies.push(
		new Typie(
			state.spritesheet,
			playerID,
			userName,
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

const reachFinishLine = (state: GameState, gameRoom: GameRoom, playerID: number): void => {
	const url = `http://${getCurrentHost()}/typie/me?auth=${playerID}`;
	// Send this player's record to server.
	const record = {
		record: (new Date().getTime() - state.startTime.getTime()) / 1000
	};
	axios
		.patch(url, record)
		.then(res => {
			let isGameEnded = true;
			gameRoom.players.forEach(player => {
				console.log(player.position);
				if (player.position !== 20) {
					isGameEnded = false;
				}
			});
			if (isGameEnded) {
				endGame(state);
			}
		})
		.catch(error => {
			console.log(error.response.data);
		});
};

const endGame = (state: GameState) => {
	const url = `http://${getCurrentHost()}/end`;
	axios.post(url).catch(err => {
		console.log(err);
	});
	setTimeout(() => {
		// Clear all typies in game state.
		state.entities.typies = [];
		window.location.replace('index.html');
	}, 3000);
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
