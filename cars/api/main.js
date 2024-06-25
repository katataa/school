const express = require('express');
const app = express();
const port = 8080;

// Serve static files
app.use('/static', express.static('static'));

app.get('/api/manufacturers', (req, res) => {
  console.log('Received request for manufacturers');
  try {
    const manufacturers = [
      { id: 1, name: 'Toyota', country: 'Japan', foundingYear: 1937 },
      { id: 2, name: 'Honda', country: 'Japan', foundingYear: 1948 },
      { id: 3, name: 'BMW', country: 'Germany', foundingYear: 1916 },
      { id: 4, name: 'Audi', country: 'Germany', foundingYear: 1910 },
      { id: 5, name: 'Mercedes-Benz', country: 'Germany', foundingYear: 1926 },
      { id: 6, name: 'Ford', country: 'United States', foundingYear: 1903 },
      { id: 7, name: 'Chevrolet', country: 'United States', foundingYear: 1911 },
      { id: 8, name: 'Hyundai', country: 'South Korea', foundingYear: 1967 },
      { id: 9, name: 'Lexus', country: 'Japan', foundingYear: 1989 },
      { id: 10, name: 'Nissan', country: 'Japan', foundingYear: 1933 }
    ];
    res.json(manufacturers);
  } catch (error) {
    console.error('Error fetching manufacturers:', error);
    res.status(500).send('Internal Server Error');
  }
});

app.get('/api/categories', (req, res) => {
  console.log('Received request for categories');
  try {
    const categories = [
      { id: 1, name: 'SUV' },
      { id: 2, name: 'Sedan' },
      { id: 3, name: 'Coupe' },
      { id: 4, name: 'Truck' },
      { id: 5, name: 'Hatchback' },
      { id: 6, name: 'Convertible' },
      { id: 7, name: 'Wagon' },
      { id: 8, name: 'Electric' },
      { id: 9, name: 'Luxury' },
      { id: 10, name: 'Sports' }
    ];
    res.json(categories);
  } catch (error) {
    console.error('Error fetching categories:', error);
    res.status(500).send('Internal Server Error');
  }
});

app.get('/api/carModels', (req, res) => {
  console.log('Received request for car models');
  try {
    const carModels = [
      { id: 1, name: 'Toyota Corolla', manufacturerId: 1, categoryId: 2, year: 2023, specifications: { engine: '1.8L Inline-4', horsepower: 139, transmission: 'CVT', drivetrain: 'Front-Wheel Drive' }, image: '/static/images/toyota_corolla.jpg' },
      { id: 2, name: 'Honda Civic', manufacturerId: 2, categoryId: 2, year: 2024, specifications: { engine: '2.0L Inline-4', horsepower: 158, transmission: '6-speed Manual', drivetrain: 'Front-Wheel Drive' }, image: '/static/images/honda_civic.jpg' },
      { id: 3, name: 'BMW 3 Series', manufacturerId: 3, categoryId: 2, year: 2023, specifications: { engine: '2.0L Inline-4', horsepower: 255, transmission: '8-speed Automatic', drivetrain: 'Rear-Wheel Drive' }, image: '/static/images/bmw_3series.jpg' },
      { id: 4, name: 'Audi A4', manufacturerId: 4, categoryId: 2, year: 2024, specifications: { engine: '2.0L Inline-4', horsepower: 201, transmission: '7-speed Automatic', drivetrain: 'All-Wheel Drive' }, image: '/static/images/audi_a4.jpg' },
      { id: 5, name: 'Mercedes-Benz E-Class', manufacturerId: 5, categoryId: 2, year: 2023, specifications: { engine: '2.0L Inline-4', horsepower: 255, transmission: '9-speed Automatic', drivetrain: 'Rear-Wheel Drive' }, image: '/static/images/mercedes_eclass.jpg' },
      { id: 6, name: 'Ford F-150', manufacturerId: 6, categoryId: 4, year: 2024, specifications: { engine: '3.5L V6', horsepower: 400, transmission: '10-speed Automatic', drivetrain: 'Rear-Wheel Drive' }, image: '/static/images/ford_f150.jpg' },
      { id: 7, name: 'Chevrolet Silverado', manufacturerId: 7, categoryId: 4, year: 2023, specifications: { engine: '5.3L V8', horsepower: 355, transmission: '8-speed Automatic', drivetrain: 'Rear-Wheel Drive' }, image: '/static/images/chevrolet_silverado.jpg' },
      { id: 8, name: 'Hyundai Sonata', manufacturerId: 8, categoryId: 2, year: 2024, specifications: { engine: '2.5L Inline-4', horsepower: 191, transmission: '8-speed Automatic', drivetrain: 'Front-Wheel Drive' }, image: '/static/images/hyundai_sonata.jpg' },
      { id: 9, name: 'Lexus RX', manufacturerId: 9, categoryId: 1, year: 2023, specifications: { engine: '3.5L V6', horsepower: 295, transmission: '8-speed Automatic', drivetrain: 'Front-Wheel Drive' }, image: '/static/images/lexus_rx.jpg' },
      { id: 10, name: 'Nissan Altima', manufacturerId: 10, categoryId: 2, year: 2024, specifications: { engine: '2.5L Inline-4', horsepower: 188, transmission: 'CVT', drivetrain: 'Front-Wheel Drive' }, image: '/static/images/nissan_altima.jpg' }
    ];
    res.json(carModels);
  } catch (error) {
    console.error('Error fetching car models:', error);
    res.status(500).send('Internal Server Error');
  }
});

app.listen(port, () => {
  console.log(`API Server is running on http://localhost:${port}`);
});
