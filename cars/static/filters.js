// all js code related to filtering
function filterCars(filter) {
  if (filter === 'all') {
    displayCars(carData);
  } else if (filter === 'liked') {
    const likedCarData = carData.filter(car => likedCars.includes(car.id));
    displayCars(likedCarData);
  } else {
    const filteredCars = carData.filter(car => car.specifications.drivetrain === filter);
    displayCars(filteredCars);
  }
}

function sortCars(order) {
  let sortedCars = [...carData];
  if (order === 'A-Z') {
    sortedCars.sort((a, b) => a.name.localeCompare(b.name));
  } else if (order === 'Z-A') {
    sortedCars.sort((a, b) => b.name.localeCompare(a.name));
  }
  displayCars(sortedCars);
}

document.querySelectorAll('#buttons .button-value').forEach(button => {
  button.addEventListener('click', () => {
    const filter = button.getAttribute('data-filter');
    const sort = button.getAttribute('data-sort');
    if (filter) {
      filterCars(filter);
    } else if (sort) {
      sortCars(sort);
    }
  });
});
