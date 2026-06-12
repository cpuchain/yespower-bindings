{
    "targets": [
        {
            "target_name": "yespower",
            "sources": [
                "yespower.cc",
                "../yespower-c/sha256.c",
                "../yespower-c/yespower.c",
                "../yespower-c/yespower-opt.c",
            ],
            "include_dirs": [
                "<!@(node -p \"require('node-addon-api').include\")"
            ],
            "cflags_cc": [
                "-std=c++17"
            ],
            "defines": ["NAPI_DISABLE_CPP_EXCEPTIONS"]
        }
    ]
}
