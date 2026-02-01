
console.log("Fetching posts (Array check)...");

let response = fetch("https://jsonplaceholder.typicode.com/posts?_limit=3");
if (response.ok) {
    let posts: any = response.json();
    console.log("Got posts!");

    // Check if we can access array elements
    let firstPost = posts[0];
    console.log("First Post Title:", firstPost.title);

    let secondPost = posts[1];
    console.log("Second Post Title:", secondPost.title);
} else {
    console.log("Fetch failed", response.status);
}
