import { createRequire } from 'node:module';
const require = createRequire(import.meta.url);
const addon = require('./build/Release/yespower.node');

export const yespower = addon.yespower;

// Asynchronous variant: runs the CPU-heavy hash on the libuv threadpool and
// returns a Promise, so the Node.js event loop is not blocked.
export const yespower_async = addon.yespower_async;
