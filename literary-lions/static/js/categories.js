function escapeHTML(str) {
    return str
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;")
        .replace(/\\/g, "\\\\")   // Double backslashes
        .replace(/\n/g, "\\n")    // Escape newlines
        .replace(/\r/g, "\\r")    // Escape carriage returns
        .replace(/\t/g, "\\t");   // Escape tabs
}

document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('search-form');
    const searchInput = document.getElementById('search-input');
    const authorFilter = document.getElementById('author-filter');
    const categoryFilter = document.getElementById('category-filter');
    const publishedDateFilter = document.getElementById('published-date-filter');
    const submitButton = document.getElementById('submit-button');

    async function handleSearch(event) {
        event.preventDefault();
        const query = searchInput.value.trim().toLowerCase();
        const author = authorFilter.value.trim().toLowerCase();
        const category = categoryFilter.value.trim().toLowerCase();
        const publishedDate = publishedDateFilter.value.trim().toLowerCase();
        const resultsDiv = document.getElementById('results');
        resultsDiv.innerHTML = 'Searching...';

        console.log("Search parameters:", {
            query,
            author,
            category,
            publishedDate
        });

        try {
            const response = await fetch(`/search-books?q=${encodeURIComponent(query)}&author=${encodeURIComponent(author)}&category=${encodeURIComponent(category)}&publishedDate=${encodeURIComponent(publishedDate)}`);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }

            const data = await response.json();
            console.log("Response data:", data);

            if (data.books.length === 0) {
                resultsDiv.innerHTML = 'Sorry! We couldn\'t find the book you are looking for!';
                return;
            }

            resultsDiv.innerHTML = '';
            data.books.forEach(book => {
                const bookDiv = document.createElement('div');
                bookDiv.classList.add('book');
                bookDiv.innerHTML = `
                    <img src="${book.cover}" alt="${book.title} cover" class="book-cover">
                    <div class="book-details">
                        <h3>${book.title}</h3>
                        <p>Author: ${book.author || 'Unknown'}</p>
                        <p>Category: ${book.category || 'Unknown'}</p>
                        <p>Published Date: ${book.publishedDate || 'Unknown'}</p>
                        <p class="book-description">${book.description || ''}</p>
                        ${book.isFavorite ? '<p class="favorite-text">Book added to favourites!</p>' : '<button type="button" class="add-to-favourites" data-book=\'' + escapeHTML(JSON.stringify(book)) + '\'>Add to Favourite Books</button>'}
                    </div>
                `;
                resultsDiv.appendChild(bookDiv);
            });

            document.querySelectorAll('.add-to-favourites').forEach(button => {
                button.addEventListener('click', async function() {
                    const book = JSON.parse(this.dataset.book);
                    try {
                        const response = await fetch('/add-to-favorites', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json'
                            },
                            body: JSON.stringify(book)
                        });
                        if (!response.ok) {
                            throw new Error('Network response was not ok');
                        }
                        const span = document.createElement('p');
                        span.textContent = 'Book added to favourites!';
                        span.classList.add('favorite-text');
                        span.style.marginTop = '20px';
                        this.replaceWith(span);
                    } catch (error) {
                        console.error("Error adding book to favourites:", error);
                        alert('An error occurred while adding the book to favourites.');
                    }
                });
            });
        } catch (error) {
            resultsDiv.innerHTML = 'Sorry! We couldn\'t find the book you are looking for!';
            console.error("Error:", error);
        }
    }

    searchForm.addEventListener('submit', handleSearch);

    [searchInput, authorFilter, categoryFilter, publishedDateFilter].forEach(input => {
        input.addEventListener('keypress', function(event) {
            if (event.key === 'Enter') {
                event.preventDefault(); 
                handleSearch(event);   
            }
        });
    });

    const postForm = document.getElementById('post-form');
    postForm.addEventListener('keypress', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault(); 
            submitButton.click();   
        }
    });
});
