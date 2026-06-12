/**
 * Example CPU mining script with Core-Geth eth_getWork RPC
 *
 * byte functions from ethereumjs
 */
import { yespower } from '../index.js';

const RPC = 'http://localhost:8545';

const ZeroHash = '0x0000000000000000000000000000000000000000000000000000000000000000';

const concatBytes = (...arrays) => {
    if (arrays.length === 1) return arrays[0];
    const length = arrays.reduce((a, arr) => a + arr.length, 0);
    const result = new Uint8Array(length);
    for (let i = 0, pad = 0; i < arrays.length; i++) {
        const arr = arrays[i];
        result.set(arr, pad);
        pad += arr.length;
    }
    return result;
};

const setLength = (msg, length, right) => {
    if (right) {
        if (msg.length < length) {
            return new Uint8Array([...msg, ...new Uint8Array(length - msg.length)]);
        }
        return msg.subarray(0, length);
    } else {
        if (msg.length < length) {
            return new Uint8Array([...new Uint8Array(length - msg.length), ...msg]);
        }
        return msg.subarray(-length);
    }
};

const hexToBytes = (hexStr) => {
    let cpyStr = typeof hexStr === 'bigint' || typeof hexStr === 'number' ? hexStr.toString(16) : hexStr;
    if (cpyStr.startsWith('0x')) {
        cpyStr = cpyStr.slice(2);
    }
    if (cpyStr.length % 2 !== 0) {
        cpyStr = '0' + cpyStr;
    }
    return Uint8Array.from(cpyStr.match(/.{1,2}/g).map((byte) => parseInt(byte, 16)));
};

const bytesToHex = (bytes) => {
    return (
        '0x' +
        Array.from(bytes)
            .map((b) => b.toString(16).padStart(2, '0'))
            .join('')
    );
};

function bytesReverse(a) {
    const length = a.length;
    const b = new Uint8Array(length);
    for (let i = 0; i < length; i++) {
        b[i] = a[length - i - 1];
    }
    return b;
}

async function getWork() {
    const resp = await (
        await fetch(RPC, {
            method: 'POST',
            headers: {
                'content-type': 'application/json',
            },
            body: JSON.stringify({
                jsonrpc: '2.0',
                method: 'eth_getWork',
                params: [],
                id: 1,
            }),
        })
    ).json();

    const newWork = {
        header: resp.result[0],
        seed: resp.result[1],
        target: resp.result[2],
        difficulty: 2n ** 256n / BigInt(resp.result[2]),
        number: BigInt(resp.result[3]),
        blockNumber: resp.result[3],
    };

    return newWork;
}

async function submitWork(solution) {
    const resp = await (
        await fetch(RPC, {
            method: 'POST',
            headers: {
                'content-type': 'application/json',
            },
            body: JSON.stringify({
                jsonrpc: '2.0',
                method: 'eth_submitWork',
                params: [solution.nonce, solution.header, solution.mixHash],
                id: 1,
            }),
        })
    ).json();

    console.log(resp);
}

async function test() {
    const work = await getWork();

    console.log(work);

    let nonce = 0;

    const target = BigInt(work.target);

    while (true) {
        const { header } = work;
        const nonceBytes = setLength(hexToBytes(nonce), 8, false);
        const concatted = concatBytes(hexToBytes(header), bytesReverse(nonceBytes));
        const result = BigInt('0x' + yespower(Buffer.from(concatted)).toString('hex'));

        if (target > result) {
            const solution = {
                header: work.header,
                mixHash: ZeroHash,
                nonce: bytesToHex(nonceBytes),
            };

            console.log(solution);

            await submitWork(solution);

            break;
        }

        nonce++;
    }
}

test();
