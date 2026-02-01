let someValue: unknown = "Hello, TypeScript!";

// Using 'as' keyword for type assertion
let strLength: number = (someValue as string).length;

console.log(strLength); // Output: 18
