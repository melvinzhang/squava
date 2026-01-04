importScripts('wasm_exec.js');

const go = new Go();
let wasmInstance;

WebAssembly.instantiateStreaming(fetch("squava.wasm"), go.importObject).then((result) => {
    wasmInstance = result.instance;
    go.run(wasmInstance);
    postMessage({ type: 'READY' });
});

onmessage = (e) => {
    const { type, payload } = e.data;
    if (type === 'NEW_GAME') {
        const hash = squavaNewGame();
        const board = squavaGetBoard();
        postMessage({ type: 'GAME_UPDATED', payload: { hash, board } });
    } else if (type === 'APPLY_MOVE') {
        const hash = squavaApplyMove(payload.idx);
        const board = squavaGetBoard();
        postMessage({ type: 'GAME_UPDATED', payload: { hash, board } });
    } else if (type === 'GET_AI_MOVE') {
        const move = squavaGetBestMove(payload.iterations || 10000);
        postMessage({ type: 'AI_MOVE_RESULT', payload: { move } });
    } else if (type === 'GET_BOARD') {
        const board = squavaGetBoard();
        postMessage({ type: 'BOARD_RESULT', payload: board });
    }
};
