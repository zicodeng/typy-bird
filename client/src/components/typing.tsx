import axios from 'axios';
import * as React from 'react';

import './typing.scss';

class Typing extends React.Component<any, any> {
	constructor(props, context) {
		super(props, context);

		this.state = {
			dictionary: ['Hello', 'World', 'How', 'Are', 'You'],
			currentWordIndex: 0
		};
	}

	public render(): JSX.Element {
		return (
			<div className="typing">
				<p>{this.state.dictionary[this.state.currentWordIndex]}</p>
				<input type="text" ref="word" autoFocus onChange={e => this.handleChangeInput()} />
			</div>
		);
	}

	private handleChangeInput = (): void => {
		let currentWordIndex = this.state.currentWordIndex;
		const currentWord = this.state.dictionary[currentWordIndex];
		const word = this.refs.word['value'].trim();
		// If the user typed word matches the currentWord that prompts the user to type
		// send a request to update this Typie's game state,
		// clear input area,
		// update currentWordIndex by 1,
		// and prompt the user with a new currentWord.
		if (word === currentWord) {
			console.log('Update Typie');

			// const host = this.props.host;
			// const sessionToken = this.props.sessionToken;
			// const url = `https://${host}/v1/messages/${selectedMessage._id}`;
			// axios
			// 	.patch(url, message, {
			// 		headers: {
			// 			Authorization: sessionToken
			// 		}
			// 	})
			// 	.then(res => {
			// 		this.refs.messageBody['value'] = '';
			// 		this.setState({
			// 			messageMode: this.MESSAGE_MODE.CREATE
			// 		});
			// 	})
			// 	.catch(error => {
			// 		window.alert(error.response.data);
			// 	});

			this.refs.word['value'] = '';
			currentWordIndex++;
			this.setState({
				currentWordIndex: currentWordIndex
			});
		}
	};
}

export default Typing;
