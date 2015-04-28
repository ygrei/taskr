'use strict';

var Forms = require('./forms.jsx');
var React = require('react');

var SignupForm = React.createClass({
    handleSubmit: function(e) {
	e.preventDefault();
	payload = {
	    'username': this.refs.username.value(),
	    'password': this.refs.password.value(),
	    'email': this.refs.email.value()
	};
	console.log(payload);
	$.ajax({
            url: "/json/signup",
            dataType: 'json',
            type: 'POST',
            data: payload,
            success: function(data) {
		if (data["status"] == "invalid") {
		    console.log("Setting status");
		    console.log( data["errors"] );
		    this.setState( data["errors"] );
		} else {

		}
            }.bind(this),
            error: function(xhr, status, err) {
		console.error(this.props.url, status, err.toString());
            }.bind(this)
	});
    },
    getInitialState: function() {
	return {
	    usernameError: '',
	    emailError: '',
	    passwordError: ''
	};
    },
    render: function() {
	return (
	    <div class="SignupForm">
		<form role="form" class="signup" onSubmit={this.handleSubmit} >
		<Forms.TextInputWithFeedback field="username" errors={ this.state.usernameError } ref="username"/>
		<Forms.TextInputWithFeedback field="email" errors={ this.state.emailError } ref="email" />
		<Forms.TextInputWithFeedback field="password" errors={ this.state.passwordError } ref="password" />

		<button type="submit" class="btn btn-primary">Sign up</button>
		</form>
	    </div>
	);
    }
});

React.render(
  <SignupForm />,
  document.getElementById('content')
);
