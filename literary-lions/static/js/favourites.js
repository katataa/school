document.addEventListener('DOMContentLoaded', function() {
    console.log("favourites.js loaded");

    const buttons = document.querySelectorAll('.btn-delete');
    console.log("Found delete buttons:", buttons.length); 

    buttons.forEach(button => {
        console.log("Attaching event listener to delete-from-favourites button");  
        button.addEventListener('click', async function() {
            const title = this.getAttribute('data-title');
            console.log("Book title on button click:", title); 

            try {
                const response = await fetch('/delete-from-favourites', {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ title })
                });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                console.log("Book deleted from favourites successfully"); 
                this.closest('.book').remove(); 
            } catch (error) {
                console.error("Error deleting book from favourites:", error);
                alert('An error occurred while deleting the book from favourites.');
            }
        });
    });
});
