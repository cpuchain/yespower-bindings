import { webcrypto as crypto } from 'crypto';
import { yespower } from '../index.js';

interface TestCases {
    input: string;
    output: string;
    N?: number;
    r?: number;
    pers?: string;
}

function getRandomBuffer() {
    return Buffer.from(crypto.getRandomValues(new Uint8Array(32)));
}

const pers = 'pers';

const num = 5;

function createCases() {
    const cases: TestCases[] = [];

    for (let i = 0; i < num; ++i) {
        const input = getRandomBuffer();
        const output = yespower(input, undefined, undefined, pers);

        cases.push({
            input: input.toString('hex'),
            output: output.toString('hex'),
            pers,
        });
    }

    console.log(cases);
}

createCases();
