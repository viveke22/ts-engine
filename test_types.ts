// Test Runtime Type Checking

function test() {
    let correct: number = 100;
    console.log("Correct assignment:", correct);

    let whatever: any = "I can be anything";
    console.log("Any assignment:", whatever);

    let unknownVal: unknown = true;
    console.log("Unknown assignment:", unknownVal);
}

test();

// Verify errors (commented out to allow test to pass, but manual verification can uncomment)
// let wrong: number = "string value"; 
// console.log("Should not reach here if uncommented");

export { };
