document.addEventListener("DOMContentLoaded", function() {
  fetch('/recommendations')
    .then(response => {
      console.log("Response status:", response.status);
      return response.json();
    })
    .then(data => {
      console.log("Fetched recommendations:", data);
      const recommendationList = document.getElementById('recommendation-list');
      recommendationList.innerHTML = '';

      if (data.length === 0) {
        recommendationList.innerHTML = '<p>No recommendations available.</p>';
        return;
      }

      data.forEach(car => {
        const carElement = document.createElement('div');
        carElement.className = 'car-card';
        carElement.innerHTML = `
          <h3>${car.name}</h3>
          <img src="${car.image}" alt="${car.name}" class="car-image">
        `;
        carElement.addEventListener('click', function() {
          window.location.href = `/details.html?id=${car.id}`;
        });
        recommendationList.appendChild(carElement);
      });
    })
    .catch(error => console.error('Error fetching recommendations:', error));
});
