class RecoveryCycle extends Chart {
	constructor(participant, norms) {
		super([1, 200], [-50, 110], Chart.scaleType.LOG)
		this.participant = participant.rc.data
		this.norms = norms.rc.data
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	updateParticipant(participant) {
		this.participant = participant.rc.data
		this.animateXYLine(this.participant, "rc")
	}

	updateNorms(norms) {
		this.norms = norms.rc.data
		this.animateNorms(this.norms, "rc")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "rc")
		this.createNorms(this.norms, "rc")

		this.drawHorizontalLine(this.linesLayer, 0)

		this.animateXYLine(this.participant, "rc")
		this.animateNorms(this.norms, "rc")
	}
}
