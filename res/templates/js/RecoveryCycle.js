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
			.x(function(d) { return self.x(d.interval); })
			.y(function(d) { return self.y(d.change); });

		// Add the valueline path.
		svg.append("path")
			.data([this.data])
			.attr("class", "line")
			.attr("d", valueline);

		svg.append("path")
			.data([
				[{ "interval": 1, "change": 0 }, { "interval": 200, "change": 0 }]
			])
			.attr("class", "line")
			.attr("d", valueline);

		const circles = svg.selectAll("circle")
			.data(this.data)
			.enter()
			.append("circle");
		circles.attr("cx", function(d) { return self.x(d.interval); })
			.attr("cy", function(d) { return self.y(d.change); })
			.attr("r", 5)
			.style("fill", "black");
	}
}
