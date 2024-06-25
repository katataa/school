// js code that shows liked cars and keeps them in memory
let carData = [];
let likedCars = JSON.parse(localStorage.getItem('likedCars')) || [];
const userID = localStorage.getItem('userID') || generateUserID();

function generateUserID() {
  const newUserID = 'user-' + Math.random().toString(36).substring(2, 15);
  localStorage.setItem('userID', newUserID);
  return newUserID;
}

document.addEventListener('DOMContentLoaded', () => {
  fetchLikedCars();
  fetchCarModels();
});

function fetchCarModels() {
  fetch('/carModels')
    .then(response => response.json())
    .then(data => {
      carData = data;
      displayCars(carData);
      updateLikeIcons();
    })
    .catch(error => console.error('Error fetching car models:', error));
}

function fetchLikedCars() {
  fetch('/likedCars', {
    headers: {
      'User-ID': userID
    }
  })
    .then(response => response.json())
    .then(data => {
      likedCars = data.map(car => car.id);
      localStorage.setItem('likedCars', JSON.stringify(likedCars));
      updateLikeIcons();
    })
    .catch(error => console.error('Error fetching liked cars:', error));
}

function displayCars(cars) {
  const container = document.getElementById('grid-container');
  container.innerHTML = '';
  cars.forEach(car => {
    const carDiv = document.createElement('div');
    carDiv.className = 'container';
    const isLiked = likedCars.includes(car.id);
    carDiv.innerHTML = `
      <img src="/img/${car.image}" alt="${car.name}" class="image">
      <div class="middle">
        <div class="button" onclick="showCarDetail(${car.id})">${car.name}</div>
        <div class="like-button" id="like-icon-${car.id}" onclick="toggleLike(${car.id})">
          <span>${isLiked ? '♥' : '♡'}</span>
        </div>
      </div>
    `;
    container.appendChild(carDiv);
  });
}

function updateLikeIcons() {
  likedCars.forEach(carId => {
    const likeIcon = document.getElementById(`like-icon-${carId}`);
    if (likeIcon) {
      likeIcon.innerHTML = '<span>♥</span>';
    }
  });
}

function toggleLike(carId) {
  const index = likedCars.indexOf(carId);
  const isLiked = index > -1;

  if (isLiked) {
    likedCars.splice(index, 1);
  } else {
    likedCars.push(carId);
  }

  localStorage.setItem('likedCars', JSON.stringify(likedCars));
  updateLikeIcons();

  fetch(`/likeCar?car_model_id=${carId}`, {
    method: 'POST',
    headers: {
      'User-ID': userID
    }
  }).then(response => {
    if (!response.ok) {
      console.error('Failed to update liked car status on the server');
    }
  }).catch(error => console.error('Error updating liked car status:', error));
}

document.querySelector('[data-filter="liked"]').addEventListener('click', function() {
  const likedCarData = carData.filter(car => likedCars.includes(car.id));
  displayCars(likedCarData);
});

function showCarDetail(carId) {
  window.location.href = `details.html?id=${carId}`;
}
