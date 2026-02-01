
// test_fs.ts - Testing File System CRUD

import * as fs from 'fs';

let filename = "test_file.txt";
let content = "Hello, File System!";

console.log("1. Writing file...");
fs.writeFileSync(filename, content);
console.log("File written.");

console.log("2. Reading file...");
let readContent = fs.readFileSync(filename);
console.log("Read content: " + readContent);

if (readContent === content) {
    console.log("PASS: Read/Write match.");
} else {
    console.log("FAIL: Content mismatch.");
}

console.log("3. Updating file (Overwriting)...");
let newContent = "Updated Content";
fs.writeFileSync(filename, newContent);
let updated = fs.readFileSync(filename);
console.log("Updated content: " + updated);

// console.log("4. Deleting file...");
// fs.removeSync(filename); // or unlinkSync
// if (fs.existsSync(filename)) {
//     console.log("FAIL: File still exists.");
// } else {
//     console.log("PASS: File deleted.");
// }
