function drawModal(name, chart) {
	document.getElementById('modal-title').innerHTML = name
	d3.selectAll("#modal svg > *").remove();
	chart.draw(d3.select('#modal svg'));
	$('#modal').modal('toggle');
}
