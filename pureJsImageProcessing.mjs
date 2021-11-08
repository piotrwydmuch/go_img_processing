import { image } from './stackblur.js';

let transformedImage = document.getElementById('transformed-image-js');
let sourceImage = document.getElementById('source-img');
let pureJsCallBtnBlur = document.getElementById('pure-js-button')

pureJsCallBtnBlur.addEventListener("click", () => {
  let startTime = performance.now();
  image(sourceImage, transformedImage, 10, true);
  let endTime = performance.now();
  console.log(`%c pure js done. t: ${(endTime - startTime)}ms`, 'color: #dedede')
})