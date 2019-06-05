class ChargeDuration extends Chart {
	constructor(participant) {
		super([0, 1], [0, 10])
		this.participant = participant.sections.CD.data
		this.norms = this.participant
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	updateParticipant(participant) {
		this.participant = participant.sections.CD.data
		this.animateXYLine(this.participant, "cd")
	}

	updateNorms(norms) {
		this.norms = norms.CD.data
		this.animateNorms(this.norms, "cd")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "cd")
		this.animateXYLine(this.participant, "cd")
		this.createNorms(this.norms, "cd", false)
		this.animateNorms(this.norms, "cd", false)
	}
}
