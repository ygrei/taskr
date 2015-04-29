var webpack = require('webpack');
var path = require('path');

module.exports = {
  entry: {
      'signup': './private/signup.jsx',
      'login': './private/login.jsx',
      'hello': './private/hello.jsx',
  },
  output: {
    path:  __dirname + '/public',
    filename: 'scripts/[name].js'
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
