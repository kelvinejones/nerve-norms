function normativeRange(data) {
	return (Array.from(data)
			.map(function(d) { return { x: d.delay, y: d.mean + 2 * d.SD } }))
		.concat(Array.from(data).reverse().map(function(d) { return { x: d.delay, y: d.mean - 2 * d.SD } }))
}
