import * as React from 'react';
import * as ReactDOM from 'react-dom';
import axios from 'axios';

import 'sass/index';

class Index extends React.Component<any, any> {
	constructor(props, context) {
		super(props, context);

		this.state = {
			typies: new Array(),
			available: false
		};
	}

	public render() {
		return (
			<div>
				<h1>Hello, New Typies</h1>
				<input type="text" ref="username" id="username" />
				<button onClick={e => this.postTypie()}>PLAY</button>
				{this.state.available ? <h3>Gameroom Open!</h3> : null}
				<table style={{ width: 100 }}>
					<thead>
					<tr>
						<th>Username</th>
						<th>Time</th>
					</tr>
					</thead>
					<tbody>
						{this.state.typies.forEach(element => {
							this.renderTableData(element);
						})}
					</tbody>
				</table>
			</div>
		);
	}

	public componentWillMount() {
		//establish websocket collection
		let host = this.getCurrentHost();

		const websocket = new WebSocket('ws://' + host + '/ws');
		websocket.addEventListener('error', function(err) {
			window.alert(err);
		});
		websocket.addEventListener('message', event => {
			const highScores = JSON.parse(event.data);
			this.setState({
				leaderboard: highScores
			});
		});
	}

	private postTypie = () => {
		let username = this.refs.username['value'].trim();
		let typie = {
			userName: ''
		};
		if (!username) {
			return;
		}
		typie.userName = username;

		const url = `http://${this.getCurrentHost()}/typie`;
		axios
			.post(url, typie)
			.then(res => {
				localStorage.setItem('TypieID', JSON.stringify(res.data.id));
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
