class ChargeDuration extends Chart {
	constructor(participant, norms) {
		super([0, 1], [0, 10])
		this.participant = participant.sections.CD.data
		this.norms = norms.sections.CD.data
		this.xName = 0
		this.yName = 1
		this.yMeanName = 'mean'
		this.xMeanName = 'stimWidth'
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	updateParticipant(participant) {
		this.participant = participant.sections.CD.data
		this.animateXYLine(this.participant, "cd")
	}

	updateNorms(norms) {
		this.norms = norms.sections.CD.data
		this.animateNorms(this.norms, "cd")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "cd")
		this.animateXYLine(this.participant, "cd")
		this.createNorms(this.norms, "cd")
		this.animateNorms(this.norms, "cd")
	}
}
