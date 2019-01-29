class RecoveryCycle extends Chart {
	constructor(data) {
		super()

		this.data = data
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	drawLines(svg) {
		const self = this

		// define the line
		const valueline = d3.line()
			.x(function(d) { return self.x(d.delay); })
			.y(function(d) { return self.y(d.value); });

		const meanline = d3.line()
			.x(function(d) { return self.x(d.delay); })
			.y(function(d) { return self.y(d.mean); });

		const intervalData = (Array.from(this.data)
				.map(function(d) { return { x: d.delay, y: d.mean + 3 * d.SD } }))
			.concat(Array.from(this.data).reverse().map(function(d) { return { x: d.delay, y: d.mean - 3 * d.SD } }))
		const interval = d3.line()
			.x(function(d) { console.log(d); return self.x(d.x); })
			.y(function(d) { return self.y(d.y); });

		// Draw the confidence interval
		svg.append("path")
			.data([intervalData])
			.attr("class", "confidenceinterval")
			.attr("d", interval);

		svg.append("path")
			.data([this.data])
			.attr("class", "meanline")
			.attr("d", meanline);

		// Add a reference line for 0
		svg.append("path")
			.data([
				[{ "delay": 1, "value": 0 }, { "delay": 200, "value": 0 }]
			])
			.attr("class", "meanline")
			.attr("d", valueline);

		// Add the valueline path.
		svg.append("path")
			.data([this.data])
			.attr("class", "line")
			.attr("d", valueline);

		const circles = svg.selectAll("circle")
			.data(this.data)
			.enter()
			.append("circle");
		circles.attr("cx", function(d) { return self.x(d.delay); })
			.attr("cy", function(d) { return self.y(d.value); })
			.attr("r", 5)
			.style("fill", "black");
	}
}
