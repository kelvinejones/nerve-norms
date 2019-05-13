class ThresholdIV extends Chart {
	constructor(participant, norms) {
		super([-100, 50], [-400, 50])
		this.participant = participant.tiv.data
		this.norms = norms.tiv.data
		this.xName = 'current'
	}

	get name() { return "Threshold I/V" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Current (% Threshold)" }

	updateParticipant(participant) {
		this.participant = participant.tiv.data
		this.animateXYLine(this.participant, "tiv")
	}

	updateNorms(norms) {
		this.norms = norms.tiv.data
		this.animateNorms(this.norms, "tiv")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "tiv")
		this.createNorms(this.norms, "tiv")
		this.drawHorizontalLine(this.linesLayer, 0)
		this.drawVerticalLine(this.linesLayer, 0)
		this.animateXYLine(this.participant, "tiv")
		this.animateNorms(this.norms, "tiv")
	}
}
