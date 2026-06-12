import { describe, before, it } from 'node:test';
import { strict as assert } from 'assert';
import { Yespower } from '../src/index.js';
import { bytesToHex } from '../src/utils.js';

// ---------------------------------------------------------------------------
// Reference vectors from yespower-c/TESTS-OK, reproduced exactly the way the C
// test driver (yespower-c/tests.c) does. See node/test/index.ts for details.
// ---------------------------------------------------------------------------

// Replicates `src` in tests.c: src[i] = i * 3 (mod 256). 80 bytes.
function refSrc(): Uint8Array {
    const src = new Uint8Array(80);
    for (let i = 0; i < src.length; i++) {
        src[i] = (i * 3) & 0xff;
    }
    return src;
}

const SRC = refSrc();

// Special sentinel meaning "use SRC itself as the personality" (the BSTY case).
const BSTY = Symbol('BSTY');

interface RefCase {
    version: 5 | 10;
    N: number;
    r: number;
    pers: string | null | typeof BSTY;
    output: string;
}

const refCases: RefCase[] = [
    {
        version: 5,
        N: 2048,
        r: 8,
        pers: 'Client Key',
        output: 'a59fec4c4fdda16e3b1405adda66d525b68e7cadfcfe6ac066c7ad118cd80590',
    },
    {
        version: 5,
        N: 2048,
        r: 8,
        pers: BSTY,
        output: '5ea2b2956a9eace30a3237ff1d441edee1dc25aab8f0ea15c12165f83a7bc265',
    },
    {
        version: 5,
        N: 4096,
        r: 16,
        pers: 'Client Key',
        output: '927e72d0ded3d80475473f40f1743c67289d453d5242d4f55af4e325e06699c5',
    },
    {
        version: 5,
        N: 4096,
        r: 24,
        pers: 'Jagaricoin',
        output: '0e1366973211e7fea8ad9d81989c84a254d968c9d333dd8ff099324f38611e04',
    },
    {
        version: 5,
        N: 4096,
        r: 32,
        pers: 'WaviBanana',
        output: '3ae05abb3c5cf6f75415a92554c98d50e38ec9552cfa78373616f480b24e559f',
    },
    {
        version: 5,
        N: 2048,
        r: 32,
        pers: 'Client Key',
        output: '560a891b5ca2e1c636111a9ff7c894a5d0a2602f43fdcfa5949b95e22fe4461e',
    },
    {
        version: 5,
        N: 1024,
        r: 32,
        pers: 'Client Key',
        output: '2a79e53d1be6669bc556ccc417bce3d22a74a232f56b8e1d39b45792675de108',
    },
    {
        version: 5,
        N: 2048,
        r: 8,
        pers: null,
        output: '5ecbd8e8d7c90baed4bbf8916a1225dcc3c65f5c9165bae81cdde3cffad128e8',
    },

    {
        version: 10,
        N: 2048,
        r: 8,
        pers: null,
        output: '69e0e895b3df7aeeb837d71fe199e9d34f7ec46ecbca7a2c4308e51857ae9b46',
    },
    {
        version: 10,
        N: 4096,
        r: 16,
        pers: null,
        output: '33fb8f063824a4a020f63dca535f5ca66ab5576468c75d1ccaac7542f76495ac',
    },
    {
        version: 10,
        N: 4096,
        r: 32,
        pers: null,
        output: '771aeefda8fe79a0825bc7f2aee162ab5578574639ffc6ca3723cc18e5e3e285',
    },
    {
        version: 10,
        N: 2048,
        r: 32,
        pers: null,
        output: 'd5efb813cd263e9b34540130233cbbc6a921fbff3431e5ec1a1abde2aea6ff4d',
    },
    {
        version: 10,
        N: 1024,
        r: 32,
        pers: null,
        output: '501b792db42e388f6e7d453c95d03a12a36016a5154a688390ddc609a40c6799',
    },
    {
        version: 10,
        N: 1024,
        r: 32,
        pers: 'personality test',
        output: '1f0269acf565c49adc0ef9b8f26ab3808cdc38394a254fddeedcc3aacff6ad9d',
    },
];

const xorV5 = 'aef13291870f557047f42e9befa616dfe5f19677e13f8ba692f7c59755a0f50e';
const xorV10 = '8d13c5fb07309675d1b8489277ba4be44033bedfae7a60438a9be21f3a7b1237';

const hex = (bytes: Uint8Array) => bytesToHex(bytes).replace('0x', '');

