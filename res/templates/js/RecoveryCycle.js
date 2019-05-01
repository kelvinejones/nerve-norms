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
		this.animateCI(svg, [Chart.normativeRange(this.data)])

		this.animateLine(svg, [Chart.dataAsXY(this.data, 'delay', 'mean')], "meanline")

		this.drawHorizontalLine(svg, 0, 1)

		// Add the valueline path.
		this.animateLine(svg, [Chart.dataAsXY(this.data)], "line")

		this.animateCircles(svg, this.data)
	}

}
