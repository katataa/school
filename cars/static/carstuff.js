// js code for all car related stuff aka comparing, details, display
let carData = [];
let likedCars = JSON.parse(localStorage.getItem('likedCars')) || [];

document.addEventListener('DOMContentLoaded', () => {
  fetchCarModels();
  fetchCategories();
});

function fetchCarModels() {
  fetch('/carModels')
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to fetch car models');
      }
      return response.json();
    })
    .then(data => {
      carData = data;
      displayCars(carData);
    })
    .catch(error => {
      console.error('Error fetching car models:', error);
      displayErrorMessage("Sorry, this page isn't currently working, try again later! We are working on it. Thank you for your patience! :)");
    });
}

function fetchCategories() {
  fetch('/categories')
    .then(response => {
      if (!response.ok) {
        throw new Error('Failed to fetch categories');
      }
      return response.json();
    })
    .then(categories => {
      const categoryFilter = document.getElementById('category-filter');
      categoryFilter.innerHTML = '<option value="">All Categories</option>';
      
      const desiredCategories = ['SUV', 'Sedan', 'Truck'];
      categories.forEach(category => {
        if (desiredCategories.includes(category.name)) {
          const option = document.createElement('option');
          option.value = category.id;
          option.text = category.name;
          categoryFilter.appendChild(option);
        }
      });
    })
    .catch(error => {
      console.error('Error fetching categories:', error);
    });
}

function displayErrorMessage(message) {
  const mainContent = document.getElementById('main-content');
  if (mainContent) {
    mainContent.classList.add('hidden');
  }

  const body = document.body;
  const errorDiv = document.createElement('div');
  errorDiv.className = 'error-message';
  errorDiv.innerText = message;
  body.appendChild(errorDiv);
}

function displayCars(cars) {
  const container = document.getElementById('grid-container');
  container.innerHTML = '';
  cars.forEach(car => {
    const carDiv = document.createElement('div');
    carDiv.className = 'container';
    carDiv.innerHTML = `
      <img src="${car.image}" alt="${car.name}" class="image">
      <div class="middle">
        <div class="button" onclick="showCarDetail(${car.id})">${car.name}</div>
      </div>
    `;
    container.appendChild(carDiv);
  });
}

function showCarDetail(carId) {
  window.location.href = `details.html?id=${carId}`;
}

document.getElementById('compare-cars').addEventListener('click', () => {
  window.location.href = '/static/compare.html';
});

document.getElementById('search-input').addEventListener('input', (event) => {
  const searchTerm = event.target.value.toLowerCase();
  const filteredCars = carData.filter(car => 
    car.name.toLowerCase().includes(searchTerm) || 
    getManufacturerName(car.manufacturerId).toLowerCase().includes(searchTerm) ||
    car.year.toString().includes(searchTerm)
  );
  displayCars(filteredCars);
});

function getManufacturerName(manufacturerId) {
  const manufacturer = data.Manufacturers.find(manufacturer => manufacturer.id === manufacturerId);
  return manufacturer ? manufacturer.name : '';
}

document.getElementById('category-filter').addEventListener('change', (event) => {
  const selectedCategory = event.target.value;
  const filteredCars = selectedCategory 
    ? carData.filter(car => car.categoryId == selectedCategory)
    : carData;
  displayCars(filteredCars);
});
