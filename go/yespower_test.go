package yespower

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

type testCase struct {
	input  string
	output string
	N      uint32
	r      uint32
	pers   string
}

const refInput string = "000306090c0f1215181b1e2124272a2d303336393c3f4245484b4e5154575a5d606366696c6f7275787b7e8184878a8d909396999c9fa2a5a8abaeb1b4b7babdc0c3c6c9cccfd2d5d8dbdee1e4e7eaed"

func getTestCases() []testCase {
	testCaseSlice := []testCase{
		{
			input: refInput,
			output: "d5efb813cd263e9b34540130233cbbc6a921fbff3431e5ec1a1abde2aea6ff4d",
		},
		{
			input: refInput,
			output: "69e0e895b3df7aeeb837d71fe199e9d34f7ec46ecbca7a2c4308e51857ae9b46",
			N: 2048,
			r: 8,
		},
		{
			input: refInput,
			output: "33fb8f063824a4a020f63dca535f5ca66ab5576468c75d1ccaac7542f76495ac",
			N: 4096,
			r: 16,
		},
		{
			input: refInput,
			output: "771aeefda8fe79a0825bc7f2aee162ab5578574639ffc6ca3723cc18e5e3e285",
			N: 4096,
			r: 32,
		},
		{
			input: refInput,
			output: "501b792db42e388f6e7d453c95d03a12a36016a5154a688390ddc609a40c6799",
			N: 1024,
			r: 32,
		},
		{
			input: refInput,
			output: "1f0269acf565c49adc0ef9b8f26ab3808cdc38394a254fddeedcc3aacff6ad9d",
			N: 1024,
			r: 32,
			pers: "personality test",
		},
		{
			input:  "eebb7bf9a8c813b5e0a03ce627bd1a0c836e0a89793743666dc82b83e28e8f00",
			output: "07d0c37029360872e56e95873d55dca1424e45b65266e7ec4f992625c5f836d7",
		},
		{
			input:  "6c0d7c5f367d5d39b723965eb4ac7d37188c5af208a7a38813ace92f4567d5ce",
			output: "f475c6e007ac79bf53388b8065c5c55c9fccbca9d9923c61b2f0fec3454c333f",
		},
		{
			input:  "f82a3ad1a4615c725d17246e87163abd0a08d01c13b5e349795f99fb579e61ea",
			output: "ca2699030b8ac88f4b390d174b1fb5a3d262f608cc604eac01ccb03d8b809c38",
		},
		{
			input:  "e85408220029f42d62c6543a4101e593c9a0bf4c10b99731c81c77b717ddfb83",
			output: "8c489bf47cf955fe5083ec68fdb5c85bb00847a3674c8514db3ec4e7bf14aed3",
		},
		{
			input:  "f881d0e7f2d6bd621ac83f39aa29778a5c0173fbfcde4e83581b1f02c6d94704",
			output: "eaa5a2db5ab756ff9179e45956661d60045d2a82507ea93dfb09e80554b6f127",
		},
		{
			input:  "82128c4a7eef3a079eb797909f7108a0d4e0c2cd91a875d3d9739cf02ff211f8",
			output: "d068f5bfb5a5e5a67a296d8f65f1775e22e3ff6c60fbc17c8dace4aa9ed8624a",
		},
		{
			input:  "04e42c070dddea438dc6f06a0dc596d184caa3a324fbbbdf91bdd9d7e1835a00",
			output: "54913feaf7f91bfd7e39985b7ca7a53062a4584b4a72c25ab0db2380cad65748",
		},
		{
			input:  "f0eba2e16c7b8527c33982a238de5b06625b637260ed9bb18847a4fb2c82aa84",
			output: "55901209b07c5fd580df3c5d62f4f499c6024349ff151fe3efbcd7cdf06d6ded",
		},
		{
			input:  "0d65680a44d784a8e4511c03bc039f9498dbea6edd0a17d655f1cc4e3c170567",
			output: "6398153b75cead0412096eef964095bb6a94d9a46687b2c54674f3aa4fd6b783",
		},
		{
			input:  "eba9c0c98615bef5364653b44f40b5c051056592b76dafcc47e9114789afc98e",
			output: "c74e2a86266870abddee6e0f15437268e3cff57bd50730d0ac6ccf03449b9837",
		},
		{
			input:  "fa0bda752b182be3036c07139a62d2c89618702909b1f9f84645cf943b32ff97",
			output: "ee235a8b499d18a9c6499d11de3c60d8cb88c19ab31c5dcd5a36f1f10845633a",
		},
		{
			input:  "9889006e194dcb39439b0463d2efb9665a3c51195e041d5ee605e1247701b542",
			output: "2dd667250910ce096d86e8bc5fe784e649a8bb2706631f1ff74bb103a80bef7b",
		},
		{
			input:  "79a8487674e136c7a5169289265db025910f09a769bbb99a104e8707d3e9c3e5",
			output: "029e51eb819b76dac6e861e8123717f8e608afa86b48440302a160c0b9952fc6",
		},
		{
			input:  "b1bff94d41f07a53c88bfdd8059efb5c11616b56e738ec9200078b0f07067457",
			output: "6803652aa7bb1aa756e5e7647c103d70ce5a3bedd0aaea0140719c524263a539",
		},
		{
			input:  "ae6448448c97f6675ef31e167346df22502e9a8b06124d11113d75800d5c1d39",
			output: "e74c05b434cb559857af4c43efb9c96d4de591fa754eb49616623d3387680021",
		},
		{
			input:  "87529030f04da3e6e0685b921871557ace38bf252dc4d3ba090c4aae6cade29f",
			output: "80cd3ebee02cf7642743623ff03c5dbe4491506db7288bdce853404cef07fb11",
		},
		{
			input:  "74e762803a284a65aeded8984fdf3c7850f2b4fc4e76cf05ce692e76673fb005",
			output: "7c1cb6f3789dc3c839f6768fa803b163d544bc047ac4d88fb770e8c578bf7109",
		},
		{
			input:  "261aacbae7b3de7491dc98d85c4774fe17f74e3b130520b9ef9ec8fe70f0d93c",
			output: "64d267ac7e5470efd4b5ce64ed3a5a758c6c11ccb01be1f2b712ec6e2fc8a5d6",
		},
		{
			input:  "66cccceae528d4934df6eec6af52cf38519985c18558c343db467ff33cb47d88",
			output: "a8d57de754f97c861408b10a64e8ac2c5cafa89c69945eb2598e5cc28db43489",
		},
		{
			input:  "65a2545d400f03ab2617cab753a6c0aa4568d464f510df0c6456f0513e99ff34",
			output: "d979afebe4db319f1eeec3044769eaddcff05e8774c6f7246d650891c524195e",
		},
		{
			input:  "c5169aebac063bbefe9645fe11da4740bdd2281a5b54ee3f7df9a5672765e3bc",
			output: "aa60ba2866fdef4aa41afb542894148db2194e8faf04d4448468d173af589b8c",
		},
		{
			input:  "d7b703ee568f1841b4a19ac1a53f5d0ae2123c2d638c2edd862ad9a5872132cd",
			output: "5f453a7dc20ca516cb495883007f9a01135af2bfb8407b4b3f0fcca97d6dc877",
		},
		{
			input:  "6a7281f25792175b2ec4f0d7f74d19a37846393a6641b85823a9947d09acdbf4",
			output: "a4497915020f20d08cec766112b070b91e4d7a652a131518f208e958f468bb42",
		},
		{
			input:  "f56d072b656d46dd8d4709e120e402e2a6d265f8cca0e3997873468cd42d4e75",
			output: "9784d3b073e146c528d5d40a33f15d48fba2f6a81101ff09449fd474e3ab257c",
		},
		{
			input:  "bd395bb3ac19602aab43dd957541c9f7097ddc6f5e1e69ae6e088743216b0871",
			output: "a86523541840d10d0f8c4c39bcb16f010b77e3bd993110bd409da5091c1e5505",
		},
		{
			input:  "c90d485af8868c780e10f21c7286446ce3e8d1a9a39ba20e8775893d44143380",
			output: "b3753382b7f532f959b8d0aa9ba9a1276ecad1f96ba5d669968c23e07af62a22",
		},
		{
			input:  "0c730fcde54b69da226291eae0205e729f08c6e1536a806fa439b7a7d6d74230",
			output: "ed13cc686258196b8197787806792c4bbbeacf9e0c7d40fedb1cbecbcd3368ee",
		},
		{
			input:  "604db8331f1cbbba4adddb6aeaf71f53e325ef1130227b27d693a5f1f7619e30",
			output: "44591d8674a1f161970ee9287250a815f9f3e6a737262e57da72493874de10cd",
		},
		{
			input:  "20a3b82a2ad57cd0fc12f4bb4cc27bf9731129391b2c7034965fb338089c3aae",
			output: "6e3d6981f276726bf01b4d7fa347f6ed54933745e7eba5f229e662a39e0225ef",
		},
		{
			input:  "2b5414a0d6fe1ada82f342a24448bc3ce52d533c268385940bded8bb3fc153c2",
			output: "53c5082a759d7aa129a125b08e44b30cb1065e0bdbebce839791979369831b3b",
		},
		{
			input:  "2538dad623dab29a3d5387804ab51dea411014fad9c47fb94b2d83e44358064b",
			output: "62c4ac19375787857e2e7c41282a58fb68638a25ba27402239edc8d2683b5126",
			pers:   "pers",
		},
		{
			input:  "28d5ce4e064ddc5ecf7a62a65d245facb3b46d6d2d70cdcf5c413e1c4e205157",
			output: "9f82d3a2d796ac89aa012c71ccad798f5e10cf515dd6fb0fedceb44271029d43",
			pers:   "pers",
		},
		{
			input:  "172aaa13d4b7d606134bc6989a84c403bac3b6c4f9f4504dd79b2618c11c63d3",
			output: "5bd719de97d5dbc0012bd576f778994ae9d3eeb31e12cfadeb58133a34c2ebef",
			pers:   "pers",
		},
		{
			input:  "c9d4866805b9e9788cf1c7f6a369e8932e5bc6e32a5c7a06c3ecb7e0ec8ebaff",
			output: "ac22b64773effd8d0dbb088bd71b27449ec5c8c21ecc0262f615e920ea2df2a6",
			pers:   "pers",
		},
		{
			input:  "fae4f16aed075d1e94b6d913242840ad3933bc675aa72f2f647f1df60c431963",
			output: "a625f58d4585b8bad21509f34f715f6a15a5e2741a89642e937587b5c4b68e8a",
			pers:   "pers",
		},
	}
	return testCaseSlice
}

