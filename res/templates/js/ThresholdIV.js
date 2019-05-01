class ThresholdIV extends Chart {
	constructor(data) {
		super([-100, 50], [-400, 50])
		this.data = data
		this.xName = 'current'
	}

	get name() { return "Threshold I/V" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Current (% Threshold)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.data)
		this.drawVerticalLine(svg, 0)
	}
}
