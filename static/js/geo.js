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

const LocationService = {
  _watcherId: -1,
  _address: "",
  _locationChangeHandlers: [],
  _location: [48.137154, 11.576124], // initial position is Munich Marienplatz

  get location() {
    return this._location;
  },

  set location(newLocation) {
    this._location = newLocation;
    this.searchAddress().then(() => this._locationChangeHandlers.forEach((handler) => handler(this._location, this._address)));
  },
  
  get address() {
    return this._address;
  },

  setLocationFromString: function (locationString) {
    if (locationString.length > 0) {
      this.location = locationString
        .split(",")
        .map((n) => Number(n))
        .reverse();
    }
  },

  addLocationChangeHandler: function (handler) {
    this._locationChangeHandlers.push(handler);
    console.log("Number of locationChangeHandlers:", this._locationChangeHandlers.length);
  },

  startLocationWatcher: function () {
    if (this._watcherId != -1) {
      return;
    }

    console.log("Start location watcher");
    this._watcherId = window.navigator.geolocation.watchPosition(
      async (position) => {
        this.location = [position.coords.latitude, position.coords.longitude];
      },
      (error) => {
        console.error("Error while watching position:", error);
      }
    );
  },

  searchAddress: async function() {
    if (this.location.length == 2) {
      this._address = await getAddress(this.location[0], this.location[1]);
    }
  },

  stopLocationWatcher: function () {
    if (this._watcherId > -1) {
      console.log("Stop location watcher");
      window.navigator.geolocation.clearWatch(this._watcherId);
      this._watcherId = -1;
    }
  }
};
