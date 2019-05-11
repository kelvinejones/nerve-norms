class ThresholdIV extends Chart {
	constructor(plots) {
		super([-100, 50], [-400, 50])
		this.participant = plots.tiv.data
		this.xName = 'current'
	}

	get name() { return "Threshold I/V" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Current (% Threshold)" }

	updatePlots(plots) {
		this.participant = plots.tiv.data
		this.animateXYLine(this.participant, "tiv")
	}

	updateNorms(norms) {
		this.animateNorms(this.participant, "tiv")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "tiv")
		this.createNorms(this.participant, "tiv")
		this.drawHorizontalLine(this.linesLayer, 0)
		this.drawVerticalLine(this.linesLayer, 0)
		this.animateXYLine(this.participant, "tiv")
		this.animateNorms(this.participant, "tiv")
	}
}
