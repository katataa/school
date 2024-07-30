document.addEventListener('DOMContentLoaded', function() {
    // Like/Dislike functionality
    document.querySelectorAll('.like-dislike-form').forEach(form => {
        form.addEventListener('submit', async function(event) {
            event.preventDefault();
            
            const formData = new FormData(this);
            const postID = formData.get('post_id');
            const commentID = formData.get('comment_id');
            const isLike = formData.get('is_like') === 'true';
            
            const payload = {};
            if (postID) payload.post_id = parseInt(postID, 10);
            if (commentID) payload.comment_id = parseInt(commentID, 10);
            payload.is_like = isLike;

            try {
                const response = await fetch('/like_dislike', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload)
                });

                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }

                const result = await response.json();

                if (result.post_id) {
                    const postDiv = document.querySelector(`.post[data-id="${result.post_id}"]`);
                    postDiv.querySelector('.post-likes').textContent = result.likes;
                    postDiv.querySelector('.post-dislikes').textContent = result.dislikes;
                } else if (result.comment_id) {
                    const commentDiv = document.querySelector(`.comment[data-id="${result.comment_id}"]`);
                    commentDiv.querySelector('.comment-likes').textContent = result.likes;
                    commentDiv.querySelector('.comment-dislikes').textContent = result.dislikes;
                }
            } catch (error) {
                console.error('Error processing like/dislike:', error);
            }
        });
    });

    // Prevent form submission on Enter key press in text fields
    document.querySelectorAll('form').forEach(form => {
        form.addEventListener('keypress', function(event) {
            if (event.key === 'Enter') {
                const activeElement = document.activeElement;
                if (activeElement.tagName === 'TEXTAREA' || activeElement.tagName === 'INPUT') {
                    event.preventDefault();  
                    form.querySelector('button[type="submit"]').click();  
                }
            }
        });
    });

    // Modal functionality for images
    var modal = document.getElementById("imageModal");
    var modalImg = document.getElementById("modalImage");
    var captionText = document.getElementById("caption");
    var closeModal = document.getElementsByClassName("close")[0];

    document.querySelectorAll(".post img").forEach(function(img) {
        img.addEventListener("click", function() {
            modal.style.display = "block";
            modalImg.src = this.src;
            captionText.innerHTML = this.alt;
        });
    });

    closeModal.onclick = function() {
        modal.style.display = "none";
    }

    modal.onclick = function(event) {
        if (event.target === modal) {
            modal.style.display = "none";
        }
    }
});
