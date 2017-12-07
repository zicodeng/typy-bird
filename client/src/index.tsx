import * as React from 'react';
import * as ReactDOM from 'react-dom';
import axios from 'axios';

import 'sass/index';

class Index extends React.Component<any, any> {
	constructor(props, context) {
		super(props, context);

		this.state = {
			gameRoom: null
		};
	}

	public render() {
		return (
			<div>
				<h1>Hello, New Typies</h1>
				<input type="text" ref="username" id="username" />
				<button onClick={e => this.postTypie()}>PLAY</button>
				{this.renderTable()}
			</div>
		);
	}

	public componentWillMount() {
		//establish websocket collection
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
			const gameRoom = JSON.parse(event.data);
			this.setState({
				gameRoom: gameRoom
			});
		});
	}

	private renderTable = (): JSX.Element => {
		const thead = (
			<thead>
				<tr>
					<th>Username</th>
					<th>Best Record</th>
				</tr>
			</thead>
		);
		const rows = {};
		const tbody = <tbody>{rows}</tbody>;
		return <table>{thead}</table>;
	};

	private postTypie = () => {
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
