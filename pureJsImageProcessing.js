import * as StackBlur from './node_modules/stackblur-canvas/dist/stackblur-es.min.js';

let transformedImage = document.getElementById('transformed-image-js');
let sourceImage = document.getElementById('source-img');
let pureJsCallBtnBlur = document.getElementById('pure-js-button')

pureJsCallBtnBlur.addEventListener("click", () => {
  let startTime = performance.now();
  StackBlur.image(sourceImage, transformedImage, 10, true);
  let endTime = performance.now();
  console.log(`%c pure js done. t: ${(endTime - startTime)}ms`, 'color: #dedede')
})