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
		if (localStorage.getItem('TypieID')) {
			const typieID = localStorage.getItem('TypieID');
			if (typieID) {
				const url = `http://${this.getCurrentHost()}/gameroom?auth=${typieID}`;
				axios.delete(url).catch(error => {
					console.log(error);
				});

				localStorage.clear();
			}
		}
		return (
			<div className="container">
				<h1 className="greeting">Hello, New Typies!</h1>
				<h4 className="gameroom-status">
					Game Room Status:&nbsp;
					{this.state.available ? <span>Available</span> : <span>Unavailable</span>}
				</h4>
				<input type="text" ref="username" id="username" placeholder="Username" />
				{this.renderButtons()}
				<h4>History Records</h4>
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
			const data = JSON.parse(event.data);
			console.log(data);
			switch (data.type) {
				case 'Leaderboard':
					this.setState({
						available: data.payload.available,
						leaderboard: data.payload.leaders
					});
					break;

				case 'GameStart':
					this.setState({
						available: false
					});
					break;

				case 'GameEnd':
					this.setState({
						available: true
					});
					break;

				default:
					break;
			}
		});
	}

	private renderButtons = (): JSX.Element => {
		if (this.state.available) {
			return (
				<button onClick={e => this.postTypie()} disabled={this.state.disabled}>
					PLAY
				</button>
			);
		} else {
			return <button>WAIT</button>;
		}
	};

	private renderTable = (): JSX.Element => {
		const thead = (
			<thead>
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
