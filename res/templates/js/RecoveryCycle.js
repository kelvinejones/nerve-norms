class RecoveryCycle extends Chart {
	constructor(data) {
		super()
		this.data = data

		this.xscale = d3.scaleLog().range([0, this.width]);
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	drawLines(svg) {
		const self = this

		this.animateCI(svg, [Chart.normativeRange(this.data)])

		this.animateLine(svg, [Chart.dataAsXY(this.data, 'delay', 'mean')], "meanline")

		this.drawHorizontalLine(svg, 0, 1)

		// Add the valueline path.
		this.animateLine(svg, [Chart.dataAsXY(this.data)], "line")

		const circles = svg.selectAll("circle")
			.data(this.data)
			.enter()
			.append("circle");
		circles.attr("cx", d => self.xscale(d.delay))
			.attr("cy", self.yscale(0))
			.attr("r", d => d.wasImputed ? 3 : 5)
			.style("fill", d => d.wasImputed ? "red" : "black");
		circles
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("cy", d => self.yscale(d.value))
	}

}
