import * as Magick from 'https://knicknic.github.io/wasm-imagemagick/magickApi.js';

let magicCallBtnBlur = document.getElementById('btn-cpp-1')
let magicCallBtnEdges = document.getElementById('btn-cpp-2')
let transformedImage = document.getElementById('transformed-image');

let DoMagickCall = async function (command) {
    let startTime = performance.now();
    let fetchedSourceImage = await fetch("lenna.png");
    let arrayBuffer = await fetchedSourceImage.arrayBuffer();
    let sourceBytes = new Uint8Array(arrayBuffer);

    const files = [{ 'name': 'srcFile.png', 'content': sourceBytes }];
    let processedFiles = await Magick.Call(files, command);

    let firstOutputImage = processedFiles[0]
    transformedImage.src = URL.createObjectURL(firstOutputImage['blob'])
    let endTime = performance.now();
    console.log(`%c imageMagic done. t: ${(endTime - startTime)}ms`, 'color: #008000')
};

magicCallBtnBlur.addEventListener("click", () => {
    const command = ["convert", "srcFile.png", "-channel", "RGBA", "-blur", "0x8", "out.png"];
    DoMagickCall(command);
})

magicCallBtnEdges.addEventListener("click", () => {
    const command = ["convert", "srcFile.png", "-canny", "0x1+10%+20%", "out.png"];
    DoMagickCall(command);
})
