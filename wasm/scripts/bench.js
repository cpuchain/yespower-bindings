import { Yespower } from '../lib/index.js';

const tests = 30;

const N = 1024;
const r = 8;

async function bench() {
    const yespower = await Yespower.init();

    const timeStart = Date.now();
    
    for (let i = 0; i < tests; ++i) {
        const buf = Buffer.allocUnsafe(4);
        buf.writeUint32BE(i)
        yespower.Hash(buf, N, r);
    }

    const timeTook = Date.now() - timeStart;
    const hps = Math.floor(1000 * tests / timeTook);

    console.log(`${tests} Tests: ${Date.now() - timeStart}ms (${hps}H/s)`);
}
bench();