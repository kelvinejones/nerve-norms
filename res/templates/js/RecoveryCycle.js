class RecoveryCycle extends Chart {
	constructor(participant, norms) {
		super([1, 200], [-50, 110], Chart.scaleType.LOG)
		this.participant = participant.sections.RC.data
		this.norms = norms.sections.RC.data
		this.xName = 0
		this.yName = 1
		this.xSDName = undefined // No standard deviation in the x direction.
		this.ySDName = 'SD'
		this.yMeanName = 'mean'
		this.xMeanName = 'delay'
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	updateParticipant(participant) {
		this.participant = participant.sections.RC.data
		this.animateXYLine(this.participant, "rc")
	}

	updateNorms(norms) {
		this.norms = norms.sections.RC.data
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
