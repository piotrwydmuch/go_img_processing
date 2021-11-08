import * as Magick from 'https://knicknic.github.io/wasm-imagemagick/magickApi.js';

let magicCallBtnBlur = document.getElementById('btn-cpp-1')
let magicCallBtnEdges = document.getElementById('btn-cpp-2')
let magicCallMachineTesting = document.getElementById('btn-cpp-2-machine-testing')
let transformedImage = document.getElementById('transformed-image');

let resultArray = [];

let DoMagickCall = async function (command, index = 1 ) {
    // performance.clearMarks(); // it have to be empty here !!
    let startTime = performance.now();
    let fetchedSourceImage = await fetch("lenna.png");
    let arrayBuffer = await fetchedSourceImage.arrayBuffer();
    let sourceBytes = new Uint8Array(arrayBuffer);

    const files = [{ 'name': 'srcFile.png', 'content': sourceBytes }];
    let processedFiles = await Magick.Call(files, command);

    let firstOutputImage = processedFiles[0]
    transformedImage.src = URL.createObjectURL(firstOutputImage['blob'])
    let endTime = performance.now();
    let performanceResult = endTime - startTime;
    resultArray.push(performanceResult);
    console.log(`%c ${index} imageMagic done. t: ${(performanceResult)}ms`, 'color: #008000');
};

magicCallBtnBlur.addEventListener("click", () => {
    const command = ["convert", "srcFile.png", "-channel", "RGBA", "-blur", "0x8", "out.png"];
    DoMagickCall(command);
})

magicCallBtnEdges.addEventListener("click", () => {
    const command = ["convert", "srcFile.png", "-canny", "0x1+10%+20%", "out.png"];
    DoMagickCall(command);
})

magicCallMachineTesting.addEventListener("click", () => {
    const command = ["convert", "srcFile.png", "-canny", "0x1+10%+20%", "out.png"];
    async function machineTesting() {
        for (let i = 1; i < 11; i++) {
            DoMagickCall(command, i);
        }
        const result = await resolveAfter2Seconds();
        console.log(result);
    }
    machineTesting();
})

function resolveAfter2Seconds() {
    return new Promise(() => {
      console.log("koniec")
    });
  }

  
