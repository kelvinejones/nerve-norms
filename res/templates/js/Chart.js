class Chart {
	constructor(xRange, yRange, xScaleType = Chart.scaleType.LINEAR, yScaleType = Chart.scaleType.LINEAR) {
		this.xRange = xRange
		this.yRange = yRange

		this.margin = { top: 20, right: 20, bottom: 50, left: 50 };
		this.width = 600 - this.margin.left - this.margin.right;
		this.height = 300 - this.margin.top - this.margin.bottom;

		// Set the default scaling
		this.xscale = this.makeScale(xScaleType).range([0, this.width]).domain(xRange);
		this.yscale = this.makeScale(yScaleType).range([this.height, 0]).domain(yRange);

		// Set default config values
		this.xName = 'delay'
		this.yName = 'value'
		this.ySDName = 'SD'
		this.yMeanName = 'mean'

		this.yAnimStart = this.animationStartValue(this.yRange)

		this.group = {}
	}

	makeScale(name) {
		let scale
		switch (name) {
			case "LINEAR":
				scale = d3.scaleLinear()
				scale.scaleType = Chart.scaleType.LINEAR
				break
			case "LOG":
				scale = d3.scaleLog()
				scale.scaleType = Chart.scaleType.LOG
				break
		}
		return scale
	}

	animationStartValue(range) {
		// The animation start value should be at zero unless that's outside of range
		if (range[0] < 0 && range[1] >= 0) {
			return 0
		} else {
			return range[0]
		}
	}

	updatePlots(plots) { throw new Error("A Chart must implement updatePlots(plots)") }
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
			svg.attr("transform", "scale(0.9) translate(110, 0)")
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

	// standardDeviationCI creates a CI polygon showing the area within a specified number of standard deviations
	standardDeviationCI(data, numSD = 2) {
		const xMean = this.xMeanName || this.xName
		const yMean = this.yMeanName || this.yName
		const ySD = this.ySDName
		const xSD = this.xSDName
		const first = data[0]
		const last = data[data.length - 1]
		if (first[xMean] === undefined || first[yMean] === undefined) {
			return []
		}
		return this.scaleArrayWithinRange((Array.from(data)
				.map(function(d) { return { x: d[xMean] - numSD * (d[xSD] || 0), y: d[yMean] + numSD * (d[ySD] || 0) } }))
			.concat({ x: last[xMean] + numSD * (last[xSD] || 0), y: last[yMean] + numSD * (last[ySD] || 0) })
			.concat(Array.from(data).reverse().map(function(d) { return { x: d[xMean] + numSD * (d[xSD] || 0), y: d[yMean] - numSD * (d[ySD] || 0) } }))
			.concat({ x: first[xMean] - numSD * (first[xSD] || 0), y: first[yMean] - numSD * (first[ySD] || 0) }))
	}

	// normativeLimits extracts the calculated limits from the dataset, which describes the range in which a healthy measure is expected
	normativeLimits(data) {
		const xMean = this.xMeanName || this.xName
		const yMean = this.yMeanName || this.yName
		const last = data[data.length - 1]
		const first = data[0]
		if (first.leftLimit !== undefined) {
			// This is a complicated limit with x and y limits.
			return this.scaleArrayWithinRange((Array.from(data)
					.map(function(d) { return { x: d.leftLimit, y: d.upperLimit || d[yMean] } }))
				.concat({ x: last.rightLimit, y: last.upperLimit || last[yMean] })
				.concat(Array.from(data).reverse().map(function(d) { return { x: d.rightLimit, y: d.lowerLimit || d[yMean] } }))
				.concat({ x: first.leftLimit, y: first.lowerLimit || first[yMean] }))
		} else if (first.upperLimit !== undefined) {
			// This is a simple limit with upper and lower bounds.
			return this.scaleArrayWithinRange((Array.from(data)
					.map(function(d) { return { x: d[xMean], y: d.lowerLimit || d[yMean] } }))
				.concat({ x: last[xMean], y: last.upperLimit || last[yMean] })
				.concat(Array.from(data).reverse().map(function(d) { return { x: d[xMean], y: d.upperLimit || d[yMean] } }))
				.concat({ x: first[xMean], y: first.lowerLimit || first[yMean] }))
		} else {
			return []
		}
	}

	scaleArrayWithinRange(ar) {
		// With a log scale, values can't be plotted at or below zero.
		if (this.xscale.scaleType == Chart.scaleType.LOG) {
			ar = this.raiseZeroValues(ar, 'x', this.xRange[0])
		}
		if (this.yscale.scaleType == Chart.scaleType.LOG) {
			ar = this.raiseZeroValues(ar, 'y', this.yRange[0])
		}
		return ar
	}

	raiseZeroValues(ar, axis, min) {
		return ar.map(function(d) {
			if (d[axis] < min) {
				d[axis] = min
			}
			return d
		})
	}

	dataAsXY(data, xName, yName) {
		if (data[0][xName] === undefined || data[0][yName] === undefined) {
			return []
		}
		return data.map(function(d) { return { x: d[xName], y: d[yName] } })
	}

	xZeroLine() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(this.yAnimStart))
	}

	xyLine() {
		return d3.line()
			.x(d => this.xscale(d.x))
			.y(d => this.yscale(d.y));
	}

	drawHorizontalLine(svg, yVal) {
		svg.append("path")
			.data([
				[{ x: this.xRange[0], y: yVal }, { x: this.xRange[1], y: yVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyLine());
	}

	drawVerticalLine(svg, xVal) {
		svg.append("path")
			.data([
				[{ y: this.yRange[0], x: xVal }, { y: this.yRange[1], x: xVal }]
			])
			.attr("class", "meanline")
			.attr("d", this.xyLine());
	}

	createGroup(svg, typeString, name) {
		svg = svg.append("g")
		this.group[typeString + "-" + name] = svg
		return svg
	}

	animateGroup(typeString, newData, name) {
		return this.group[typeString + "-" + name].selectAll(typeString)
			.data(newData)
			.transition()
			.delay(Chart.delayTime)
			.duration(Chart.transitionTime)
	}

	animateCI(ciNormRange, name) {
		this.animateGroup("path", ciNormRange, name + "-" + "confidenceinterval")
			.attr("d", this.xyLine());
	}

	createLine(svg, xyLine, groupName, className) {
		this.createGroup(svg, "path", groupName + "-" + className)
			.append("path")
			.data(xyLine)
			.attr("class", className)
			.attr("d", this.xZeroLine())
	}

	animateLine(xyLine, groupName, className) {
		this.animateGroup("path", xyLine, groupName + "-" + className)
			.attr("d", this.xyLine());
	}

	createCircles(svg, circleLocations, name) {
		// create circle locations at init for each name, with right amount and position of circles
		// Add circles into a separate SVG group
		this.createGroup(svg, "circle", name)
			.selectAll("circle")
			.data(circleLocations)
			.enter()
			.append("circle")
			.attr("cx", d => this.xscale(d[this.xName]))
			.attr("cy", this.yscale(this.yAnimStart))
			.attr("r", d => 3)
			.style("fill", d => "black");
	}

	animateCircles(circleLocations, name) {
		this.animateGroup("circle", circleLocations, name)
			.attr("r", d => d.wasImputed ? 3 : 5)
			.style("fill", d => d.wasImputed ? "red" : "black")
			.attr("cy", d => this.yscale(d[this.yName]))
	}

	createXYLineWithMean(lineData, name) {
		this.createLine(this.ciLayer, [this.normativeLimits(lineData)], name, "confidenceinterval")
		this.createLine(this.meanLayer, [this.dataAsXY(lineData, this.xMeanName || this.xName, this.yMeanName)], name, "meanline")
		this.createLine(this.valueLayer, [this.dataAsXY(lineData, this.xName, this.yName)], name, "line")
		this.createCircles(this.circlesLayer, lineData, name)
	}

	animateXYLineWithMean(lineData, name) {
		this.animateCI([this.normativeLimits(lineData)], name)
		this.animateLine([this.dataAsXY(lineData, this.xMeanName || this.xName, this.yMeanName)], name, "meanline")
		this.animateLine([this.dataAsXY(lineData, this.xName, this.yName)], name, "line")
		this.animateCircles(lineData, name)
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
Object.defineProperty(Chart, 'scaleType', {
	value: {
		LINEAR: "LINEAR",
		LOG: "LOG",
	},
	enumerable: true,
})
