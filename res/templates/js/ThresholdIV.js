class ThresholdIV extends Chart {
	constructor(participant, norms) {
		super([-400, 50], [-100, 50])
		this.participant = participant.sections.IV.data
		this.norms = norms.IV.data
		this.xIndex = 1
		this.yIndex = 0
		this.ySDName = undefined
		this.yMeanName = 3
		this.xSDName = 1
		this.xMeanName = 0
	}

	get name() { return "Threshold I/V" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Current (% Threshold)" }

	updateParticipant(participant) {
		this.participant = participant.sections.IV.data
		this.animateXYLine(this.participant, "tiv")
	}

	updateNorms(norms) {
		this.norms = norms.IV.data
		this.animateNorms(this.norms, "tiv")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "tiv")
		this.createNorms(this.norms, "tiv", false)
		this.drawHorizontalLine(this.linesLayer, 0)
		this.drawVerticalLine(this.linesLayer, 0)
		this.animateXYLine(this.participant, "tiv")
		this.animateNorms(this.norms, "tiv", false)
	}
}
