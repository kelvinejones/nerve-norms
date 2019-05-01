class RecoveryCycle extends Chart {
	constructor(data) {
		super()
		this.data = data
		this.xscale = d3.scaleLog().range([0, this.width]).domain([1, 200]);
		this.yscale = this.yscale.domain([-50, 110]);
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	drawLines(svg) {
		this.animateCI(svg, [Chart.normativeRange(this.data)])
		this.animateLine(svg, [Chart.dataAsXY(this.data, 'delay', 'mean')], "meanline")
		this.drawHorizontalLine(svg, 0, 1)
		this.animateLine(svg, [Chart.dataAsXY(this.data)], "line")
		this.animateCircles(svg, this.data)
	}
}
