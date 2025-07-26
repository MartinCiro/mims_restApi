module.exports = {
  apps: [{
    name: 'api',
    script: 'dist/src/index.js',
    interpreter: "node",
    node_args: '',  // <- Eliminado tsconfig-paths
    env: {
      NODE_ENV: 'production',
      // TS_CONFIG_PATHS ya no es necesario
    },
    error_file: 'logs/err.log',
    out_file: 'logs/out.log',
    log_date_format: 'YYYY-MM-DD HH:mm:ss'
  }]
};