func TestYespower(t *testing.T) {
	// NOTE: 'in' values copied directly from C implementation tests
	in := []byte{0x00, 0x03, 0x06, 0x09, 0x0c, 0x0f, 0x12, 0x15,
		0x18, 0x1b, 0x1e, 0x21, 0x24, 0x27, 0x2a, 0x2d,
		0x30, 0x33, 0x36, 0x39, 0x3c, 0x3f, 0x42, 0x45,
		0x48, 0x4b, 0x4e, 0x51, 0x54, 0x57, 0x5a, 0x5d,
		0x60, 0x63, 0x66, 0x69, 0x6c, 0x6f, 0x72, 0x75,
		0x78, 0x7b, 0x7e, 0x81, 0x84, 0x87, 0x8a, 0x8d,
		0x90, 0x93, 0x96, 0x99, 0x9c, 0x9f, 0xa2, 0xa5,
		0xa8, 0xab, 0xae, 0xb1, 0xb4, 0xb7, 0xba, 0xbd,
		0xc0, 0xc3, 0xc6, 0xc9, 0xcc, 0xcf, 0xd2, 0xd5,
		0xd8, 0xdb, 0xde, 0xe1, 0xe4, 0xe7, 0xea, 0xed}

	out := hex.EncodeToString(Hash(in, 2048, 32, ""))

	want := "d5efb813cd263e9b34540130233cbbc6a921fbff3431e5ec1a1abde2aea6ff4d"

	if out != want {
		t.Errorf("got %s want %s", out, want)
	}

	fmt.Println(hex.EncodeToString(in), out)

	fmt.Println("testing 30 cases")

	tests := getTestCases()

	for i, tt := range tests {
		in, err := hex.DecodeString(tt.input)

		N := tt.N
		r := tt.r

		if N == uint32(0) {
			N = 2048
		}

		if r == uint32(0) {
			r = 32
		}

		if err != nil {
			t.Errorf("test %d: error %x", i, err)
		}

		out := hex.EncodeToString(Hash(in, N, r, tt.pers))

		if out != tt.output {
			t.Errorf("test %d: got %s want %s", in, out, tt.output)
		}

		fmt.Println(tt.input, out)
	}
}

