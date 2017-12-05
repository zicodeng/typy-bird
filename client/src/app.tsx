import * as React from 'react';
import * as ReactDOM from 'react-dom';

import * as Game from 'game/game';
import Typing from 'components/typing';

import 'sass/app';

class App extends React.Component<any, any> {
	constructor(props, context) {
		super(props, context);

		this.state = {
			playerState: 'waiting'
		};
	}

	public render() {
		return (
			<div>
				<h1 className="countdown" />
				<h1>Typie Bird</h1>
				{this.renderButtons()}
				<Typing playerState={this.state.playerState} />
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
		// If all players are ready, start the game.
		if (this.checkPlayersState()) {
			console.log('start');
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
	};

	private handleClickCancel = (): void => {
		this.setState({
			playerState: 'waiting'
		});
	};

	private checkPlayersState = (): boolean => {
		return true;
	};
}

ReactDOM.render(<App />, document.getElementById('app'));
