class RecoveryCycle extends Chart {
	constructor(plots) {
		super([1, 200], [-50, 110], Chart.scaleType.LOG)
		this.data = plots.rc.data
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	updatePlots(plots) {
		this.data = plots.rc.data
		this.animateXYLineWithMean(this.data, "rc")
	}

	drawLines(svg) {
		this.createXYLineWithMean(this.data, "rc")
		this.animateXYLineWithMean(this.data, "rc")
		this.drawHorizontalLine(this.linesLayer, 0)
	}
}
