function add(x, y) {
    return x + y;
}

let sum = add(5, 10);
console.log(sum);

// Function as expression
let sub = function (a, b) {
    return a - b;
};

console.log(sub(20, 5));

function hello() {
    console.log("Hello from function");
}
hello();

export { };
