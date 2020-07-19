import babel from 'rollup-plugin-babel';
import commonJS from 'rollup-plugin-commonjs';
import nodeResolve from 'rollup-plugin-node-resolve';
import { uglify } from 'rollup-plugin-uglify';

module.exports = {
  input: 'src/main.js',
  output: {
    name: 'bundle',
    format: 'iife',
    file: 'assets/js/ws-client.min.js'
  },
  plugins: [
    babel({
      exclude: 'node_modules/**'
    }),
    nodeResolve({ browser: true }),
    commonJS({
      include: './node_modules/**'
    }),
    uglify()
  ],
}