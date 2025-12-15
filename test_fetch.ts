console.log("Fetching example.com...");
// Standard TS fetch returns Promise<Response>. Engine now returns a Hash object with status, ok, statusText.
let res: any = await fetch("https://example.com/l");

// Access properties
let status: number = res.status;
let ok: boolean = res.ok;
let msg: string = res.statusText;

console.log("Status Code:", status);
console.log("OK:", ok);
console.log("Status Text:", msg);

if (status === 200) {
    console.log("Success!");
} else {
    console.log("Failed or different status");
}

export { };
