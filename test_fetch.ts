console.log("Fetching example.com...");
// Standard TS fetch returns Promise<Response>, engine returns number.
// We use an intermediate 'any' variable to bridge the type gap without using 'as'.
let raw: any = await fetch("https://httpbin.org/status/500");
let status: number = raw;
console.log("Status Code:", status);

if (status === 200) {
    console.log("Success!");
} else {
    console.log("Failed or different status");
}

export { };
