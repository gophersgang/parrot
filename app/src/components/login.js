import React, { PropTypes } from 'react';
import TextField from 'material-ui/TextField';
import RaisedButton from 'material-ui/RaisedButton';

class Login extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                <TextField
                    hintText="Your email"
                    floatingLabelText="Email"
                /><br />
                <TextField
                    hintText="Your password"
                    floatingLabelText="Password"
                    type="password"
                /><br />
                <RaisedButton
                    label="Login"
                    primary={true}
                    onClick={this.props.onSubmit}
                />
            </div>
        );
    }
}

export default Login;