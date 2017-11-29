import * as React from 'react';
import * as ReactDOM from 'react-dom';

import 'sass/index';

class Index extends React.Component<any, any> {
    constructor(props, context) {
        super(props, context);
    }

    public render() {
        return (
            <div>
                <h1>Hello, New Typies</h1>
                <a href="./app.html">PLAY GAME</a>
            </div>
        );
    }
}

ReactDOM.render(<Index />, document.getElementById('index'));
