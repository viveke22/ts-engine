let obj: any = { x: 0, a: { b: 1 } };
obj.x = 10;
obj.a.b = 20;
console.log(obj.x);
console.log(obj.a.b);

let arr: any = [1, 2];
arr[0] = 9;
arr[2] = 3;
console.log(arr[0]);
console.log(arr[1]);
console.log(arr[2]);

arr[4] = 5;
console.log(arr[3]);
console.log(arr[4]);

let bag: any = {};
bag["k"] = 7;
console.log(bag["k"]);

const n: number = 1;
console.log(n);

let u;
console.log(u);
console.log(null);
console.log(undefined === undefined);
console.log(null == undefined);

let s: number = 0;
for (let i: number = 0; i < 10; i++) {
    if (i == 3) {
        break;
    }
    if (i == 1) {
        continue;
    }
    s = s + i;
}
console.log(s);

let w: number = 0;
while (w < 5) {
    w = w + 1;
    if (w == 2) {
        continue;
    }
    if (w == 4) {
        break;
    }
}
console.log(w);

let complex: any = { list: [1, 2], matrix: [[0, 0], [0, 0]] };
complex.list[1] = 99;
complex.matrix[0][1] = 42;
console.log(complex.list[1]);
console.log(complex.matrix[0][1]);

let loopSum: number = 0;
for (let i: number = 0; i < 3; i++) {
    let factor: number = 10;
    const term: number = i * factor;
    loopSum = loopSum + term;
}
console.log(loopSum);

function emptyRet() {
    return;
}
console.log(emptyRet());
