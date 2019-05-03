class StimulusRelative extends Chart {
	constructor(plots) {
		super([0, 200], [0, 100])
		const stimFor50PercentMax = plots.sr.data[24].valueX // Could also be extracted from excitability variables
		this.data = plots.sr.data.map(function callback(d, i) {
			return {
				'x': d.valueX / stimFor50PercentMax * 100, // Normalize each element
				'y': (i + 1) * 2,
			}
		})
		this.xName = 'x'
		this.yName = 'y'
	}

	get name() { return "Relative Stimulus Response" }
	get xLabel() { return "Stimulus (% Mean Threshold)" }
	get yLabel() { return "Peak Response (% Max)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.data, "srel")
	}
}
