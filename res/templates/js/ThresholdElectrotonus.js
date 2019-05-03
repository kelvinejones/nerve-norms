class ThresholdElectrotonus extends Chart {
	constructor(plots) {
		super([0, 200], [-150, 100])

		this.hy40 = plots.te.data.h40
		this.de40 = plots.te.data.d40
		this.hy20 = plots.te.data.h20
		this.de20 = plots.te.data.d20
	}

	get name() { return "Threshold Electrotonus" }
	get xLabel() { return "Threshold Reduction (%)" }
	get yLabel() { return "Delay (ms)" }

	updatePlots(plots) {
		this.hy40 = plots.te.data.h40
		this.de40 = plots.te.data.d40
		this.hy20 = plots.te.data.h20
		this.de20 = plots.te.data.d20
		this.animateLines()
	}

	drawLines(svg) {
		this.createXYLineWithMean(this.hy40, "hy40")
		this.createXYLineWithMean(this.de40, "de40")
		this.createXYLineWithMean(this.hy20, "hy20")
		this.createXYLineWithMean(this.de20, "de20")

		this.drawHorizontalLine(this.linesLayer, 0)

		this.animateLines()
	}

	animateLines() {
		this.animateXYLineWithMean(this.hy40, "hy40")
		this.animateXYLineWithMean(this.de40, "de40")
		this.animateXYLineWithMean(this.hy20, "hy20")
		this.animateXYLineWithMean(this.de20, "de20")
	}
}
