const path = require('path');

const config = {
  mode: 'development',
  entry: './src/app.js',
  output: {
    filename: 'app.js',
    path: path.resolve(__dirname, '../public/dist'),
  },
  target: ['web', 'es5'],
  module: {
    rules: [
      {
        test: /\.js|\.jsx$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: [
              ['@babel/preset-env', {"targets": "defaults"}],
              '@babel/preset-react'
            ]
          }
        }
      },
      {
        test: /\.css|\.s[ac]ss$/i,
        use: ['style-loader', 'css-loader', 'sass-loader'],
      },
    ],
  },
  watchOptions: {
    poll: true,
  }
};

module.exports = (env, argv) => {
  if (argv.mode) {
    config.mode = argv.mode;
  }

  return config;
};