async function getPosition(address) {
  const url = `https://nominatim.openstreetmap.org/search?format=json&q=${address}`;
  const response = await fetch(url);

  if (response.ok) {
    const addresses = await response.json();
    if (addresses.length === 0) {
      return;
    }

    return { lat: +addresses[0].lat, lon: +addresses[0].lon };
  }
}

async function getAddress(coordinates) {
  if (!coordinates) {
    throw new Error("coordinates is either null or length is not 2");
  }

  const url = `https://nominatim.openstreetmap.org/reverse?format=json&lat=${coordinates.lat}&lon=${coordinates.lon}`;
  const response = await fetch(url);

  if (response.ok) {
    const address = await response.json();
    if (!address) {
      return;
    }

    const { road, house_number, city, town, village, postcode } = address.address;
    return `${road} ${house_number || ""}, ${postcode || ""} ${city || town || village || ""}`;
  }
}

const LocationService = {
  _watcherId: -1,
  _address: "",
  _locationChangeHandlers: [],
  _location: { lat: 48.137154, lon: 11.576124 }, // initial position is Munich Marienplatz
  _testIntervalId: -1,

  get location() {
    return this._location;
  },

  set location(newLocation) {
    if (!newLocation) {
      throw new Error("Invalid newLocation");
    }

    this._location = newLocation;
    this.searchAddress().then(() =>
      this._locationChangeHandlers.forEach((handler) => handler(this._location, this._address))
    );

    const form = new FormData();
    form.set("lon", newLocation.lon);
    form.set("lat", newLocation.lat);
    fetch("/api/accounts/location", {
      method: "POST",
      body: form
    });
  },

  get address() {
    return this._address;
  },

  addChangeHandler(handler) {
    this._locationChangeHandlers.push(handler);
    console.log("Number of locationChangeHandlers:", this._locationChangeHandlers.length);
  },

  startLocationWatcher() {
    if (this._watcherId != -1) {
      return;
    }

    console.log("Start location watcher");
    this._watcherId = window.navigator.geolocation.watchPosition(
      async (position) => {
        this.location = { lat: position.coords.latitude, lon: position.coords.longitude };
      },
      (error) => {
        console.error("Error while watching position:", error);
      }
    );
  },

  async searchAddress() {
    this._address = await getAddress(this.location);
  },

  stopLocationWatcher() {
    if (this._watcherId > -1) {
      console.log("Stop location watcher");
      window.navigator.geolocation.clearWatch(this._watcherId);
      this._watcherId = -1;
    }
  },

  toggleLocationServiceSimulation() {
    if (this._testIntervalId > 0) {
      console.log("Stop Location Service Simulation ...");
      clearInterval(this._testIntervalId);
      this._testIntervalId = -1;
    } else {
      console.log("Start Location Service Simulation ...");
      this._testIntervalId = setInterval(() => {
        this.location = {
          lon: this._location.lon - 0.0001,
          lat: this._location.lat + 0.0001
        };
      }, 2000);
    }
  }
};

const FilterService = {
  _searchRadius: 500,
  _selectedCategories: [],
  _searchRadiusChangeListeners: [],
  _selectedCategoriesChangeListeners: [],

  get searchRadius() {
    return this._searchRadius;
  },

  set searchRadius(newSearchRadius) {
    this._searchRadius = newSearchRadius;
    this._searchRadiusChangeListeners.forEach((handler) => handler(newSearchRadius));
  },

  get selectedCategories() {
    return this._selectedCategories;
  },

  set selectedCategories(newSelectedCategories) {
    this._selectedCategories = newSelectedCategories;
    this._selectedCategoriesChangeListeners.forEach((handler) => handler(newSelectedCategories));
  },

  toggleSelectedCategory(category) {
    const index = this._selectedCategories.indexOf(category);

    if (index > -1) {
      this._selectedCategories.splice(index, 1);
    } else {
      this._selectedCategories.push(category);
    }
  },

  addSearchRadiusChangeListener(handler) {
    this._searchRadiusChangeListeners.push(handler);
  }
};
