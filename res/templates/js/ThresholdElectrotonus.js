class ThresholdElectrotonus extends Chart {
	constructor(plots) {
		super([0, 200], [-150, 100])
		this.participant = plots.te.data
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	updatePlots(plots) {
		this.participant = plots.te.data
		this.animateParticipant()
	}

	updateNorms(norms) {
		this.animateUpdatedNorms()
	}

	drawLines(svg) {
		Object.keys(this.participant).forEach(key => {
			this.createXYLine(this.participant[key], key)
			this.createNorms(this.participant[key], key)
		})
		this.drawHorizontalLine(this.linesLayer, 0)
		this.animateParticipant()
		this.animateUpdatedNorms()
	}

	animateParticipant() {
		Object.keys(this.participant).forEach(key => {
			this.animateXYLine(this.participant[key], key)
		})
	}

	animateUpdatedNorms() {
		Object.keys(this.participant).forEach(key => {
			this.animateNorms(this.participant[key], key)
		})
	}
}
