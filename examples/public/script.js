console.log('Duwa HTTP Server - JavaScript file loaded successfully!');

document.addEventListener('DOMContentLoaded', function() {
    const body = document.body;
    
    // Add a success message
    const successDiv = document.createElement('div');
    successDiv.className = 'success';
    successDiv.innerHTML = '✅ JavaScript is working! File served successfully by Duwa HTTP server.';
    
    // Insert after the first element
    if (body.children.length > 0) {
        body.insertBefore(successDiv, body.children[1]);
    } else {
        body.appendChild(successDiv);
    }
    
    // Add click handler to all links
    const links = document.querySelectorAll('a');
    links.forEach(link => {
        link.addEventListener('click', function(e) {
            console.log('Clicked link:', this.href);
        });
    });
});