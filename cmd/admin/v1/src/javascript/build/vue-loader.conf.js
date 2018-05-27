'use strict';
const utils = require('./utils');
const config = require('../config');
const isProduction = process.env.NODE_ENV === 'production';

module.exports = {
    // // only enable if doing custom css
    // loaders: utils.cssLoaders({
    //     sourceMap: isProduction
    //         ? config.build.productionSourceMap
    //         : config.dev.cssSourceMap,
    //     extract: isProduction,
    // }),
    transformToRequire: {
        video: 'src',
        source: 'src',
        img: 'src',
        image: 'xlink:href',
    },
};