func TestYespowerNative(t *testing.T) {
	// NOTE: 'in' values copied directly from C implementation tests
	in := []byte{0x00, 0x03, 0x06, 0x09, 0x0c, 0x0f, 0x12, 0x15,
		0x18, 0x1b, 0x1e, 0x21, 0x24, 0x27, 0x2a, 0x2d,
		0x30, 0x33, 0x36, 0x39, 0x3c, 0x3f, 0x42, 0x45,
		0x48, 0x4b, 0x4e, 0x51, 0x54, 0x57, 0x5a, 0x5d,
		0x60, 0x63, 0x66, 0x69, 0x6c, 0x6f, 0x72, 0x75,
		0x78, 0x7b, 0x7e, 0x81, 0x84, 0x87, 0x8a, 0x8d,
		0x90, 0x93, 0x96, 0x99, 0x9c, 0x9f, 0xa2, 0xa5,
		0xa8, 0xab, 0xae, 0xb1, 0xb4, 0xb7, 0xba, 0xbd,
		0xc0, 0xc3, 0xc6, 0xc9, 0xcc, 0xcf, 0xd2, 0xd5,
		0xd8, 0xdb, 0xde, 0xe1, 0xe4, 0xe7, 0xea, 0xed}

	out := hex.EncodeToString(YespowerNative(in, 2048, 32, ""))

	want := "d5efb813cd263e9b34540130233cbbc6a921fbff3431e5ec1a1abde2aea6ff4d"

	if out != want {
		t.Errorf("got %s want %s", out, want)
	}

	fmt.Println(hex.EncodeToString(in), out)

	fmt.Println("testing 30 cases")

	tests := getTestCases()

	for i, tt := range tests {
		in, err := hex.DecodeString(tt.input)

		N := tt.N
		r := tt.r

		if N == uint32(0) {
			N = 2048
		}

		if r == uint32(0) {
			r = 32
		}

		// Skip native test when r is being used as output length
		if r != uint32(32) {
			continue
		}

		if err != nil {
			t.Errorf("test %d: error %x", i, err)
		}

		out := hex.EncodeToString(YespowerNative(in, int(N), int(r), tt.pers))

		if out != tt.output {
			t.Errorf("test %d: got %s want %s", in, out, tt.output)
		}

		fmt.Println(tt.input, out)
	}
}

