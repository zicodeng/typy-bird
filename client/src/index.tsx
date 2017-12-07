import * as React from 'react';
import * as ReactDOM from 'react-dom';
import axios from 'axios';

import 'sass/index';

class Index extends React.Component<any, any> {
	constructor(props, context) {
		super(props, context);

		this.state = {
			disabled: false,
			available: true,
			leaderboard: []
		};
	}

	private fetchTopScores = (): void => {
		const url = `http://${this.getCurrentHost()}/leaderboard`;
		axios
			.get(url)
			.then(res => {
				this.setState({
					leaderboard: res.data
				});
			})
			.catch(error => {
				console.log(error);
			});
	};

	public render() {
		var button;
		if (this.state.available) {
			button = (
				<button onClick={e => this.postTypie()} disabled={this.state.disabled}>
					PLAY
				</button>
			);
		} else {
			button = (
				<button disabled={true} onClick={e => this.postTypie()}>
					PLAY
				</button>
			);
		}
		return (
			<div>
				<h1>Hello, New Typies</h1>
				<input type="text" ref="username" id="username" />
				{button}
				{this.renderTable()}
			</div>
		);
	}

	public componentWillMount() {
		//establish websocket collection
		this.fetchTopScores();
		let host = this.getCurrentHost();

		const websocket = new WebSocket('ws://' + host + '/ws');

		websocket.addEventListener('error', function(err) {
			console.log(err);
		});
		websocket.addEventListener('open', function() {
			console.log('Websocket connection established');
		});
		websocket.addEventListener('close', function() {
			console.log('Websocket connection closed');
		});
		websocket.addEventListener('message', event => {
			const topScores = JSON.parse(event.data);
			console.log(topScores);
			if (topScores.type == 'Leaderboard') {
				console.log(topScores.payload.available);
				this.setState({
					available: topScores.payload.available,
					leaderboard: topScores.payload.leaders
				});
			}
		});
	}

	private renderTable = (): JSX.Element => {
		const thead = (
			<thead>
				<tr>
					<td>
						<h3>
							{this.state.available ? 'Gameroom: Available' : 'Gameroom: Unavailable'}
						</h3>
					</td>
				</tr>
				<tr>
					<th>Rank</th>
					<th>Username</th>
					<th>Best Record</th>
				</tr>
			</thead>
		);
		var count = 0;
		var scores = this.state.leaderboard.map(leader => {
			count++;
			return (
				<tr key={count}>
					<td>{count}</td>
					<td>{leader.userName}</td>
					<td>{leader.record}</td>
				</tr>
			);
		});

		const tbody = <tbody>{scores}</tbody>;
		return (
			<table>
				{thead}
				{tbody}
			</table>
		);
	};

	private postTypie = () => {
		this.setState({ disabled: true });
		console.log('fuck');
		let username = this.refs.username['value'].trim();
		if (!username) {
			return;
		}

		let typie = {
			userName: username
		};
		const url = `http://${this.getCurrentHost()}/typie`;
		axios
			.post(url, typie)
			.then(res => {
				localStorage.setItem('TypieID', res.data.id);
				window.location.replace('app.html');
			})
			.catch(error => {
				window.alert(error);
			});
	};

	private renderTableData = (typie): JSX.Element => {
		return (
			<tr>
				<td>{typie.userName}</td>
				<td>{typie.record}</td>
			</tr>
		);
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
}

ReactDOM.render(<Index />, document.getElementById('index'));
