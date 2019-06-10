class ChargeDuration extends Chart {
	constructor(participant, norms) {
		super([0, 1], [0, 10])
		this.participant = participant.sections.CD.data
		this.norms = (norms == null) ? undefined : norms.CD.data
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	updateParticipant(participant) {
		if (participant.sections.CD == null) {
			this.participant = undefined
		} else {
			this.participant = participant.sections.CD.data
		}
		this.animateXYLine(this.participant, "cd")
	}

	updateNorms(norms) {
		if (norms.CD == null) {
			this.norms = undefined
		} else {
			this.norms = norms.CD.data
		}
		this.animateNorms(this.norms, "cd")
	}

	drawLines(svg) {
		const isNull = (this.norms == null)
		const norms = isNull ? this.participant : this.norms
		this.createXYLine(this.participant, "cd")
		this.animateXYLine(this.participant, "cd")
		this.createNorms(norms, "cd", !isNull)
		this.animateNorms(norms, "cd", !isNull)
	}
}
