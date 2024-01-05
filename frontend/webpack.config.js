// webpack.config.js
const path = require('path');

module.exports = {
  mode: 'development',  // or 'production' or 'none'
  entry: './script.js',  // adjust the path based on your project structure
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist'),  // adjust the path based on your project structure
  },
};
