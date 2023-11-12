async function getPosition(street, houseNumber, city, postcode) {
  const query =
    typeof address == "string" ? address : `${address.street} ${address.house_number}, ${address.zip} ${address.city}`;

  const url = `https://nominatim.openstreetmap.org/search?format=json&q=${query}`;
  const response = await fetch(url);

  if (response.ok) {
    const addresses = await response.json();
    if (addresses.length === 0) {
      return;
    }

    return {
      latitude: +addresses[0].lat,
      longitude: +addresses[0].lon
    };
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