// --- TESTS-OK coverage ---------------------------------------------------
//
// The cases below mirror ./yespower/TESTS-OK, which is produced by the
// reference C driver in ./yespower/tests.c. That driver hashes an 80-byte
// source buffer where src[i] = i*3, using the given version/N/r/pers.
// HashVersion takes version as 5 (YESPOWER_0_5) or 10 (YESPOWER_1_0),
// matching yespower_version_t.

// refSrc reproduces the C test source: src[i] = (i*3) & 0xff for 80 bytes.
func refSrc() []byte {
	src := make([]byte, 80)
	for i := range src {
		src[i] = byte(i * 3)
	}
	return src
}

// parseHexSpaced parses a space-separated hex string like the TESTS-OK digests.
func parseHexSpaced(t *testing.T, s string) []byte {
	t.Helper()
	clean := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			clean = append(clean, s[i])
		}
	}
	b, err := hex.DecodeString(string(clean))
	if err != nil {
		t.Fatalf("bad expected hex %q: %v", s, err)
	}
	return b
}

type testsOKCase struct {
	version uint32
	N, r    uint32
	pers    string
	// persIsSrc is true for the BSTY case, where pers is the 80-byte source.
	persIsSrc bool
	want      string
}

func (c testsOKCase) persToken(src []byte) string {
	if c.persIsSrc {
		return string(src)
	}
	return c.pers
}

