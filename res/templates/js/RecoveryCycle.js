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

		this.drawCI(svg, [Chart.normativeRange(this.data)])

		svg.append("path")
			.data([Chart.dataAsXY(this.data, 'delay', 'mean')])
			.attr("class", "meanline")
			.attr("d", this.xZeroLine())
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", this.xyLine());

		// Add a reference line for 0
		svg.append("path")
			.data([
				[{ x: 1, y: 0 }, { x: 200, y: 0 }]
			])
			.attr("class", "meanline")
			.attr("d", this.xZeroLine());

		// Add the valueline path.
		svg.append("path")
			.data([Chart.dataAsXY(this.data)])
			.attr("class", "line")
			.attr("d", this.xZeroLine())
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", this.xyLine());

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
