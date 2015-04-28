var webpack = require('webpack');
var path = require('path');

module.exports = {
  entry: ['./private/signup.jsx'],
  output: {
    path:  __dirname + '/public',
    filename: 'scripts/signup.js'
  },
  plugins: [
    new webpack.NoErrorsPlugin()
  ],
  resolve: {
    extensions: ['', '.js', '.jsx']
  },
  module: {
    loaders: [
      { test: /\.jsx?$/, loaders: ['react-hot', 'babel'], exclude: /node_modules/ },
    ]
  }
};
