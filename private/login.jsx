'use strict';

var React = require('react');

var Forms = require('./forms.jsx');
var Auth = require('./auth.jsx');

var LoginForm = React.createClass({
    handleSubmit: function(e) {
	e.preventDefault();
	Auth.authorize(this.refs.handle.value(), this.refs.password.value()).then(function(data) {
	    if (data["status"] == "success") {
		window.location = 'hello.html';
	    } else {
		this.setState(data["errors"]);
	    }
	});
    },
    getInitialState: function() {
	return {
	    handle: '',
	    password: ''
	};
    },
    render: function() {
	return (
	    <div class="LoginForm">
		<form role="form" class="signup" onSubmit={this.handleSubmit} >
		<Forms.TextInputWithFeedback field="handle" errors={ this.state.handle } ref="handle"/>
		<Forms.TextInputWithFeedback field="password" errors={ this.state.password } ref="password" />

		<button type="submit" class="btn btn-primary">Login</button>
		</form>
	    </div>
	);
    }
});

React.render(
  <div class="content"><LoginForm /></div>,
  document.getElementById('content')
);
