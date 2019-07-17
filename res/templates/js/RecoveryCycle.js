class RecoveryCycle extends Chart {
	constructor(participant, norms) {
		super([1, 200], [-50, 110], Chart.scaleType.LOG)
		this.participant = participant.sections.RC.data
		this.norms = (norms == null) ? undefined : norms.RC.data
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Interstimulus Interval (ms)" }
	get yLabel() { return "Threshold Change (%)" }

	updateParticipant(participant) {
		if (participant.sections.RC == null) {
			this.participant = undefined
		} else {
			this.participant = participant.sections.RC.data
		}
		this.animateXYLine(this.participant, "rc")
	}

	updateNorms(norms) {
		if (norms.RC == null) {
			this.norms = undefined
		} else {
			this.norms = norms.RC.data
		}
		this.animateNorms(this.norms, "rc")
	}

	updateScore(scores) {
		if (scores != null && scores.RC != null) {
			this.score = scores.RC.Overall
		}
	}

	drawLines(svg) {
		const isNull = (this.norms == null)
		const norms = isNull ? this.participant : this.norms
		this.createXYLine(this.participant, "rc")
		this.createNorms(norms, "rc", !isNull)

		this.drawHorizontalLine(this.linesLayer, 0)

		this.animateXYLine(this.participant, "rc")
		this.animateNorms(norms, "rc", !isNull)
	}
}