describe('Yespower (WASM) — yespower-c/TESTS-OK reference vectors', () => {
    let yespower: Yespower;

    before(async () => {
        yespower = await Yespower.init();
    });

    it('single-hash vectors (v0.5 + v1.0, incl. binary BSTY pers)', () => {
        for (const c of refCases) {
            const pers = c.pers === BSTY ? SRC : c.pers === null ? '' : c.pers;
            const hashed = yespower.Hash(SRC, c.N, c.r, pers, c.version);
            assert.strictEqual(
                hex(hashed),
                c.output,
                `yespower(${c.version}, ${c.N}, ${c.r}, ${String(c.pers)})`,
            );
        }
    });

    it('XOR of yespower(5, ...) aggregate', () => {
        const src = refSrc();
        src[0] = 43;
        const xor = new Uint8Array(32);
        for (let N = 1024; N <= 4096; N <<= 1) {
            for (let r = 8; r <= 32; r++) {
                const dst = yespower.Hash(src, N, r, 'Client Key', 5);
                for (let i = 0; i < 32; i++) {
                    xor[i] ^= dst[i];
                }
            }
        }
        assert.strictEqual(hex(xor), xorV5);
    });

    it('XOR of yespower(10, ...) aggregate', () => {
        const src = refSrc();
        src[0] = 43;
        const xor = new Uint8Array(32);
        for (let N = 1024; N <= 4096; N <<= 1) {
            for (let r = 8; r <= 32; r++) {
                const dst = yespower.Hash(src, N, r, '', 10);
                for (let i = 0; i < 32; i++) {
                    xor[i] ^= dst[i];
                }
            }
        }
        assert.strictEqual(hex(xor), xorV10);
    });
});

// ---------------------------------------------------------------------------
// Additional fixed-input vectors (v1.0 defaults N=2048, r=32) kept from the
// original test suite, including personality="pers" cases.
// ---------------------------------------------------------------------------

interface TestCases {
    input: string;
    output: string;
    N?: number;
    r?: number;
    pers?: string;
}

const cases: TestCases[] = [
    {
        input: 'eebb7bf9a8c813b5e0a03ce627bd1a0c836e0a89793743666dc82b83e28e8f00',
        output: '07d0c37029360872e56e95873d55dca1424e45b65266e7ec4f992625c5f836d7',
    },
    {
        input: '6c0d7c5f367d5d39b723965eb4ac7d37188c5af208a7a38813ace92f4567d5ce',
        output: 'f475c6e007ac79bf53388b8065c5c55c9fccbca9d9923c61b2f0fec3454c333f',
    },
    {
        input: 'f82a3ad1a4615c725d17246e87163abd0a08d01c13b5e349795f99fb579e61ea',
        output: 'ca2699030b8ac88f4b390d174b1fb5a3d262f608cc604eac01ccb03d8b809c38',
    },
    {
        input: 'e85408220029f42d62c6543a4101e593c9a0bf4c10b99731c81c77b717ddfb83',
        output: '8c489bf47cf955fe5083ec68fdb5c85bb00847a3674c8514db3ec4e7bf14aed3',
    },
    {
        input: 'f881d0e7f2d6bd621ac83f39aa29778a5c0173fbfcde4e83581b1f02c6d94704',
        output: 'eaa5a2db5ab756ff9179e45956661d60045d2a82507ea93dfb09e80554b6f127',
    },
    {
        input: '82128c4a7eef3a079eb797909f7108a0d4e0c2cd91a875d3d9739cf02ff211f8',
        output: 'd068f5bfb5a5e5a67a296d8f65f1775e22e3ff6c60fbc17c8dace4aa9ed8624a',
    },
    {
        input: '2538dad623dab29a3d5387804ab51dea411014fad9c47fb94b2d83e44358064b',
        output: '62c4ac19375787857e2e7c41282a58fb68638a25ba27402239edc8d2683b5126',
        pers: 'pers',
    },
    {
        input: '28d5ce4e064ddc5ecf7a62a65d245facb3b46d6d2d70cdcf5c413e1c4e205157',
        output: '9f82d3a2d796ac89aa012c71ccad798f5e10cf515dd6fb0fedceb44271029d43',
        pers: 'pers',
    },
    {
        input: '172aaa13d4b7d606134bc6989a84c403bac3b6c4f9f4504dd79b2618c11c63d3',
        output: '5bd719de97d5dbc0012bd576f778994ae9d3eeb31e12cfadeb58133a34c2ebef',
        pers: 'pers',
    },
    {
        input: 'c9d4866805b9e9788cf1c7f6a369e8932e5bc6e32a5c7a06c3ecb7e0ec8ebaff',
        output: 'ac22b64773effd8d0dbb088bd71b27449ec5c8c21ecc0262f615e920ea2df2a6',
        pers: 'pers',
    },
    {
        input: 'fae4f16aed075d1e94b6d913242840ad3933bc675aa72f2f647f1df60c431963',
        output: 'a625f58d4585b8bad21509f34f715f6a15a5e2741a89642e937587b5c4b68e8a',
        pers: 'pers',
    },
];

describe('Yespower (WASM) — additional fixed-input vectors', () => {
    let yespower: Yespower;

    before(async () => {
        yespower = await Yespower.init();
    });

    it('Test Cases', () => {
        for (const { input, output, N, r, pers } of cases) {
            const hashed = yespower.Hash(Buffer.from(input, 'hex'), N, r, pers);
            assert.strictEqual(output, hex(hashed));
        }
    });
});
