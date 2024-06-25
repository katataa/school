// js code for search bar and suggestions
const searchInput = document.querySelector("[data-search]");
const suggestionsBox = document.getElementById('suggestions');

searchInput.addEventListener('input', (e) => {
  const value = e.target.value.toLowerCase();
  suggestionsBox.innerHTML = '';
  if (value) {
    const suggestions = carData.filter(car => car.name.toLowerCase().includes(value));
    if (suggestions.length > 0) {
      suggestionsBox.classList.add('has-suggestions');
    } else {
      suggestionsBox.classList.remove('has-suggestions');
    }
    suggestions.forEach(car => {
      const suggestionDiv = document.createElement('div');
      suggestionDiv.className = 'suggestion';
      suggestionDiv.textContent = car.name;
      suggestionDiv.addEventListener('click', () => {
        showCarDetail(car.id);
      });
      suggestionsBox.appendChild(suggestionDiv);
    });
  } else {
    suggestionsBox.classList.remove('has-suggestions');
  }
});

searchInput.addEventListener('keypress', (e) => {
  if (e.key === 'Enter') {
    const value = searchInput.value.toLowerCase();
    const car = carData.find(car => car.name.toLowerCase().includes(value));
    if (car) {
      showCarDetail(car.id);
    }
  }
});
