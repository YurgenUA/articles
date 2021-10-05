import assert from 'assert';

function Sum(num1: number, num2: number): number { 
  assert(typeof num1 === "number");
  assert(typeof num2 === "number");
  return num1 + num2;
}

// @ts-ignore
console.log('true + nothing will be:', Sum(true, 'nothing'));