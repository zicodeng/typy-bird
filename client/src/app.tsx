import axios from 'axios';
import * as React from 'react';
import * as ReactDOM from 'react-dom';

import * as Game from 'game/game';
import Typing from 'components/typing';

import 'sass/app';

class App extends React.Component<any, any> {
	constructor(props, context) {
		super(props, context);

		this.state = {
			isGameInProgress: false,
			playerState: 'waiting',
			counter: null,
			counterVal: 3,
			gameRoom: null,
			player: null
		};
	}

	public render() {
		return (
			<div>
				<h1 className="countdown" />
				<h4 className="game-title">Typy Bird</h4>
				{this.renderButtons()}
				{this.checkPlayersState() && this.state.counterVal !== 0 ? (
					<h1 className="counter">{this.state.counterVal}</h1>
				) : null}
				{this.checkPlayersState() && this.state.counterVal === 0 ? (
					<Typing
						playerState={this.state.playerState}
						getCurrentHost={() => this.getCurrentHost()}
						playerPos={this.state.player.position}
					/>
				) : null}
				<canvas id="bg-canvas" />
				<canvas id="fg-canvas" />
			</div>
		);
	}

	public componentWillMount() {
		this.fetchGame();
	}

	public componentDidUpdate() {
		const counter = this.state.counter;
		// Start the game when the counterVal is 0,
		// and stop counter.
		if (this.state.gameRoom.available && this.state.counterVal === 0) {
			clearInterval(counter);
			this.startGame();
		}
	}

	// Inform the server that our game has started.
	private startGame = (): void => {
		const url = `http://${this.getCurrentHost()}/start`;
		axios.post(url).catch(err => {
			console.log(err);
		});
	};

	private establishWebsocket = (): WebSocket => {
		const websocket = new WebSocket(`ws://${this.getCurrentHost()}/ws`);
		websocket.addEventListener('error', function(err) {
			console.log(err);
			// If the connection is lost,
			// the player will be forced to quit the game.
			localStorage.removeItem('TypieID');
			window.location.replace('index.html');
		});
		websocket.addEventListener('open', function() {
			console.log('Websocket connection established');
		});
		websocket.addEventListener('close', function() {
			console.log('Websocket connection closed');
		});
		websocket.addEventListener('message', event => {
			const data = JSON.parse(event.data);
			const gameRoom = data.payload;
			switch (data.type) {
				case 'Ready':
					this.setState({
						gameRoom: gameRoom
					});
					// If all players are ready, start the game.
					if (this.checkPlayersState()) {
						let counterVal = 3;
						this.setState({
							counter: setInterval(() => {
								if (counterVal !== 0) {
									counterVal--;
									this.setState({
										counterVal: counterVal
									});
								}
							}, 1000)
						});
					}
					break;

				case 'Position':
					gameRoom.players.forEach(player => {
						if (this.state.player.id === player.id) {
							this.setState({
								player: player
							});
						}
					});
					break;

				case 'GameStart':
					this.setState({
						gameRoom: gameRoom
					});
					break;

				default:
					break;
			}
		});
		return websocket;
	};

	// When a new player first joins the game room,
	// fetch the most updated game room.
	private fetchGame = (): void => {
		const url = `http://${this.getCurrentHost()}/gameroom`;
		axios
			.get(url)
			.then(res => {
				const gameRoom = res.data;
				const websocket = this.establishWebsocket();
				Game.Init(websocket, gameRoom);
				// Store the current player in memory.
				const typieID = localStorage.getItem('TypieID');
				gameRoom.players.forEach(player => {
					if (player.id === typieID) {
						this.setState({
							gameRoom: gameRoom,
							player: player
						});
						return;
					}
				});
				if (this.checkPlayersState()) {
					this.setState({
						playerState: 'ready',
						counterVal: 0
					});
				}
			})
			.catch(error => {
				// If no such player
				// redirect to waiting room page.
				localStorage.removeItem('TypieID');
				window.location.replace('index.html');
				console.log(error.response.data);
			});
	};

	private getCurrentHost = (): string => {
		let host: string;
		if (window.location.hostname === 'typy-bird.zicodeng.me') {
			host = 'typy-bird-api.zicodeng.me';
		} else {
			host = 'localhost:3000';
		}
		return host;
	};

	private renderButtons = (): JSX.Element | null => {
		// Prevent players cancelling the game if it is already started;
		if (this.state.counterVal === 0) {
			return null;
		}
		const playerState = this.state.playerState;
		if (playerState === 'waiting') {
			return (
				<button className="btn btn--ready" onClick={e => this.handleClickReady()}>
					READY
				</button>
			);
		}
		return (
			<button className="btn btn--cancel" onClick={e => this.handleClickCancel()}>
				CANCEL
			</button>
		);
	};

	private handleClickReady = (): void => {
		this.setState({
			playerState: 'ready'
		});

		const typieID = localStorage.getItem('TypieID');
		if (!typieID) {
			return;
		}
		// Update server.
		const url = `http://${this.getCurrentHost()}/ready?auth=${typieID}`;
		axios.patch(url).catch(error => {
			console.log(error.response.data);
		});
	};

	private handleClickCancel = (): void => {
		const counter = this.state.counter;
		clearInterval(counter);
		this.setState({
			playerState: 'waiting',
			counterVal: 3
		});
		const typieID = localStorage.getItem('TypieID');
		if (!typieID) {
			return;
		}
		// Update server.
		const url = `http://${this.getCurrentHost()}/ready?auth=${typieID}`;
		axios.patch(url).catch(error => {
			console.log(error.response.data);
		});
	};

	private checkPlayersState = (): boolean => {
		let result = true;
		const gameRoom = this.state.gameRoom;
		if (!gameRoom || !gameRoom.players) {
			return false;
		}
		gameRoom.players.forEach(player => {
			if (!player.isReady) {
				result = false;
				return;
			}
		});
		return result;
	};
}

ReactDOM.render(<App />, document.getElementById('app'));
