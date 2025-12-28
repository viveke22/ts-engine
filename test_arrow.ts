// Arrow function with typed parameters and return type
const add = (a: number, b: number) => a + b;
console.log("Add:", add(5, 3)); // Output: Add: 8

// Multi-line arrow function with explicit return
const square: any = (x: number): number => {
    return x * x;
};
console.log("Square:", square(4)); // Output: Square: 16

// Arrow function with no parameters
const greet: any = (): void => console.log("Hello Arrow!");
greet(); // Output: Hello Arrow!