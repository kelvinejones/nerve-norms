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

		// define the line
		const valueData = Chart.valueData(this.data)
		const meanData = this.data.map(function(d) { return { x: d.delay, y: d.mean } })

		const normRange = Chart.normativeRange(this.data)

		const xZeroLine = this.xZeroLine()
		const xyLine = this.xyLine()

		// Draw the confidence interval
		svg.append("path")
			.data([normRange])
			.attr("class", "confidenceinterval")
			.attr("d", xZeroLine)
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", xyLine);

		svg.append("path")
			.data([meanData])
			.attr("class", "meanline")
			.attr("d", xZeroLine)
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", xyLine);

		// Add a reference line for 0
		svg.append("path")
			.data([
				[{ x: 1, y: 0 }, { x: 200, y: 0 }]
			])
			.attr("class", "meanline")
			.attr("d", xZeroLine);

		// Add the valueline path.
		svg.append("path")
			.data([valueData])
			.attr("class", "line")
			.attr("d", xZeroLine)
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", xyLine);

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
