async function getPosition(address) {
  const url = `https://nominatim.openstreetmap.org/search?format=json&q=${address}`;
  const response = await fetch(url);

  if (response.ok) {
    const addresses = await response.json();
    if (addresses.length === 0) {
      return;
    }

    return [+addresses[0].lat, +addresses[0].lon];
  }
}

async function getAddress(latitude, longitude) {
  const url = `https://nominatim.openstreetmap.org/reverse?format=json&lat=${latitude}&lon=${longitude}`;
  const response = await fetch(url);

  if (response.ok) {
    const address = await response.json();
    if (!address) {
      return;
    }

    const { road, house_number, city, postcode } = address.address;
    return `${road} ${house_number}, ${postcode} ${city}`;
  }
}

const Location = {
  _changeHandlers: [],
  _value: [],
  get location() {
    return this._value;
  },
  set location(newLocation) {
    this._value = newLocation;
    this._changeHandlers.forEach(handler => handler(this._value));
  },
  addChangeHandler: function(handler) {
    this._changeHandlers.push(handler);
  }
}
