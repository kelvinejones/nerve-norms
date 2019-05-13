class ChargeDuration extends Chart {
	constructor(participant, norms) {
		super([0, 1], [0, 10])
		this.participant = participant.cd.data
		this.norms = norms.cd.data
		this.xName = 'stimWidth'
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	updateParticipant(participant) {
		this.participant = participant.cd.data
		this.animateXYLine(this.participant, "cd")
	}

	updateNorms(norms) {
		this.norms = norms.cd.data
		this.animateNorms(this.norms, "cd")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "cd")
		this.animateXYLine(this.participant, "cd")
		this.createNorms(this.norms, "cd")
		this.animateNorms(this.norms, "cd")
	}
}
