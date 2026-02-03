export function loop (a: number[]) {
    for (let i = 20; i > a.length; i--) {
        console.log(i);
    }
}
export function loop2 (a: number) {
    let counter: number = 0; // Initialize a counter variable

while (counter < a) { // Loop while counter is less than 5
  console.log(counter); // Print the current counter value
  counter++; // Increment the counter (crucial to avoid an infinite loop)
}

console.log("Loop finished!");

}