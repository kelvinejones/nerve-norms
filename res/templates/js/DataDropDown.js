class DataDropDown {
	// elementID is the HTML ID.
	// data is an object indexed by the drop-down's values
	// action is a function that receives the selection name and data as its arguments
	constructor(elementID, data, action) {
		this.elementID = elementID
		this.data = data
		this.action = action
		this.val = document.getElementById(this.elementID).value
		document.getElementById(elementID).addEventListener("change", this.updateData);

		this.updateData = this.updateData.bind(this);
	}

	updateData = (ev) => {
		this.val = ev.srcElement.value
		this.action(this.val, this.data[this.val])
	}
}
