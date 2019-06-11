class Participant {
	static load(name, array) {
		return new Participant(array[name], name)
	}
	constructor(data, name, dataIsLocal = true) {
		this._data = data
		this._name = name
		this._localData = dataIsLocal
	}

	get dataIsLocal() { return this._localData }
	get name() { return this._name }
	get data() { return this._data }
}
