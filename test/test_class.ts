class Point {
	x: number;
	y: number;
	constructor(x: number, y: number) {
		this.x = x;
		this.y = y;
	}
	sum() {
		return this.x + this.y;
	}
}

const p = new Point(2, 3);
console.log(p.sum());
