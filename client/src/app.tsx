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
			playerState: 'waiting',
			counterVal: 3,
			gameRoom: null,
			player: null
		};
	}

	public render() {
		return (
			<div>
				<h1 className="countdown" />
				<h1>Typie Bird</h1>
				{this.renderButtons()}
				{this.checkPlayersState() && this.state.counterVal !== 0 ? (
					<h1 className="counter">{this.state.counterVal}</h1>
				) : null}
				{this.checkPlayersState() && this.state.counterVal === 0 ? (
					<Typing
						playerState={this.state.playerState}
						getCurrentHost={() => this.getCurrentHost()}
					/>
				) : null}
				<canvas id="bg-canvas" />
				<canvas id="fg-canvas" />
			</div>
		);
	}

	public componentWillMount() {
		// this.fetchGameRoom();
		this.fetchPlayer();
	}

	public componentDidUpdate() {
		const counter = this.state.counter;
		// Start the game when the counterVal is 0,
		// and stop counter.
		if (this.state.counterVal === 0) {
			clearInterval(counter);
		}
	}

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
		return websocket;
	};

	// When a new player first joins the game room,
	// fetch the most updated game room.
	// This is a hacky way to get game room.
	// Sending request to this url will cause websocket to broadcast game room.
	// We are not really getting any response data back.
	private fetchGameRoom = (): void => {
		const url = `http://${this.getCurrentHost()}/gameroom`;
		axios
			.get(url)
			.then(res => {
				// Fetch game state and store it locally.
				const websocket = this.establishWebsocket();
				websocket.addEventListener('message', event => {
					const gameRoom = JSON.parse(event.data);
					console.log(gameRoom);
					this.setState({
						gameRoom: gameRoom
					});
				});
				this.setState({
					gameRoom: res.data
				});
				console.log(res.data);
				Game.Init(websocket, res.data);
			})
			.catch(error => {
				console.log(error.response.data);
			});
	};

	private fetchPlayer = (): void => {
		const typieID = localStorage.getItem('TypieID');
		if (!typieID) {
			return;
		}
		const url = `http://${this.getCurrentHost()}/typie/me?auth=${typieID}`;
		axios
			.get(url)
			.then(res => {
				this.setState({
					player: res.data
				});
				console.log(res);
			})
			.catch(error => {
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

	private renderButtons = (): JSX.Element => {
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
		axios
			.patch(url)
			.then(res => {
				console.log(res.data);
			})
			.catch(error => {
				console.log(error.response.data);
			});

		// If all players are ready, start the game.
		if (this.checkPlayersState()) {
			let counterVal = this.state.counterVal;
			const counter = setInterval(() => {
				if (counterVal !== 0) {
					counterVal--;
					this.setState({
						counterVal: counterVal
					});
				}
			}, 1000);
			this.setState({
				counter: counter
			});
		}
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
		axios
			.patch(url)
			.then(res => {
				console.log(res.data);
			})
			.catch(error => {
				console.log(error.response.data);
			});
	};

	private checkPlayersState = (): boolean => {
		let result = true;
		const gameRoom = this.state.gameRoom;
		if (!gameRoom || !gameRoom.Players) {
			return result;
		}
		gameRoom.Players.forEach(player => {
			if (!player.isReady) {
				result = false;
				return;
			}
		});
		return result;
	};
}

ReactDOM.render(<App />, document.getElementById('app'));
