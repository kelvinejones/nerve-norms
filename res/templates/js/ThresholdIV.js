class ThresholdIV extends Chart {
	constructor(data) {
		super()
		this.data = data
		this.xscale = this.xscale.domain([-100, 50]);
		this.yscale = this.yscale.domain([-400, 50]);
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
