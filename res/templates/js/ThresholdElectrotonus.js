class ThresholdElectrotonus extends Chart {
	constructor(hy40, de40, hy20, de20) {
		super()

		this.hy40 = hy40
		this.de40 = de40
		this.hy20 = hy20
		this.de20 = de20

		this.xscale = this.xscale.domain([1, 200]);
		this.yscale = this.yscale.domain([-150, 100]);
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "??Threshold Change (%)" }
	get yLabel() { return "??Interstimulus Interval (ms)" }

	drawLines(svg) {
		this.animateCI(svg, [Chart.normativeRange(this.hy40)])
		this.animateCI(svg, [Chart.normativeRange(this.de40)])
		this.animateCI(svg, [Chart.normativeRange(this.hy20)])
		this.animateCI(svg, [Chart.normativeRange(this.de20)])

		this.animateLine(svg, [Chart.dataAsXY(this.hy40, 'delay', 'mean')], "meanline")
		this.animateLine(svg, [Chart.dataAsXY(this.de40, 'delay', 'mean')], "meanline")
		this.animateLine(svg, [Chart.dataAsXY(this.hy20, 'delay', 'mean')], "meanline")
		this.animateLine(svg, [Chart.dataAsXY(this.de20, 'delay', 'mean')], "meanline")

		this.drawHorizontalLine(svg, 0)

		this.animateLine(svg, [Chart.dataAsXY(this.hy40)], "line")
		this.animateLine(svg, [Chart.dataAsXY(this.de40)], "line")
		this.animateLine(svg, [Chart.dataAsXY(this.hy20)], "line")
		this.animateLine(svg, [Chart.dataAsXY(this.de20)], "line")

		this.animateCircles(svg, this.hy40)
		this.animateCircles(svg, this.de40)
		this.animateCircles(svg, this.hy20)
		this.animateCircles(svg, this.de20)
	}
}
