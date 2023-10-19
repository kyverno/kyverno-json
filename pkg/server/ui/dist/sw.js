importScripts('dist/wasm_exec.js')
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.1.0/sw.js')

// Skip installed stage and jump to activating stage
addEventListener('install', (event) => {
    event.waitUntil(skipWaiting())
})

// Start controlling clients as soon as the SW is activated
addEventListener('activate', event => {
    event.waitUntil(clients.claim())
})

registerWasmHTTPListener('assets/main.wasm', { base: 'api' })
