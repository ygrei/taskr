var Forms = require('./forms.jsx');
var React = require('react');

var LoginForm = React.createClass({
    handleSubmit: function(e) {
	e.preventDefault();
	payload = {
	    'handle': this.refs.handle.value(),
	    'password': this.refs.password.value(),
	};
	console.log(payload);
	$.ajax({
            url: "/json/login",
            dataType: 'json',
            type: 'POST',
            data: payload,
            success: function(data) {
		if (data["status"] == "invalid") {
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
	    handleError: '',
	    passwordError: ''
	};
    },
    render: function() {
	return (
	    <div class="LoginForm">
		<form role="form" class="signup" onSubmit={this.handleSubmit} >
		<Forms.TextInputWithFeedback field="handle" errors={ this.state.handleError } ref="handle"/>
		<Forms.TextInputWithFeedback field="password" errors={ this.state.passwordError } ref="password" />

		<button type="submit" class="btn btn-primary">Sign up</button>
		</form>
	    </div>
	);
    }
});

React.render(
  <LoginForm />,
  document.getElementById('content')
);
