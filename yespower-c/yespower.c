#include <stdlib.h>

#include "yespower.h"

/*
 * Thin wrappers around the reference yespower_tls()/yespower() API used by the
 * Node.js (NAPI) and WebAssembly (Emscripten) bindings.
 *
 * `version` is the raw yespower_version_t value: 5 (YESPOWER_0_5) or
 * 10 (YESPOWER_1_0).
 *
 * `pers` is treated as a length-delimited byte buffer (persLen bytes); it may
 * contain embedded NUL bytes and may be NULL when persLen == 0.
 *
 * A thread-local yespower_local_t is initialised once and reused across calls
 * (mirroring the reference yespower() usage) to avoid per-call setup overhead.
 * Under Emscripten's default single-threaded build __thread degrades to a plain
 * static, which is fine.
 *
 * Returns 0 on success, or -1 on error (propagated from yespower()).
 */

/* Portable thread-local storage: C11 _Thread_local, MSVC __declspec(thread),
 * GCC/Clang __thread. Under Emscripten's default single-threaded build this is
 * effectively a plain static, which is fine. */
#if defined(_MSC_VER)
#  define YP_THREAD_LOCAL __declspec(thread)
#elif defined(__STDC_VERSION__) && __STDC_VERSION__ >= 201112L && !defined(__STDC_NO_THREADS__)
#  define YP_THREAD_LOCAL _Thread_local
#elif defined(__GNUC__) || defined(__clang__)
#  define YP_THREAD_LOCAL __thread
#else
#  define YP_THREAD_LOCAL
#endif

static YP_THREAD_LOCAL yespower_local_t yp_local;
static YP_THREAD_LOCAL int yp_local_ready = 0;

int yespower_hash(const char* input, uint32_t inputLen, uint32_t N, uint32_t r,
    const char* pers, uint32_t persLen, uint32_t version, char* output) {
    if (!yp_local_ready) {
        if (yespower_init_local(&yp_local))
            return -1;
        yp_local_ready = 1;
    }

    const yespower_params_t params = {
        .version = (yespower_version_t)version,
        .N = N,
        .r = r,
        .pers = persLen ? (const uint8_t*)pers : NULL,
        .perslen = persLen
    };

    return yespower(&yp_local, (const uint8_t*)input, inputLen, &params,
        (yespower_binary_t*)output);
}

const char* yespower_wasm(const char* input, uint32_t inputLen, uint32_t N,
    uint32_t r, const char* pers, uint32_t persLen, uint32_t version) {
    char* output = malloc(32);
    if (!output)
        return NULL;

    if (yespower_hash(input, inputLen, N, r, pers, persLen, version, output)) {
        /* Zero the buffer on failure so callers don't read stale heap. */
        for (int i = 0; i < 32; i++)
            output[i] = 0;
    }

    return output;
}