// testsOKCases enumerates every print_yespower(...) line in tests.c in order.
var testsOKCases = []testsOKCase{
	{Version05, 2048, 8, "Client Key", false, "a5 9f ec 4c 4f dd a1 6e 3b 14 05 ad da 66 d5 25 b6 8e 7c ad fc fe 6a c0 66 c7 ad 11 8c d8 05 90"},
	{Version05, 2048, 8, "", true, "5e a2 b2 95 6a 9e ac e3 0a 32 37 ff 1d 44 1e de e1 dc 25 aa b8 f0 ea 15 c1 21 65 f8 3a 7b c2 65"},
	{Version05, 4096, 16, "Client Key", false, "92 7e 72 d0 de d3 d8 04 75 47 3f 40 f1 74 3c 67 28 9d 45 3d 52 42 d4 f5 5a f4 e3 25 e0 66 99 c5"},
	{Version05, 4096, 24, "Jagaricoin", false, "0e 13 66 97 32 11 e7 fe a8 ad 9d 81 98 9c 84 a2 54 d9 68 c9 d3 33 dd 8f f0 99 32 4f 38 61 1e 04"},
	{Version05, 4096, 32, "WaviBanana", false, "3a e0 5a bb 3c 5c f6 f7 54 15 a9 25 54 c9 8d 50 e3 8e c9 55 2c fa 78 37 36 16 f4 80 b2 4e 55 9f"},
	{Version05, 2048, 32, "Client Key", false, "56 0a 89 1b 5c a2 e1 c6 36 11 1a 9f f7 c8 94 a5 d0 a2 60 2f 43 fd cf a5 94 9b 95 e2 2f e4 46 1e"},
	{Version05, 1024, 32, "Client Key", false, "2a 79 e5 3d 1b e6 66 9b c5 56 cc c4 17 bc e3 d2 2a 74 a2 32 f5 6b 8e 1d 39 b4 57 92 67 5d e1 08"},
	{Version05, 2048, 8, "", false, "5e cb d8 e8 d7 c9 0b ae d4 bb f8 91 6a 12 25 dc c3 c6 5f 5c 91 65 ba e8 1c dd e3 cf fa d1 28 e8"},

	{Version10, 2048, 8, "", false, "69 e0 e8 95 b3 df 7a ee b8 37 d7 1f e1 99 e9 d3 4f 7e c4 6e cb ca 7a 2c 43 08 e5 18 57 ae 9b 46"},
	{Version10, 4096, 16, "", false, "33 fb 8f 06 38 24 a4 a0 20 f6 3d ca 53 5f 5c a6 6a b5 57 64 68 c7 5d 1c ca ac 75 42 f7 64 95 ac"},
	{Version10, 4096, 32, "", false, "77 1a ee fd a8 fe 79 a0 82 5b c7 f2 ae e1 62 ab 55 78 57 46 39 ff c6 ca 37 23 cc 18 e5 e3 e2 85"},
	{Version10, 2048, 32, "", false, "d5 ef b8 13 cd 26 3e 9b 34 54 01 30 23 3c bb c6 a9 21 fb ff 34 31 e5 ec 1a 1a bd e2 ae a6 ff 4d"},
	{Version10, 1024, 32, "", false, "50 1b 79 2d b4 2e 38 8f 6e 7d 45 3c 95 d0 3a 12 a3 60 16 a5 15 4a 68 83 90 dd c6 09 a4 0c 67 99"},
	{Version10, 1024, 32, "personality test", false, "1f 02 69 ac f5 65 c4 9a dc 0e f9 b8 f2 6a b3 80 8c dc 38 39 4a 25 4f dd ee dc c3 aa cf f6 ad 9d"},
}

