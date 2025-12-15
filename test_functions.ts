function add(x: number, y: string) {
    return x + y;
}

let sum = add(5, "h");
console.log(sum);

// Function as expression
let sub = function (a: number, b: number) {
    return a - b;
};

console.log(sub(20, 5));

function hello() {
    console.log("Hello from function");
}
hello();

export { };
