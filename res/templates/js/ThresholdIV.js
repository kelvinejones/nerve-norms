class ThresholdIV extends Chart {
	constructor(plots) {
		super([-100, 50], [-400, 50])
		this.data = plots.tiv.data
		this.xName = 'current'
	}

	get name() { return "Threshold I/V" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Current (% Threshold)" }

	updatePlots(plots) {
		this.data = plots.tiv.data
		this.animateXYLineWithMean(this.data, "tiv")
	}

	drawLines(svg) {
		this.createXYLineWithMean(this.data, "tiv")
		this.drawHorizontalLine(this.linesLayer, 0)
		this.drawVerticalLine(this.linesLayer, 0)
		this.animateXYLineWithMean(this.data, "tiv")
	}
}
