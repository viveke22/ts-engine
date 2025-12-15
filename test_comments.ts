// This is a single line comment
console.log("Start");

/* This is a 
   multi-line comment */
console.log("Middle");

let a = 10; // Comment after code
console.log(a);

/* Inline comment */ let b = 20;
console.log(b);

let c = 5 + /* 5 + */ 5;
console.log(c); // Should be 10

// Comment with / inside: /path/to/something
// Comment with * inside: 5 * 5

/* 
   Comment with * inside 
   and / inside
*/

console.log("End");
export { };
