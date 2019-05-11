class RecoveryCycle extends Chart {
	constructor(plots) {
		super([1, 200], [-50, 110], Chart.scaleType.LOG)
		this.participant = plots.rc.data
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	updatePlots(plots) {
		this.participant = plots.rc.data
		this.animateXYLine(this.participant, "rc")
		this.updateNorms()
	}

	updateNorms(norms) {
		this.animateNorms(this.participant, "rc")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "rc")
		this.createNorms(this.participant, "rc")

		this.drawHorizontalLine(this.linesLayer, 0)

		this.animateXYLine(this.participant, "rc")
		this.animateNorms(this.participant, "rc")
	}
}
