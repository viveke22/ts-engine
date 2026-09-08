let nums: any = [1, 2];
let lenAfterPush: number = nums.push(3);
console.log(lenAfterPush);
console.log(nums.join("-"));

let doubled: any = nums.map((n: number) => n * 2);
console.log(doubled.join(","));

let last: any = nums.pop();
console.log(last);
console.log(nums.length);

let name: string = "ts-engine";
console.log(`hello ${name} ${1 + 2}`);
console.log(`plain`);
