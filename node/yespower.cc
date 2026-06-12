#include <stdint.h>
#include <vector>
#include <string>

#include <napi.h>

extern "C" {
    #include "../yespower-c/yespower.h"
}

// Parsed + validated parameters, with an owned copy of the input and
// personality so that asynchronous workers don't depend on the lifetime of the
// JS Buffers / strings.
struct YespowerArgs {
    std::vector<uint8_t> input;
    std::vector<uint8_t> pers;
    uint32_t N;
    uint32_t r;
    uint32_t version;
    bool ok;
};

static YespowerArgs ParseArgs(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();
    YespowerArgs args;
    args.ok = false;

    // Required input buffer.
    if (info.Length() < 1 || !info[0].IsBuffer()) {
        Napi::TypeError::New(env, "First argument must be a Buffer").ThrowAsJavaScriptException();
        return args;
    }
    auto inputBuf = info[0].As<Napi::Buffer<uint8_t>>();
    args.input.assign(inputBuf.Data(), inputBuf.Data() + inputBuf.Length());

    // Optional N (default 2048).
    args.N = 2048;
    if (info.Length() > 1 && info[1].IsNumber()) {
        args.N = info[1].As<Napi::Number>().Uint32Value();
    }

    // Optional r (default 32).
    args.r = 32;
    if (info.Length() > 2 && info[2].IsNumber()) {
        args.r = info[2].As<Napi::Number>().Uint32Value();
    }

    // Optional pers: string (UTF-8) or Buffer / TypedArray (binary-safe).
    // Copied so async workers own the bytes.
    if (info.Length() > 3 && !info[3].IsUndefined() && !info[3].IsNull()) {
        if (info[3].IsString()) {
            std::string s = info[3].As<Napi::String>().Utf8Value();
            args.pers.assign(s.begin(), s.end());
        } else if (info[3].IsBuffer()) {
            auto b = info[3].As<Napi::Buffer<uint8_t>>();
            args.pers.assign(b.Data(), b.Data() + b.Length());
        } else if (info[3].IsTypedArray()) {
            auto ta = info[3].As<Napi::TypedArray>();
            auto data = reinterpret_cast<const uint8_t*>(ta.ArrayBuffer().Data()) + ta.ByteOffset();
            args.pers.assign(data, data + ta.ByteLength());
        } else {
            Napi::TypeError::New(env, "pers must be a string, Buffer or Uint8Array")
                .ThrowAsJavaScriptException();
            return args;
        }
    }

    // Optional version (default 10 = YESPOWER_1_0).
    args.version = 10;
    if (info.Length() > 4 && info[4].IsNumber()) {
        args.version = info[4].As<Napi::Number>().Uint32Value();
    }
    if (args.version != 5 && args.version != 10) {
        Napi::RangeError::New(env, "yespower: version must be 5 (YESPOWER_0_5) or 10 (YESPOWER_1_0)")
            .ThrowAsJavaScriptException();
        return args;
    }

    args.ok = true;
    return args;
}

// Runs the hash into a 32-byte output. Returns true on success.
static bool RunHash(const YespowerArgs& args, uint8_t out[32]) {
    const char* input = reinterpret_cast<const char*>(args.input.data());
    const char* pers = args.pers.empty() ? nullptr : reinterpret_cast<const char*>(args.pers.data());
    int rc = yespower_hash(input, (uint32_t)args.input.size(), args.N, args.r,
        pers, (uint32_t)args.pers.size(), args.version, reinterpret_cast<char*>(out));
    return rc == 0;
}

static const char* kFailureMessage =
    "yespower failed: invalid parameters (N, r or version) or out of memory";

// ----------------------------- synchronous ----------------------------- //

Napi::Value YespowerFunc(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();
    YespowerArgs args = ParseArgs(info);
    if (!args.ok || env.IsExceptionPending()) {
        return env.Null();
    }

    uint8_t out[32];
    if (!RunHash(args, out)) {
        Napi::Error::New(env, kFailureMessage).ThrowAsJavaScriptException();
        return env.Null();
    }

    return Napi::Buffer<uint8_t>::Copy(env, out, 32);
}

// ----------------------------- asynchronous ---------------------------- //
// Runs the (CPU-heavy) hash on the libuv threadpool so the event loop is not
// blocked. The worker owns copies of all inputs (YespowerArgs holds vectors).

class YespowerWorker : public Napi::AsyncWorker {
public:
    YespowerWorker(Napi::Env env, YespowerArgs&& args)
        : Napi::AsyncWorker(env),
          deferred_(Napi::Promise::Deferred::New(env)),
          args_(std::move(args)),
          success_(false) {}

    Napi::Promise GetPromise() { return deferred_.Promise(); }

protected:
    void Execute() override {
        success_ = RunHash(args_, out_);
        if (!success_) {
            SetError(kFailureMessage);
        }
    }

    void OnOK() override {
        Napi::Env env = Env();
        Napi::HandleScope scope(env);
        deferred_.Resolve(Napi::Buffer<uint8_t>::Copy(env, out_, 32));
    }

    void OnError(const Napi::Error& e) override {
        deferred_.Reject(e.Value());
    }

private:
    Napi::Promise::Deferred deferred_;
    YespowerArgs args_;
    bool success_;
    uint8_t out_[32];
};

Napi::Value YespowerAsyncFunc(const Napi::CallbackInfo& info) {
    Napi::Env env = info.Env();
    YespowerArgs args = ParseArgs(info);
    if (!args.ok || env.IsExceptionPending()) {
        // Convert the (pending) validation error into a rejected Promise so
        // callers can always rely on a Promise being returned. Clear the pending
        // exception before creating the Deferred to avoid operating on the env
        // in a faulted state under NAPI_DISABLE_CPP_EXCEPTIONS.
        Napi::Value reason;
        if (env.IsExceptionPending()) {
            reason = env.GetAndClearPendingException().Value();
        } else {
            reason = Napi::Error::New(env, "yespower: invalid arguments").Value();
        }
        Napi::Promise::Deferred deferred = Napi::Promise::Deferred::New(env);
        deferred.Reject(reason);
        return deferred.Promise();
    }

    YespowerWorker* worker = new YespowerWorker(env, std::move(args));
    Napi::Promise promise = worker->GetPromise();
    worker->Queue();
    return promise;
}

Napi::Object Init(Napi::Env env, Napi::Object exports) {
    exports.Set("yespower", Napi::Function::New(env, YespowerFunc));
    // Asynchronous variant: runs on the libuv threadpool and returns a Promise.
    exports.Set("yespower_async", Napi::Function::New(env, YespowerAsyncFunc));
    return exports;
}

NODE_API_MODULE(yespower, Init)
