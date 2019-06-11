class Participant {
	static load(name, array) {
		return new Participant(array[name], name)
	}
	constructor(data, name) {
		this._data = data
		this._localData = (name != null)
		if (this._localData) {
			this._name = name
		} else {
			this._name = this._data.header.name
		}
	}

	get dataIsLocal() { return this._localData }
	get name() { return this._name }
	get data() { return this._data }
}
