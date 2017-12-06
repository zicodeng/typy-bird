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
			counterVal: 3
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
					<Typing playerState={this.state.playerState} />
				) : null}
				<canvas id="bg-canvas" />
				<canvas id="fg-canvas" />
			</div>
		);
	}

	public componentDidMount() {
		// Fetch game state and store it locally.
		Game.Init();
	}

	public componentDidUpdate() {
		const counter = this.state.counter;
		// Start the game when the counterVal is 0,
		// and stop counter.
		if (this.state.counterVal === 0) {
			clearInterval(counter);
		}
	}

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
	};

	private checkPlayersState = (): boolean => {
		return true;
	};
}

ReactDOM.render(<App />, document.getElementById('app'));
