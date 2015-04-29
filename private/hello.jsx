'use strict';

var React = require('react');
const store = require('store');

var handle = store.get('handle');

React.render(
  <p>hello {handle}</p>,
  document.getElementById('content')
);