// TestTestsOK verifies the per-parameter digests listed in yespower/TESTS-OK.
func TestTestsOK(t *testing.T) {
	src := refSrc()
	for _, c := range testsOKCases {
		got := HashVersion(src, c.N, c.r, c.persToken(src), c.version)
		want := parseHexSpaced(t, c.want)
		if hex.EncodeToString(got) != hex.EncodeToString(want) {
			t.Errorf("yespower(%d, %d, %d, %q): got %x want %x",
				c.version, c.N, c.r, c.pers, got, want)
		}
	}
}

// TestTestsOKXOR verifies the two "XOR of yespower(...)" lines in TESTS-OK,
// reproducing print_yespower_loop from tests.c: src[0]=43, src[i]=i*3, then
// XOR every digest over N in {1024,2048,4096} and r in [8,32].
func TestTestsOKXOR(t *testing.T) {
	src := refSrc()
	src[0] = 43

	cases := []struct {
		version uint32
		pers    string
		want    string
	}{
		{Version05, "Client Key", "ae f1 32 91 87 0f 55 70 47 f4 2e 9b ef a6 16 df e5 f1 96 77 e1 3f 8b a6 92 f7 c5 97 55 a0 f5 0e"},
		{Version10, "", "8d 13 c5 fb 07 30 96 75 d1 b8 48 92 77 ba 4b e4 40 33 be df ae 7a 60 43 8a 9b e2 1f 3a 7b 12 37"},
	}

	for _, c := range cases {
		var xor [32]byte
		for n := uint32(1024); n <= 4096; n <<= 1 {
			for r := uint32(8); r <= 32; r++ {
				d := HashVersion(src, n, r, c.pers, c.version)
				for i := 0; i < 32; i++ {
					xor[i] ^= d[i]
				}
			}
		}
		want := parseHexSpaced(t, c.want)
		if hex.EncodeToString(xor[:]) != hex.EncodeToString(want) {
			t.Errorf("XOR of yespower(%d, ...): got %x want %x", c.version, xor[:], want)
		}
	}
}

var result []byte

func bench(b *testing.B, N int, r int) {
	for i := 0; i < b.N; i++ {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, uint32(i))
		ignore := Hash(bs, uint32(N), uint32(r), "")
		result = ignore
	}
}

func benchNative(b *testing.B, N int, r int) {
	for i := 0; i < b.N; i++ {
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, uint32(i))
		ignore := YespowerNative(bs, N, r, "")
		result = ignore
	}
}

func BenchmarkYespower_1024(b *testing.B)       { bench(b, 1024, 8) }
func BenchmarkYespower_2048(b *testing.B)       { bench(b, 2048, 32) }
func BenchmarkYespowerNative_1024(b *testing.B) { benchNative(b, 1024, 8) }
func BenchmarkYespowerNative_2048(b *testing.B) { benchNative(b, 2048, 32) }
