module.exports = {
    plugins: [
      require('postcss-preset-env')({
        stage: 3,
        features: {
          'nesting-rules': true, 
          'custom-properties': true, 
        },
      }),
      require('autoprefixer')({
        overrideBrowserslist: ['> 1%', 'last 2 versions'], 
      }),
      require('cssnano')({ 
        preset: 'default',
      }),
    ],
  };