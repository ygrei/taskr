'use strict';

const store = require('store');

function authorize(handle, password) {
    var payload = {
	'handle': handle,
	'password': password
    };
    var result = $.getJSON('/json/login', payload).then(function(data) {
	if (data["status"] == "success") {
	    store.set('handle', payload['handle']);
	    store.set('password', payload['password']);
	} else {
	    store.remove('handle');
	    store.remove('password');
	}
	return data;
    });
    return result;
}

module.exports = {authorize: authorize};
