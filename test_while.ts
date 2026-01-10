
console.log("Testing While Loop...");

let i = 0;
while (i < 5) {
    console.log("i = " + i);
    i++;
}

if (i == 5) {
    console.log("PASS: While Loop");
} else {
    console.log("FAIL: While Loop i=" + i);
}

// Countdown
let c = 3;
while (c > 0) {
    console.log("Countdown: " + c);
    c--;
}
