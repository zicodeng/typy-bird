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
				<button className="btn btn--ready" onClick={e => this.handleClickReady()}>
					READY
				</button>
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

	private handleClickReady = (): void => {
		console.log('ready');
		this.setState({
			playerState: 'ready'
		});
	};

	private checkPlayersState = (): boolean => {
		return true;
	};
}

ReactDOM.render(<App />, document.getElementById('app'));
