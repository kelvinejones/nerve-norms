class ChargeDuration extends Chart {
	constructor(plots) {
		super([0, 1], [0, 10])
		this.participant = plots.cd.data
		this.xName = 'stimWidth'
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	updatePlots(plots) {
		this.participant = plots.cd.data
		this.animateXYLine(this.participant, "cd")
		this.updateNorms()
	}

	updateNorms(norms) {
		this.animateNorms(this.participant, "cd")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "cd")
		this.animateXYLine(this.participant, "cd")
		this.createNorms(this.participant, "cd")
		this.animateNorms(this.participant, "cd")
	}
}
