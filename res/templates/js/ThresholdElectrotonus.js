class ThresholdElectrotonus extends Chart {
	constructor(participant) {
		super([0, 200], [-150, 100])
		this.participant = participant.sections.TE
		this.norms = this.participant
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	updateParticipant(participant) {
		this.participant = participant.sections.TE
		this.animateParticipant()
	}

	updateNorms(norms) {
		this.norms = norms.TE
		this.animateUpdatedNorms()
	}

	drawLines(svg) {
		Object.keys(this.participant).forEach(key => {
			this.createXYLine(this.participant[key].data, key)
			this.createNorms(this.norms[key].data, key, false)
		})
		this.drawHorizontalLine(this.linesLayer, 0)
		this.animateParticipant()
		this.animateUpdatedNorms(false)
	}

	animateParticipant() {
		Object.keys(this.participant).forEach(key => {
			this.animateXYLine(this.participant[key].data, key)
		})
	}

	animateUpdatedNorms(useSD = true) {
		Object.keys(this.norms).forEach(key => {
			this.animateNorms(this.norms[key].data, key, useSD)
		})
	}
}
