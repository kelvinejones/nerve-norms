class RecoveryCycle extends Chart {
	constructor(participant, norms) {
		super([1, 200], [-50, 110], Chart.scaleType.LOG)
		this.participant = participant.sections.RC.data
		this.norms = (norms === undefined) ? undefined : norms.RC.data
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	updateParticipant(participant) {
		if (participant.sections.RC === undefined) {
			this.participant = undefined
		} else {
			this.participant = participant.sections.RC.data
		}
		this.animateXYLine(this.participant, "rc")
	}

	updateNorms(norms) {
		if (norms.RC === undefined) {
			this.norms = undefined
		} else {
			this.norms = norms.RC.data
		}
		this.animateNorms(this.norms, "rc")
	}

	drawLines(svg) {
		const useSD = (this.norms !== undefined)
		const norms = (this.norms === undefined) ? this.participant : this.norms
		this.createXYLine(this.participant, "rc")
		this.createNorms(norms, "rc", useSD)

		this.drawHorizontalLine(this.linesLayer, 0)

		this.animateXYLine(this.participant, "rc")
		this.animateNorms(norms, "rc", useSD)
	}
}
