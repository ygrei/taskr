'use strict';

var React = require('react');

var TextInputWithFeedback = React.createClass({
    value: function() {
	return React.findDOMNode(this.refs.param).value.trim();
    },
    render: function() {
	var field = this.props.field.toString();
	var errors = this.props.errors.toString();
	return (
	     <div class={ "form-group" + errors ? "has-error has-feedback" : "" } onSubmit={this.props.callback} >
		<label class="control-label" for={ field }>{{ field }}</label>
		<input type="text" ref="param" class="form-control" id={ field } name={ field } placeholder={"Enter " + field } />
		<span class="help-block"> { errors } </span>
	    </div>
	);
    }
});

module.exports = {TextInputWithFeedback: TextInputWithFeedback};
