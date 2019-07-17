class ThresholdIV extends Chart {
	constructor(participant, norms) {
		super([-400, 50], [-100, 50])
		this.participant = participant.sections.IV.data
		this.norms = (norms == null) ? undefined : norms.IV.data
		this.xIndex = 1
		this.yIndex = 0
		this.ySDIndex = undefined
		this.yMeanIndex = 3
		this.xSDIndex = 1
		this.xMeanIndex = 0
	}

	get name() { return "Threshold I/V" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Current (% Threshold)" }

	updateParticipant(participant) {
		if (participant.sections.IV == null) {
			this.participant = undefined
		} else {
			this.participant = participant.sections.IV.data
		}
		this.animateXYLine(this.participant, "tiv")
	}

	updateNorms(norms) {
		this.norms = norms.IV.data
		this.animateNorms(this.norms, "tiv")
	}

	updateScore(scores) {
		if (scores != null && scores.IV != null) {
			this.score = scores.IV.Overall
		}
	}

	drawLines(svg) {
		const isNull = (this.norms == null)
		const norms = isNull ? this.participant : this.norms
		this.createXYLine(this.participant, "tiv")
		this.createNorms(norms, "tiv", !isNull)
		this.drawHorizontalLine(this.linesLayer, 0)
		this.drawVerticalLine(this.linesLayer, 0)
		this.animateXYLine(this.participant, "tiv")
		this.animateNorms(norms, "tiv", !isNull)
	}
}
