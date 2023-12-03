// const { defineConfig } = require('@vue/cli-service')
// module.exports = defineConfig({
//   transpileDependencies: true,
//   lintOnSave: false
// })

// module.exports ={
//   devServer:{
//     host: '0.0.0.0',
//     port: 8080,
//     open: true},
//   lintOnSave: false,
// }

const NodePolyfillPlugin = require('node-polyfill-webpack-plugin');

module.exports = {
  devServer: {
    host: '0.0.0.0',
    port: 8080,
    open: true
  },
  lintOnSave: false,
  configureWebpack: {
    plugins: [
      new NodePolyfillPlugin()
    ]
  }
};
