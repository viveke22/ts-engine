import { add, sub, mul, div, mod, float, pi } from "./test_export_function";

let result = add(5, 3);
console.log("Result:", result);
let result2 = sub(5, 3);
console.log("Result:", result2);
let result3 = mul(5, 3);
console.log("Result:", result3);
let result4 = div(5, 3);
console.log("Result:", result4);
let result5 = mod(5, 3);
console.log("Result:", result5);
let result6 = float(5.0, 3.5);
console.log("Result:", result6);
console.log("Result:", pi * pi);
