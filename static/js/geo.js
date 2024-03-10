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
  _location: { lon: 48.137154, lat: 11.576124 }, // initial position is Munich Marienplatz

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

  addChangeHandler: function (handler) {
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
        this.location = { lat: position.coords.latitude, lon: position.coords.longitude };
      },
      (error) => {
        console.error("Error while watching position:", error);
      }
    );
  },

  searchAddress: async function () {
    this._address = await getAddress(this.location);
  },

  stopLocationWatcher: function () {
    if (this._watcherId > -1) {
      console.log("Stop location watcher");
      window.navigator.geolocation.clearWatch(this._watcherId);
      this._watcherId = -1;
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

  toggleSelectedCategory: function (category) {
    const index = this._selectedCategories.indexOf(category);

    if (index > -1) {
      this._selectedCategories.splice(index, 1);
    } else {
      this._selectedCategories.push(category);
    }
  },

  addSearchRadiusChangeListener: function (handler) {
    this._searchRadiusChangeListeners.push(handler);
  }
};
