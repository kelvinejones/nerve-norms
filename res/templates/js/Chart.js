class Chart {
	constructor() {
		this.margin = { top: 20, right: 20, bottom: 50, left: 50 };
		this.width = 600 - this.margin.left - this.margin.right;
		this.height = 300 - this.margin.top - this.margin.bottom;

		// Set the default scaling
		this.xscale = d3.scaleLinear().range([0, this.width]);
		this.yscale = d3.scaleLinear().range([this.height, 0]);
	}

	get name() { throw new Error("A Chart must implement name()") }
	get xLabel() { throw new Error("A Chart must implement xLabel()") }
	get yLabel() { throw new Error("A Chart must implement yLabel()") }
	drawLines(svg) { throw new Error("A Chart must implement drawLines(svg)") }

	draw(svg, hideLabels) {
		// append the svg object to the body of the page
		// appends a 'group' element to 'svg'
		// moves the 'group' element to the top left margin
		svg = svg
			.append("g")
			.attr("transform",
				"translate(" + this.margin.left + "," + this.margin.top + ")");

		// Add layers for various elements
		this.ciLayer = svg.append("g")
		this.meanLayer = svg.append("g")
		this.linesLayer = svg.append("g")
		this.valueLayer = svg.append("g")
		this.circlesLayer = svg.append("g")

		// Add the X Axis
		var xelements = svg.append("g")
			.attr("transform", "translate(0," + this.height + ")")
			.call(d3.axisBottom(this.xscale).ticks(2)
				.tickFormat(d3.format("")));

		// Add the Y Axis
		var yelements = svg.append("g")
			.call(d3.axisLeft(this.yscale));

		if (!hideLabels) {
			this.labels(svg);
		} else {
			xelements.selectAll("text").remove();
			yelements.selectAll("text").remove();
		}

		this.drawLines(svg)
	}

	labels(svg) {
		// text label for the x axis
		svg.append("text")
			.attr("transform",
				"translate(" + (this.width / 2) + " ," +
				(this.height + this.margin.top + 20) + ")")
			.style("text-anchor", "middle")
			.text(this.xLabel);

		// text label for the y axis
		svg.append("text")
			.attr("transform", "rotate(-90)")
			.attr("y", 0 - this.margin.left)
			.attr("x", 0 - (this.height / 2))
			.attr("dy", "1em")
			.style("text-anchor", "middle")
			.text(this.yLabel);
	}

	static normativeRange(data, xName = 'delay', yName = 'mean', ySDName = 'SD') {
		return (Array.from(data)
				.map(function(d) { return { x: d[xName], y: d[yName] + 2 * (d[ySDName] || 0) } }))
			.concat(Array.from(data).reverse().map(function(d) { return { x: d[xName], y: d[yName] - 2 * (d[ySDName] || 0) } }))
	}

	static dataAsXY(data, xName = 'delay', yName = 'value') {
		return data.map(function(d) { return { x: d['delay'], y: d['value'] } })
	}

	xZeroLine() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(0))
	}

	xyLine() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(d.y));
	}

	animateCI(svg, ciNormRange) {
		svg.append("path")
			.data(ciNormRange)
			.attr("class", "confidenceinterval")
			.attr("d", this.xZeroLine())
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", this.xyLine());
	}

	animateLine(svg, xyLine, className) {
		svg.append("path")
			.data(xyLine)
			.attr("class", className)
			.attr("d", this.xZeroLine())
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("d", this.xyLine());
	}

	drawHorizontalLine(svg, yVal, xMin = 0, xMax = 1000000) {
		svg.append("path")
			.data([
				[{ x: xMin, y: yVal }, { x: xMax, y: yVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xZeroLine());
	}


	animateCircles(svg, circleLocations) {
		// Add circles into a separate SVG group
		svg = svg.append("g")
		const self = this

		const circles = svg.selectAll("circle")
			.data(circleLocations)
			.enter()
			.append("circle");
		circles.attr("cx", d => self.xscale(d.delay))
			.attr("cy", self.yscale(0))
			.attr("r", d => d.wasImputed ? 3 : 5)
			.style("fill", d => d.wasImputed ? "red" : "black");
		circles
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
			.attr("cy", d => self.yscale(d.value))
	}

	animateXYLineWithMean(lineData, className, xMin = 0) {
		this.animateCI(this.ciLayer, [Chart.normativeRange(lineData)])
		this.animateLine(this.meanLayer, [Chart.dataAsXY(lineData, 'delay', 'mean')], "meanline")
		this.drawHorizontalLine(this.linesLayer, 0, xMin)
		this.animateLine(this.valueLayer, [Chart.dataAsXY(lineData)], "line")
		this.animateCircles(this.circlesLayer, lineData)
	}
}

// Set some constants for the class
Object.defineProperty(Chart, 'delayTime', {
	value: 750,
	enumerable: true,
})
Object.defineProperty(Chart, 'transitionTime', {
	value: 2500,
	enumerable: true,
})
