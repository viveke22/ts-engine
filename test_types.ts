let x: number = 10;
let y: string = "hello";
let z: boolean = true;

console.log(x);
console.log(y);

function add(a: number, b: number): number {
    return a + b;
}

console.log(add(5, 5));

let sub: any = function (a: number, b: number): number {
    return a - b;
};

console.log(sub(10, 2));

export { };
