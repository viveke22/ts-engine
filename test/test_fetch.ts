// test_fetch.ts

console.log("Fetching a programming joke...");

let response: any = fetch("https://v2.jokeapi.dev/joke/Programming?type=single");

console.log("Status:", response.status);
console.log("Status Text:", response.statusText);

if (response.ok) {
    let data: any = response.json();
    console.log("Full Data:", data);
    console.log("--------------------------------------------------");
    console.log("Joke Category:", data.category);
    console.log("Joke:", data.joke);
    console.log("--------------------------------------------------");
} else {
    console.log("Fetch failed!");
}
